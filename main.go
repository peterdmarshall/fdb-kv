package main

import (
	"errors"
	"github.com/apple/foundationdb/bindings/go/src/fdb"

	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

type GetKeyResponse struct {
	Body struct {
		Key   string `json:"key" example:"foo"`
		Value string `json:"value" example:"bar"`
	}
}

type SetKeyResponse struct {
	Body struct {
		Result string `json:"result" example:"ok"`
	}
}

type SetKeyRequestBody struct {
	Value string `json:"value"`
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		router.Use(middleware.Logger)
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Different API versions may expose different runtime behaviors.
		fdb.MustAPIVersion(710)

		// Open the default database from the system cluster
		db := fdb.MustOpenDefault()

		huma.Put(api, "/{key}", func(ctx context.Context, input *struct {
			Key  string            `path:"key"`
			Body SetKeyRequestBody `contentType:"application/json"`
		}) (*SetKeyResponse, error) {
			resp := &SetKeyResponse{}
			_, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				tr.Set(fdb.Key(input.Key), []byte(input.Body.Value))
				return nil, nil
			})
			if err != nil {
				resp.Body.Result = "err"
				return resp, errors.New("failed to set key")
			}
			resp.Body.Result = "ok"
			return resp, nil
		})

		huma.Get(api, "/{key}", func(ctx context.Context, input *struct {
			Key string `path:"key"`
		}) (*GetKeyResponse, error) {
			resp := &GetKeyResponse{}
			ret, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				return tr.Get(fdb.Key(input.Key)).MustGet(), nil
			})
			if err != nil {
				return nil, errors.New("failed to get key")
			}
			resp.Body.Key = input.Key
			resp.Body.Value = string(ret.([]byte))
			return resp, nil
		})

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			fmt.Println("Starting server on port", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
