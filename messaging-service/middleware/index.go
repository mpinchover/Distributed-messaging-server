package middleware

type Middleware interface {
	execute(HTTPHandler) HTTPHandler
}
