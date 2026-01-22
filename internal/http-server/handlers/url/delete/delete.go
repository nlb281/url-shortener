package delete

import (
	"errors"
	"log/slog"
	"net/http"
	res "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Response struct {
	res.Response
	Alias string `json:"alias,omitempty"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(slog.String("op", op))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, res.Error("not found"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotDeleted) {
			log.Info("failed to delete url", "alias", alias)

			render.JSON(w, r, res.Error("internal error"))

			return
		}
		log.Info("URL successfully deleted")

		render.JSON(w, r, Response{
			Response: res.Delete(),
			Alias: alias,
		})
	}
}