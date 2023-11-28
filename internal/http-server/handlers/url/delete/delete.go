package delete

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlGetter URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		err := urlGetter.DeleteURL(alias)
		if err != nil {
			log.Info("failed to delete URL by alias")

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("delete URL by alias", slog.String("alias", alias))
	}
}
