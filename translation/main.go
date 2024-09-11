package main

// call this via HTTP GET with a URL like:
//     http://localhost:9999/get-spanish-greeting?name=Bob
//     http://localhost:9999/get-spanish-farewell?name=Bob

import (
	"fmt"
	"net/http"
)

var languagesGreeting = map[string]string{
	"pt": "Olá %s!",
	"en": "Hello %s!",
	"es": "¡Hola, %s!",
}
var languagesFarewell = map[string]string{
	"pt": "Tchau %s!",
	"en": "Goodbye %s!",
	"es": "¡Adiós, %s!",
}

func spanishGreetingHandler(w http.ResponseWriter, r *http.Request) {
	lang := "en-us"
	langs, ok := r.URL.Query()["lang"]
	if ok && languagesGreeting[langs[0]] != "" {
		lang = langs[0]
	}

	keys, ok := r.URL.Query()["name"]
	if ok {
		name := keys[0]
		translation := fmt.Sprintf(languagesGreeting[lang], name)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, translation)
	} else {
		http.Error(w, "Missing required 'name' parameter.", http.StatusBadRequest)
	}
}

func spanishFarewellHandler(w http.ResponseWriter, r *http.Request) {
	lang := "en"
	langs, ok := r.URL.Query()["lang"]
	if ok && languagesFarewell[langs[0]] != "" {
		lang = langs[0]
	}

	keys, ok := r.URL.Query()["name"]
	if ok {
		name := keys[0]
		translation := fmt.Sprintf(languagesFarewell[lang], name)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, translation)
	} else {
		http.Error(w, "Missing required 'name' parameter.", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/get-spanish-greeting", spanishGreetingHandler)
	http.HandleFunc("/get-spanish-farewell", spanishFarewellHandler)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
