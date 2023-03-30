package httpserver

import (
	"context"
	"crawler/app"
	"fmt"
	"net"
	"net/http"
)

func New(ctx context.Context, app *app.App, port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: newMux(app),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}
