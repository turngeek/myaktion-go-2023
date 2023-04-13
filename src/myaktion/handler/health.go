package handler

import (
	"io"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := io.WriteString(w, `{"alive": true}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
