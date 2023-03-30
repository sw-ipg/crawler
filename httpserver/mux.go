package httpserver

import (
	"crawler/app"
	"crawler/httpserver/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func newMux(a *app.App) http.Handler {
	router := mux.NewRouter()

	router.Handle("/urls", handlers.AddUrlHandler(a)).Methods(http.MethodPost)

	return router
}
