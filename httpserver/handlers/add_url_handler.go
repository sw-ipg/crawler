package handlers

import (
	"crawler/app"
	"fmt"
	"io"
	"net/http"
)

func AddUrlHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("cannot read body: %s", err)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = app.AddUrl(r.Context(), string(url)); err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("cannot add you url: %s", err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte("url successfully added!"))
	}
}
