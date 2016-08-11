package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"goface/static"

	"gopkg.in/yaml.v2"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	appSecret       string
	pageAccessToken string
	validationToken string
	allowedMethods  = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowedHeaders  = []string{
		"Accept",
		"Authorization",
		"X-Real-IP",
		"Content-Type",
		"X-Custom-Header",
		"Query",
		"Language",
		"Origin",
	}
	pipe map[string]chan string
)

func init() {
	pipe = make(map[string]chan string)
}

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
	pipe["1"] = make(chan string)
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
	app := config()
	appSecret = app.AppSecret
	pageAccessToken = app.PageAccessToken
	validationToken = app.ValidationToken

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

	// api.Use(&Middleware{})

	go static.Static()

	router, err := rest.MakeRouter(
		rest.Get("/events", SendSSE),
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

type facebookApp struct {
	AppSecret       string `yaml:"appSecret"`
	PageAccessToken string `yaml:"pageAccessToken"`
	ValidationToken string `yaml:"validationToken"`
}

func config() facebookApp {
	var b []byte
	var err error
	app := facebookApp{
		AppSecret:       os.Getenv("appSecret"),
		PageAccessToken: os.Getenv("pageAccessToken"),
		ValidationToken: os.Getenv("validationToken"),
	}
	if app.PageAccessToken != "" {
		return app
	}

	if b, err = ioutil.ReadFile("./config.yaml"); err != nil {
		log.Fatal("you need yaml config file.")
	}

	if err = yaml.Unmarshal(b, &app); err != nil {
		log.Fatal("Please check your config.yaml file.")
	}

	return app
}
