package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
)

type Middleware struct {
}

func (m *Middleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		signature := r.Header.Get("x-hub-signature")

		if signature == "" {
			w.WriteJson(map[string]string{"error": "Couldn't validate the signature."})
			return
		}

		elements := strings.Split(signature, "=")
		// method := elements[0]
		signatureHash := elements[1]

		mac := hmac.New(sha256.New, nil)
		mac.Write([]byte(appSecret))
		expectedMAC := mac.Sum(nil)

		fmt.Println("signatureHash", signatureHash)
		fmt.Println("expectedMAC", hex.EncodeToString(expectedMAC))

		// var expectedHash = crypto.createHmac('sha1', APP_SECRET)
		//     .update(buf)
		//     .digest('hex');
		//
		// if (signatureHash != expectedHash) {
		//     throw new Error("Couldn't validate the request signature.");
		// }

		handler(w, r)
	}
}
