package deleteOne

import (
	"contact-api/internal/app/http-server/common/server"
	"contact-api/internal/app/storage"
	"contact-api/internal/pkg/logger/sl"
	"errors"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
)

type DeleterByID interface {
	Delete(id string) (bool, error)
}

type Resp struct {
	OK  bool   `json:"ok"`
	MSG string `json:"msg"`
}

// New создает обработчик HTTP для удаления контакта по ID
// @Summary Удалить контакт по ID
// @Description Удаляет контакт из базы данных по указанному ID
// @Tags contacts
// @Accept json
// @Produce json
// @Param id body string true "ID контакта для удаления"
// @Success 200 {object} Resp "Успешное удаление контакта"
// @Failure 400 {object} server.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} server.ErrorResponse "Ошибка сервера"
// @Router /v1/contact/{uid} [delete]
func New(log *slog.Logger, deleter DeleterByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.one.delete.New"
		log = log.With(
			slog.String("op: ", op))

		uid := chi.URLParam(r, "uid")

		if uid == "" {
			log.Info("empty id", slog.String("id", uid))

			server.BadRequest("uncorrected url, id is empty", nil, w, r)

			return
		}

		res, err := deleter.Delete(uid)
		if err != nil {
			if err != nil && errors.Is(err, storage.ErrContactNotFound) {
				log.Info("contact not found", slog.String("id", uid))
				server.BadRequest("contact not found", err, w, r)
				return
			}

			log.Info("error deleting item", sl.Err(err))
			server.InternalError("error deleting item", err, w, r)
			return
		}

		server.RespondOK(Resp{
			OK:  res,
			MSG: "complete deleting item with id: " + uid,
		}, w, r)

	}
}
