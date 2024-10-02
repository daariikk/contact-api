package save

import (
	"contact-api/internal/app/domain/models"
	"contact-api/internal/app/http-server/common/server"
	"contact-api/internal/pkg/logger/sl"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ContactSaver interface {
	Save(contact models.Contact) (string, error)
}

type RespOK struct {
	ID  string `json:"id"`
	MSG string `json:"msg"`
}

// New создает обработчик HTTP для сохранения контакта
// @Summary Сохранить контакт
// @Description Сохраняет новый контакт в базе данных
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact body models.Contact true "Контакт для сохранения"
// @Success 200 {object} RespOK "Успешное сохранение контакта"
// @Failure 400 {object} server.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} server.ErrorResponse "Ошибка сервера"
// @Router /v1/contact [post]
func New(log *slog.Logger, saver ContactSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.all.save.New"
		log = log.With(
			slog.String("op: ", op))

		contact := models.Contact{}

		err := json.NewDecoder(r.Body).Decode(&contact)
		if err != nil {
			log.Info("error reading json", sl.Err(err))

			server.BadRequest("request error", err, w, r)

			return
		}

		id, err := saver.Save(contact)
		if err != nil {

			log.Info("error saving contact", sl.Err(err))

			server.InternalError("error saving contact", err, w, r)

			return
		}

		log.Info("saved contact successfully")

		server.RespondOK(RespOK{
			ID:  id,
			MSG: "successful save contact",
		}, w, r)
	}
}
