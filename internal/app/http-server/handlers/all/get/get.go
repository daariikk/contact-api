package getAll

import (
	"contact-api/internal/app/domain/models"
	"contact-api/internal/app/http-server/common/server"
	"contact-api/internal/pkg/logger/sl"
	"log/slog"
	"net/http"
)

type ContactsAll interface {
	GetAll() ([]models.Contact, error)
}

// New создает обработчик HTTP для получения всех контактов
// @Summary Получить все контакты
// @Description Возвращает список всех контактов в формате JSON
// @Tags contacts
// @Accept json
// @Produce json
// @Success 200 {array} models.Contact "Успешно получен список контактов"
// @Failure 500 {object} server.ErrorResponse "Ошибка сервера"
// @Router /v1/contact [get]
func New(log *slog.Logger, getAller ContactsAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.all.get.New"
		log = log.With(
			slog.String("op:", op))

		contacts, err := getAller.GetAll()
		if err != nil {
			log.Info("error getting lines", sl.Err(err))

			server.InternalError("error get any record", err, w, r)

			return
		}

		log.Info("successfully getting all records")

		server.RespondOK(contacts, w, r)
	}
}
