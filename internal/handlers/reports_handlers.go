package handlers

import "net/http"

func LoadTMR(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello form LoadTMR"))
}
