package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HelloHandler returns a greeting. If a 'name' query parameter is provided, it greets the given name.
func HelloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		greeting := "Hello, World!"
		if name != "" {
			greeting = fmt.Sprintf("Hello, %s!", name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": greeting,
		})
	}
}
