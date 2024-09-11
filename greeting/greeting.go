package greeting

import (
	"context"
	"errors"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	TaskQueue = "greeting-tasks"
)

type WorkflowInput struct {
	Name         string
	LanguageCode string
}

type WorkflowOutput struct {
	GreetingMessage string
	GoodbyeMessage  string
}

type ActivityInput struct {
	Name         string
	LanguageCode string
}

type ActivityOutput struct {
	Message string
}

func GreetSomeone(ctx workflow.Context, input WorkflowInput) (WorkflowOutput, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	greetingInput := ActivityInput{Name: input.Name, LanguageCode: input.LanguageCode}
	var spanishGreeting ActivityOutput
	err := workflow.ExecuteActivity(ctx, GreetInSpanish, greetingInput).Get(ctx, &spanishGreeting)
	if err != nil {
		return WorkflowOutput{}, err
	}

	goodbyeInput := ActivityInput{Name: input.Name, LanguageCode: input.LanguageCode}
	var spanishFarewell ActivityOutput
	err = workflow.ExecuteActivity(ctx, FarewellInSpanish, goodbyeInput).Get(ctx, &spanishFarewell)
	if err != nil {
		return WorkflowOutput{}, err
	}

	output := WorkflowOutput{
		GreetingMessage: spanishGreeting.Message,
		GoodbyeMessage:  spanishFarewell.Message,
	}
	return output, nil
}

func FarewellInSpanish(ctx context.Context, input ActivityInput) (ActivityOutput, error) {
	return callService(ctx, "get-spanish-farewell", input)
}

func GreetInSpanish(ctx context.Context, input ActivityInput) (ActivityOutput, error) {
	return callService(ctx, "get-spanish-greeting", input)
}

func callService(ctx context.Context, stem string, input ActivityInput) (ActivityOutput, error) {
	base := "http://localhost:9999/%s?name=%s&lang=%s"
	endpoint := fmt.Sprintf(base, stem, url.QueryEscape(input.Name), url.QueryEscape(input.LanguageCode))

	resp, err := http.Get(endpoint)
	if err != nil {
		return ActivityOutput{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ActivityOutput{}, err
	}

	translation := string(body)
	status := resp.StatusCode
	if status >= 400 {
		message := fmt.Sprintf("HTTP Error %d: %s", status, translation)
		return ActivityOutput{}, errors.New(message)
	}

	return ActivityOutput{Message: translation}, nil
}
