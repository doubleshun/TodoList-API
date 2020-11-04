package routes

import (
	"net/http"

	controller "api/controller"

	"github.com/gorilla/mux"
)

//Route is ...
type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("POST", "/api/todo", controller.AddTodo, nil)
	register("DELETE", "/api/todo/{UserID}/{ItemID}", controller.DeleteTodoById, nil)
	register("PATCH", "/api/todo", controller.UpdateTodo, nil)
	register("Get", "/api/todo/{UserID}", controller.GetAllTodo, nil)
}

//NewRouter is ...
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		r.Methods(route.Method).
			Path(route.Pattern).
			Handler(route.Handler)
		if route.Middleware != nil {
			r.Use(route.Middleware)
		}
	}
	return r
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
