package get

import (
	"net/http"
	"strings"

	resp "github.com/Nicktsim/kodetest/lib/api/response"
	"github.com/Nicktsim/kodetest/logger/sl"
	"github.com/Nicktsim/kodetest/storage/psql"
	"github.com/Nicktsim/kodetest/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	resp.Response
	Notes []psql.Note
}

func GetUserNotes(log *slog.Logger, storage *psql.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.create.GetUserNotes"

        authorizationHeader := r.Header.Get("Authorization")
        if authorizationHeader == "" {
            log.Error("failed to find token")

			render.JSON(w, r, resp.Error("failed to find token"))

			return
        }

        tokenString := strings.Split(authorizationHeader, " ")[1]

        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            log.Error("token validation error", sl.Err(err))

			render.JSON(w, r, resp.Error("token validation error"))

			return
        }

        userID := claims.UserID
        log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		notes, err := storage.GetNotes(userID)
        if err != nil {
			log.Error("failed to get notes list", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get notes list"))

			return
		}
        log.Info("successfully got notes")

		responseOK(w, r, notes)
    }
}

func responseOK(w http.ResponseWriter, r *http.Request, notes []psql.Note) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Notes: notes,
	})
}