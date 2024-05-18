package delete

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	resp "simple_RESTapi/internal/lib/api/response"
	"simple_RESTapi/internal/lib/logger/sl"
	"simple_RESTapi/internal/storage"
)

// URLDelete is an interface for delete url by alias.
//

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=URLDelete
type URLDelete interface {
	URLDelete(alias string) error
}

func New(log *slog.Logger, URLDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		// Лучше больше логов, чем меньше - лишнее мы легко сможем почистить,
		// при необходимости. А вот недостающую информацию мы уже не получим.
		log.Info("request body decoded", slog.Any("req", req))

		// ...
		alias := req.Alias

		// Пробуем удалить строку по alias
		err = URLDelete.URLDelete(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			// Не нашли URL, сообщаем об этом клиенту
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("delete alias", slog.String("alias", alias))

		responseOK(w, r, alias)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Status: resp.OK().Status,
		Error:  "",
	})
}
