package static

import "net/http"

func Static() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":3000", nil)
}

// func Static(w rest.ResponseWriter, r *rest.Request) {
// 	// hw := w.(http.ResponseWriter)
// 	// http.ServeFile(hw, r.Request, r.PathParam("file"))
//     http.FileServer(http.Dir("./"+r.PathParam("file")))
//     // w.WriteJson(`{"aa":"aaa"}`)
// }
