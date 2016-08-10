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
	pageAccessToken = "EAACEdEose0cBAJ3W5QcN24OJJpl37V6sQDmAmPIXUm7Rp8MoMZBQIAQEeDzcMWapGthOEfHApglrAH1q4cRYrWwhHD4EOFWEqO3Bl2cXJZC7bqllrPCcik9nxAc86aKbM1gNhq7eQQTVwiFDXBYadpcB6LUr7TlMJDrjSoeAZDZD"
	validationToken = "EAACEdEose0cBAG8N2QEOGPAUp29vQ5s1XpuIJcusGD4tcXoZBlXmEK78lSaVy8fRMQZAvyQdMDjbsGzEKZBUPNJ3EbT388kDS6HXHNlNS8Nn0VNNvof8EGZAM6n1lXkn9KcTlAPTeRn4UynEf5IqIOG2zEd8cO9ksU0jsjy2aQZDZD"
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
