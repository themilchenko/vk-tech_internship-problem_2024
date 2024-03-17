package app

import (
	"fmt"
	"net/http"
)

type Route struct {
	Name               string
	Method             string
	Pattern            string
	isAuth             bool
	isAccessRestricted bool
	HandlerFunc        http.HandlerFunc
}

type Routes []Route

func NewRouter(authMiddleware func(http.HandlerFunc) http.HandlerFunc, routes Routes) http.Handler {
	mux := http.NewServeMux()
	for _, route := range routes {
		var handler http.HandlerFunc
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		if authMiddleware != nil && route.isAuth {
			handler = authMiddleware(handler)
		}

		mux.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != route.Method {
				http.Error(
					w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
	return mux
}

func (s *Server) getRoutes() Routes {
	return Routes{
		Route{
			"SignUp",
			"POST",
			"/signup",
			false,
			false,
			s.authHandler.Signup,
		},
		Route{
			"Login",
			"POST",
			"/login",
			false,
			false,
			s.authHandler.Login,
		},
		Route{
			"Auth",
			"GET",
			"/auth",
			true,
			false,
			s.authHandler.Auth,
		},
		Route{
			"Logout",
			"DELETE",
			"/logout",
			true,
			false,
			s.authHandler.Logout,
		},
	}
}

func Logger(inner http.Handler, name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Executing %s\n", name)
		inner.ServeHTTP(w, r)
	})
}
