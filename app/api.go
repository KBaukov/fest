package app

import (
	"context"
	"fest/config"
	"fest/utils"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"sync"
)

func ServApi(ctx context.Context, wg *sync.WaitGroup) { //ch chan struct{},

	// defer close(ch)

	settings, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Initializing App configuration error"))
	}

	server := fmt.Sprintf("%s:%d", settings.Server.Host, settings.Server.Port)

	router := mux.NewRouter()

	//router.Use(middleware.LoggingMiddleware)

	BuildCommonRoutes(router, "")

	router.NotFoundHandler = utils.NotFound
	router.MethodNotAllowedHandler = utils.NotAllowedMethod

	// Launch the app
	fmt.Printf("#################### Server listen: %s #############################\n", server)

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("#################### Server is Stopped ########################################")
			//close(ch)
			wg.Done()
			return
		}
	}()

	// or Log unsuccess start info
	if settings.Server.TLS {
		log.Fatal(http.ListenAndServeTLS(server, settings.Server.CertificatePath, settings.Server.KeyPath, router))
	} else {
		log.Fatal(http.ListenAndServe(server, router))
	}

}
