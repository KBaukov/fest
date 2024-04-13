package app

import (
	common "fest/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Wrapper mux.MiddlewareFunc

// Wrap takes Handler functions and chains them to the main handler.
func Wrap(handler http.Handler, middlewares ...mux.MiddlewareFunc) http.Handler {
	// The loop is reversed so the adapters/middleware gets executed in the same
	// order as provided in the array.
	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}

	return handler
}

func BuildCommonRoutes(router *mux.Router, prefix string) {
	endpoint := fmt.Sprintf("%s", prefix)

	router.HandleFunc(endpoint+"/ws", common.ServeWs)

	router.HandleFunc(endpoint+"/api/user", common.ServeUserCreate).Methods("POST")
	//router.HandleFunc(endpoint+"/api/user/{token}", common.ServeUserCreate).Methods("POST")

	router.HandleFunc(endpoint+"/api/auth", common.ServeUserAuth).Methods("POST")

	router.HandleFunc(endpoint+"/", common.GetHome).Methods("GET")
	router.HandleFunc(endpoint+"/map", common.GetMap).Methods("GET")
	router.HandleFunc(endpoint, common.GetHome).Methods("GET")
	//// ############# get  web Resource ########################
	router.HandleFunc(endpoint+"/js/{file}", common.GetWebResources).Methods("GET")
	router.HandleFunc(endpoint+"/css/{file}", common.GetWebResources).Methods("GET")
	router.HandleFunc(endpoint+"/css/images/{file}", common.GetWebResources).Methods("GET")
	router.HandleFunc(endpoint+"/img/{file}", common.GetWebResources).Methods("GET")
	router.HandleFunc(endpoint+"/jquery-ui/{file}", common.GetWebResources).Methods("GET")

}
