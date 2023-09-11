package middleware

// import (
// 	"encoding/json"
// 	"net/http"
// )

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, contentType, Content-Type, Accept, Authorization")
// }

// // type Middleware func(HTTPHandler) HTTPHandler
// type HTTPHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

// // type HTTPHandlerV2 func(interface{}) (interface{}, error)

// // type HandlerFunc func(ResponseWriter, *Request)
// type RootHandler struct {
// 	handler HTTPHandler
// 	// handlerV2  HTTPHandlerV2
// 	middleware []Middleware
// }

// // maybe just do it here
// func New(handler HTTPHandler, middleware []Middleware) RootHandler {
// 	wrappedHandler := InitializeHandler(handler, middleware)

// 	// wrappedHandler := InitializeHandler(middleware)
// 	return RootHandler{
// 		middleware: middleware,
// 		handler:    wrappedHandler,
// 	}
// }

// // func TestMiddleware1(h HTTPHandler) HTTPHandler {
// // 	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// // 		ctx := context.WithValue(r.Context(), "Username1", "Bob Moses")
// // 		fmt.Println("EXECUTING M1 HERE")
// // 		r = r.WithContext(ctx)
// // 		h(w, r)
// // 		return nil, nil
// // 	}
// // }

// // func TestMiddleware2(h HTTPHandler) HTTPHandler {
// // 	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// // 		ctx := context.WithValue(r.Context(), "Username2", "Bobbie Moses")
// // 		fmt.Println("EXECUTING M2 HERE")
// // 		r = r.WithContext(ctx)
// // 		h(w, r)
// // 		return nil, nil
// // 	}
// // }

// // func (rh RootHandler) Handler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// // 	// rh.
// // }

// func (rh RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w)
// 	w.Header().Set("Content-Type", "application/json")

// 	// do it here
// 	res, err := rh.handler(w, r)
// 	// res, err := rh.Handler(w, r)
// 	if err != nil {
// 		// err = errors.Wrap(err, 100)
// 		// err := serrors.GetStackTrace(err)
// 		// log.Println(errors.Unwrap(err))
// 		// stack := serrors.GetStackTrace(err)

// 		// log.Println(stack.(*errors.Error).ErrorStack())
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = json.NewEncoder(w).Encode(res)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
// }

// // return the handler/fn, not the return types of the fn
// // func InitializeHandler(handler HTTPHandler, middleware []Middleware) HTTPHandler {
// // 	if len(middleware) < 1 {
// // 		return handler
// // 	}

// // 	wrapped := handler

// // 	// loop in reverse to preserve middleware order
// // 	for i := len(middleware) - 1; i >= 0; i-- {
// // 		wrapped = middleware[i](wrapped)
// // 	}

// // 	return wrapped
// // }

// func InitializeHandler(handler HTTPHandler, middleware []Middleware) HTTPHandler {
// 	if len(middleware) < 1 {
// 		return handler
// 	}

// 	wrapped := handler

// 	// loop in reverse to preserve middleware order
// 	for i := len(middleware) - 1; i >= 0; i-- {
// 		wrapped = middleware[i].execute(wrapped)
// 	}

// 	return wrapped
// }

// // func (mw MWType) execute(h HTTPHandler) HTTPHandler {
// // 	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
// // 		ctx := context.WithValue(r.Context(), "Username1", "Bob Moses")
// // 		fmt.Println("EXECUTING M1 HERE")
// // 		r = r.WithContext(ctx)
// // 		h(w, r)
// // 		return nil, nil
// // 	}
// // }
