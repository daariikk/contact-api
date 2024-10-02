package update

import (
	"contact-api/internal/app/domain/models"
	"contact-api/internal/app/http-server/common/server"
	"contact-api/internal/app/storage"
	"contact-api/internal/pkg/logger/sl"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
)

type Updater interface {
	Update(contact models.Contact) (bool, error)
}

type Resp struct {
	OK  bool   `json:"ok"`
	MSG string `json:"msg"`
}

// New создает обработчик HTTP для обновления контакта
// @Summary Обновить контакт
// @Description Обновляет существующий контакт в базе данных
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact body models.Contact true "Контакт для обновления"
// @Success 200 {object} Resp "Успешное обновление контакта"
// @Failure 400 {object} server.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} server.ErrorResponse "Ошибка сервера"
// @Router /v1/contact/{uid} [put]
func New(log *slog.Logger, updater Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.one.update.New"
		log = log.With(
			slog.String("op: ", op))

		var contact = models.Contact{}

		err := json.NewDecoder(r.Body).Decode(&contact)
		if err != nil {
			log.Info("error parsing request body", sl.Err(err))

			server.InternalError("error parsing request body", err, w, r)

			return
		}

		uid := chi.URLParam(r, "uid")

		contact.ID = uid

		log.Info("request body parsing complete successfully")

		res, err := updater.Update(contact)
		if err != nil {
			if err != nil && errors.Is(err, storage.ErrContactNotFound) {
				log.Info("contact not found", slog.String("id", uid))
				server.BadRequest("contact not found", err, w, r)
				return
			}

			log.Info("error updating contact", sl.Err(err))

			server.InternalError("error updating contact", err, w, r)

			return
		}

		log.Info("update item complete successfully")

		server.RespondOK(Resp{
			OK:  res,
			MSG: "successful update item with id: " + contact.ID,
		}, w, r)

	}
}
