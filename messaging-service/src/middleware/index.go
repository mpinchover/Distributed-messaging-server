package middleware

import "net/http"

type Middleware interface {
	Execute(HTTPHandler) HTTPHandler
}

type HTTPHandler func(http.ResponseWriter, *http.Request) (interface{}, error)
