//go:generate swagger generate spec
package main

import (
	"api/controller"
	_ "api/docs"
	routes "api/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zabawaba99/firego"
)

//Fb is ...
var Fb firego.Firebase

func init() {
	Fb := firego.New("https://zdl-todolist.firebaseio.com/", nil)
	controller.Fb = *Fb //access global variable from main
}

// @title testProject API.
// @version 0.1
// @description just test golang and swagger for fun

// @contact.name Doubleshun
// @contact.url https://github.com/doubleshun

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// schemes http
func main() {
	router := routes.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err) //throw error
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		panic(http.ListenAndServe(":8080", router))
	}

}
