package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	//"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"

	"github.com/cyberbono3/infura-coding/client-api/data"
	"github.com/cyberbono3/infura-coding/client-api/env"
	"github.com/cyberbono3/infura-coding/client-api/handlers"
	"github.com/cyberbono3/infura-coding/config"
	"github.com/cyberbono3/infura-coding/rpc"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind address for the server")

func main() {

	env.Parse()

	cfg := config.NewConfig(config.Mainnet, "10dea4f9680f440bb9eaf33ecc51c1ef")
	c := rpc.New(cfg.GetURL())
	l := hclog.Default()
	v := data.NewValidation()

	clientHandler := handlers.NewClientHandler(c, l, v)


	sm := mux.NewRouter()


	transR := sm.Methods(http.MethodPost).Subrouter()
	transR.HandleFunc("/transaction", clientHandler.GetTransaction)
	transR.Use(clientHandler.MiddlewareValidateTransactionData)

	blockR := sm.Methods(http.MethodPost).Subrouter()
	blockR.HandleFunc("/block", clientHandler.GetBlock)
	blockR.Use(clientHandler.MiddlewareValidateBlockData)

	s := http.Server{
		Addr:         *bindAddress,                                     // configure the bind address
		Handler:      sm,                                               // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	// Block until a signal is received.
	sig := <-ch
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
