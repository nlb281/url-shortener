package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	res "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(slog.String("op", op))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, res.Error("not found"))

			return
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, res.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", err)

			render.JSON(w, r, res.Error("internal error"))

			return
		}
		log.Info("got url", slog.String("url", url))

		// redirect
		http.Redirect(w, r, url, http.StatusFound)
	}
}