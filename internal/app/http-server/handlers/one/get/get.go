package getOne

import (
	"contact-api/internal/app/domain/models"
	"contact-api/internal/app/http-server/common/server"
	"contact-api/internal/app/storage"
	"contact-api/internal/pkg/logger/sl"
	"errors"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
)

type GetterByID interface {
	ContactById(id string) (models.Contact, error)
}

// New создает обработчик HTTP для получения контакта по ID
// @Summary Получить контакт по ID
// @Description Возвращает контакт из базы данных по указанному ID
// @Tags contacts
// @Accept json
// @Produce json
// @Param id body string true "ID контакта для получения"
// @Success 200 {object} models.Contact "Успешное получение контакта"
// @Failure 400 {object} server.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} server.ErrorResponse "Ошибка сервера"
// @Router /v1/contact/{uid} [get]
func New(log *slog.Logger, getter GetterByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.one.getOne.New"
		log = log.With(
			slog.String("op: ", op))

		uid := chi.URLParam(r, "uid")

		if uid == "" {
			log.Info("empty id", slog.String("id", uid))

			server.BadRequest("uncorrected uri, id is empty", nil, w, r)

			return
		}

		res, err := getter.ContactById(uid)
		if err != nil && errors.Is(err, storage.ErrContactNotFound) {
			log.Info("contact not found", slog.String("id", uid))
			server.BadRequest("contact not found", err, w, r)
			return
		}
		if err != nil {
			log.Info("error getting item", sl.Err(err))

			server.InternalError("error getting item", err, w, r)

			return
		}

		log.Info("get contact by ID complete successful")

		server.RespondOK(res, w, r)

	}
}
