package create

import (
	"errors"
	"io"
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

type Request struct {
    Title   	string `json:"title" validate:"required,title"`
    Description string `json:"description" validate:"required,description"`
}

type Response struct {
	resp.Response
	ID int `json:"id"`
}

func NewNote(log *slog.Logger, storage *psql.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.create.NewNote"

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
		var req Request
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

        if req.Title == "" {
            log.Error("empty title for note")

			render.JSON(w, r, resp.Error("empty title for note"))

			return
        }

        if req.Description == "" {
            log.Error("empty description for note", sl.Err(err))

			render.JSON(w, r, resp.Error("empty description for note"))

			return
        }

        if err := resp.ValidateNote(req.Title, req.Description); err != nil {
            log.Error("error during validation note", sl.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))

			return
        } else {
            log.Info("validation successful!")
        }

		id, err := storage.CreateNote(req.Title,req.Description, userID)
        if err != nil {
			log.Error("failed to add note", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add note"))

			return
		}
        log.Info("note added", slog.Int("id", id))

		responseOK(w, r, id)
    }
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		ID: id,
	})
}