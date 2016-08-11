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
	pageAccessToken = "EAAESpoV8ZBZC0BAKAuxDj7JiZBjIqU3h05ZB0TCK40LiG7pUPfbSMKcVCKx2MnQwGgIEC1JlY2xUV775uCtZCPrnaDGXxVqZAKJZCbX2qX7GdXsyP8LJB1HlqHSIDszhCPvhj0ORTiTC3YloJxO0eJKzwfdRZAjuT0MZD"
	validationToken = "EAAESpoV8ZBZC0BANLkzrzFouB8sz8N0ZCUWV8ZCJQE6mMUZAqyBZAM6pPZCkVoZBlgdMIQZAFMvGOqw1MIuXvLIIdMvZA34GowLRYbgSvdiqgkGu1NmuDA9HnjlfW0ATb94ZApFCiFtD9oc7zDxaeL7eRek94GxXdsTnPkZD"
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
	if r.FormValue("hub[mode]") == "subscribe" && r.FormValue("hub[verify_token]") == validationToken {
		log.Println("Validating webhook")
		w.WriteJson(r.FormValue("hub[challenge]"))
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
				err = receivedMessage(entry.ID, msg)
			}
		}
	}

	if err != nil {
		w.WriteJson(map[string]string{"error": err.Error()})
		w.WriteHeader(500)
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

	api.Use(&Middleware{})

	router, err := rest.MakeRouter(
		rest.Get("/webhook/hub", GetHook),

		rest.Post("/webhook/", PostHook),
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
