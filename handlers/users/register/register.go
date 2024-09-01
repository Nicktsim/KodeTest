package register

import (
	"errors"
	"io"
	"net/http"

	resp "github.com/Nicktsim/kodetest/lib/api/response"
	"github.com/Nicktsim/kodetest/logger/sl"
	"github.com/Nicktsim/kodetest/storage/psql"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
    Login   	string `json:"login" validate:"required,login"`
    Password 	string `json:"password" validate:"required,password"`
	Username 	string `json:"username" validate:"required,username"`
}

type Response struct {
    resp.Response
    Username  string `json:"Username"`
}


func SignUp(log *slog.Logger, storage *psql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.register.SignUp"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := storage.SignUp(req.Login,req.Password,req.Username)
		if errors.Is(err, psql.ErrUserExists) {
			log.Info("user already exists", slog.String("login", req.Login))

			render.JSON(w, r, resp.Error("login already exists"))

			return
		}
		if err != nil {
			log.Error("failed to sign up user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to sign up user"))

			return
		}

		log.Info("User added", slog.Int("id", id))

		responseOK(w, r, req.Username)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, username string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Username:    username,
	})
}