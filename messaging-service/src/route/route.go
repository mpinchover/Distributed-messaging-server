package route

import (
	"encoding/json"
	"messaging-service/src/middleware"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, contentType, Content-Type, Accept, Authorization")
}

type RootHandler struct {
	handler    middleware.HTTPHandler
	middleware []middleware.Middleware
}

func New(handler middleware.HTTPHandler, middleware []middleware.Middleware) RootHandler {
	wrappedHandler := InitializeHandler(handler, middleware)
	return RootHandler{
		middleware: middleware,
		handler:    wrappedHandler,
	}
}

func (rh RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO - if in dev mode
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	res, err := rh.handler(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func InitializeHandler(handler middleware.HTTPHandler, middleware []middleware.Middleware) middleware.HTTPHandler {
	if len(middleware) < 1 {
		return handler
	}

	wrapped := handler

	// loop in reverse to preserve middleware order
	for i := len(middleware) - 1; i >= 0; i-- {
		wrapped = middleware[i].Execute(wrapped)
	}

	return wrapped
}
