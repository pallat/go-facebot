package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	appSecret       = "dfcff3882e7e584ea07ca040f3aeb058"
	pageAccessToken = "EAACEdEose0cBAMogBChnP3szo5nifgOMwSl1uZARWRS5o21FUSWpWRWCgRFqgQJ0eriejNDBBmAxW86hD58IH3RBhKohQckRoKMEfxtt2ZBfUBUHkkwHcavLrMnPI7z8iYBFmP2MANXRuWek5hZCyXZCYK11GqD6mAgZBItTb7QZDZD"
	validationToken = "EAACEdEose0cBAPWQN6lEeX3FO3NHwzTKZAbKHgaQpddehI4kFDxyUttiN1FD9Hk1bp5pHLgpt2BaZCuvDKu2JxAD0CuTwpzDEaGmPFVryB8wAwI665YMQYxtL5LUHAu8TMjSqGZAUq1ZADUCLotWssrLugPwO6SjM5SV92bgLwZDZD"
)

var (
	allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowedHeaders = []string{
		"Accept",
		"Authorization",
		"X-Real-IP",
		"Content-Type",
		"X-Custom-Header",
		"Query",
		"Language",
		"Origin",
	}
)

func GetHook(w rest.ResponseWriter, r *rest.Request) {
	if r.FormValue("hub.mode") == "subscribe" && r.FormValue("hub.verify_token") == validationToken {
		log.Println("Validating webhook")
		w.WriteJson(r.FormValue("hub.challenge"))
		return
	}
	log.Println("Failed validation. Make sure the validation tokens match.")
	w.WriteHeader(403)
}

func PostHook(w rest.ResponseWriter, r *rest.Request) {
	data := Common{}
	err := r.DecodeJsonPayload(&data)
	if err != nil {
		w.WriteJson(map[string]string{"error": err.Error()})
		w.WriteHeader(500)
		return
	}

	if data.Object != "page" {
		return
	}

	for _, entry := range data.Entry {
		for _, msg := range entry.Messagings {
			if msg.Message != nil {
				receivedMessage(msg)
			}
		}
	}

}

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods:                allowedMethods,
		AllowedHeaders:                allowedHeaders,
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	router, err := rest.MakeRouter(
		rest.Get("/webhook", GetHook),

		rest.Post("/webhook", PostHook),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9292"
	}
	fmt.Println(":"+port, "Listening....")
	log.Fatal(http.ListenAndServe(":"+port, api.MakeHandler()))
}
