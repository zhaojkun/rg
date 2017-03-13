package example

import "net/http"

// index get /index
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index"))

}
