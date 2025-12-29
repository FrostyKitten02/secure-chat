package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"secure-chat/api"
	"secure-chat/env"
	"secure-chat/logs"
	"secure-chat/repo"
	"secure-chat/ws"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func main() {
	logs.InitLogger()

	cli := humacli.New(func(hooks humacli.Hooks, options *env.Options) {
		dbErr := repo.InitDB(options.DbConnectionString)
		if dbErr != nil {
			slog.Error("Failed to init DB", "error", dbErr)
			panic(dbErr)
		}
		// Create a new router & API
		router := chi.NewMux()
		apiHuma := humachi.New(router, huma.DefaultConfig("Secure chat API", "1.0.0"))

		router.Handle("/ws", http.HandlerFunc(ws.Handler))
		api.RegisterAll(apiHuma, options)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			runningOn := fmt.Sprintf("127.0.0.1:%d", options.HttpPort)
			slog.Info("Server starting on http://" + runningOn)
			http.ListenAndServe(runningOn, router)
		})
	})

	cli.Run()
}
