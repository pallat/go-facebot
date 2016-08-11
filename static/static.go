package static

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func Static(w rest.ResponseWriter, r *rest.Request) {
	// hw := w.(http.ResponseWriter)
	// http.ServeFile(hw, r.Request, r.PathParam("file"))
	http.FileServer(http.Dir("./" + r.PathParam("file")))
	// w.WriteJson(`{"aa":"aaa"}`)
}
