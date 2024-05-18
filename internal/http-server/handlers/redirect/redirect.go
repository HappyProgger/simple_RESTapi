package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	resp "simple_RESTapi/internal/lib/api/response"
	"simple_RESTapi/internal/lib/logger/sl"
	"simple_RESTapi/internal/storage"
)

// URLRedirecter is an interface for getting url by alias.
//

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=URLRedirecter
type URLRedirecter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, URLRedirecter URLRedirecter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Роутер chi позволяет делать вот такие финты -
		// получать GET-параметры по их именам.
		// Имена определяются при добавлении хэндлера в роутер, это будет ниже.
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		// Находим URL по алиасу в БД
		resURL, err := URLRedirecter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			// Не нашли URL, сообщаем об этом клиенту
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		// Делаем редирект на найденный URL
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
