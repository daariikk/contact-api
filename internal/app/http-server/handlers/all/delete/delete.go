package deleteAll

import (
	"contact-api/internal/app/http-server/common/server"
	"contact-api/internal/pkg/logger/sl"
	"fmt"
	"log/slog"
	"net/http"
)

type ContactsDeleter interface {
	DeleteAll() (int64, error)
}

// New создает обработчик HTTP для удаления всех контактов
// @Summary Удалить все контакты
// @Description Удаляет все контакты из базы данных
// @Tags contacts
// @Accept json
// @Produce json
// @Success 200 {string} string "Успешное удаление всех контактов"
// @Failure 500 {object} server.ErrorResponse "Ошибка сервера"
// @Router /v1/contact [delete]
func New(log *slog.Logger, deleter ContactsDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.all.delete.New"
		log = log.With(
			slog.String("op: ", op))

		count, err := deleter.DeleteAll()
		if err != nil {
			log.Info("error deleting all contacts: ", sl.Err(err))

			server.InternalError("error deleting records", err, w, r)

			return
		}

		resp := fmt.Sprintf("deleting %d records complete successfully", count)

		log.Info("deleting records complete successfully")

		server.RespondOK(resp, w, r)
	}
}
