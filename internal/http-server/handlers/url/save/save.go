package save

import (
	"errors"
	"log/slog"
	"net/http"
	res "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL 	string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	res.Response
	Alias string `json:alias,omitempty`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// TODO: move to config
const (
	aliasLength = 8
	maxAttempts = 5
)

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(slog.String("op", op))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, res.Error("failed to decode request"))

			return
		}
		log.Info("request body decode", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", err)

			render.JSON(w, r, res.Error("Invalid request"))

			return
		}

		var idx int64
		alias := req.Alias
		for i := 0; i < maxAttempts; i++{
			if alias == "" {
				alias, err = random.RandomAlias(aliasLength)
				if err != nil {
					log.Error("failed to generate alias", err)
					render.JSON(w, r, res.Error("internal error"))
					return
				}
			} 

			id, err := urlSaver.SaveURL(req.URL, alias)
			if err == nil {
				idx = id
				break
			}

			if errors.Is(err, storage.ErrAliasExists) {
				if req.Alias != "" {
						render.Status(r, http.StatusConflict)
						render.JSON(w, r, res.Error("alias already exists"))
						return
				}

				alias = ""
				continue
			}
			log.Error("failed to save url", err)
			render.JSON(w, r, res.Error("internal error"))
			return
		}
		
		log.Info("url added", slog.Int64("id", idx))

		render.JSON(w, r, Response{
			Response: res.OK(),
			Alias: 		alias,
		})
	}
}