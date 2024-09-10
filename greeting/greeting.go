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

func GreetSomeone(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var spanishGreeting string
	err := workflow.ExecuteActivity(ctx, GreetInSpanish, name).Get(ctx, &spanishGreeting)
	if err != nil {
		return "", err
	}

	var spanishFarewell string
	err = workflow.ExecuteActivity(ctx, FarewellInSpanish, name).Get(ctx, &spanishFarewell)
	if err != nil {
		return "", err
	}
	var helloGoodbye = "\n" + spanishGreeting + "\n" + spanishFarewell
	return helloGoodbye, nil
}

func FarewellInSpanish(ctx context.Context, name string) (string, error) {
	greeting, err := callService(ctx, "get-spanish-farewell", name)
	return greeting, err
}

func GreetInSpanish(ctx context.Context, name string) (string, error) {
	return callService(ctx, "get-spanish-greeting", name)
}

func callService(ctx context.Context, stem, name string) (string, error) {
	base := "http://localhost:9999/%s?name=%s"
	endpoint := fmt.Sprintf(base, stem, url.QueryEscape(name))

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	translation := string(body)
	status := resp.StatusCode
	if status >= 400 {
		message := fmt.Sprintf("HTTP Error %d: %s", status, translation)
		return "", errors.New(message)
	}

	return translation, nil
}
