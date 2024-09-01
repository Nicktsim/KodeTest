package login

import (
	"errors"
	"io"
	"net/http"

	resp "github.com/Nicktsim/kodetest/lib/api/response"
	"github.com/Nicktsim/kodetest/logger/sl"
	"github.com/Nicktsim/kodetest/storage/psql"
	"github.com/Nicktsim/kodetest/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	Login    string `json:"login" validate:"required,login"`
	Password string `json:"password" validate:"required,password"`
}

type Response struct {
	resp.Response
	Token string `json:"Token"`
}

func SignIn(log *slog.Logger, storage *psql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "users.login.SignIn"

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

		user, err := storage.SignIn(req.Login, req.Password)
		if err != nil {
			log.Error("wrong logn/password", sl.Err(err))

			render.JSON(w, r, resp.Error("wrong logn/password"))

			return
		}

		tokenString, err := utils.CreateToken(user)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to generate token"))

			return
		}

		log.Info("login succesful", slog.String("login", user.Login))

		responseOK(w, r, tokenString)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request, tokenString string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Token:    tokenString,
	})
}
