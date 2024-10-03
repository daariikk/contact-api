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
