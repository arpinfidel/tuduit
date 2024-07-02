package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/repo/example"
	taskrepo "github.com/arpinfidel/tuduit/repo/task"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	deps Dependencies
}

type Dependencies struct {
	ExampleRepo example.Repo
	TaskRepo    taskrepo.Repo

	DB *db.DB // TODO: temp
}

func ParseJSON[P any, Q any, Req any, Res any](f func(context.Context, P, Q, Req) (Res, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pm := map[string]any{}
		if rctx := chi.RouteContext(r.Context()); rctx != nil {
			for i, k := range rctx.URLParams.Keys {
				v := rctx.URLParams.Values[i]
				pm[k] = v
			}
		}
		pj, err := json.Marshal(pm)
		if err != nil {
			log.Printf("[error] %s\n", err.Error())
		}

		qs := r.URL.Query()
		qm := map[string]any{}
		for k, v := range qs {
			if len(v) > 1 {
				qm[k] = v
				continue
			}
			qm[k] = v[0]
		}
		qj, _ := json.Marshal(qm)

		var p P
		err = json.Unmarshal(pj, &p)
		if err != nil {
			log.Printf("[error] %s\n", err.Error())
		}

		var q Q
		err = json.Unmarshal(qj, &q)
		if err != nil {
			log.Printf("[error] %s\n", err.Error())
		}

		var req Req
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("[error] %s\n", err.Error())
		}

		ctx := r.Context()
		res, err := f(ctx, p, q, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := map[string]string{
				"error": err.Error(),
			}
			json.NewEncoder(w).Encode(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
