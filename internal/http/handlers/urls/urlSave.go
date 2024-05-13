package urls

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
	resp "url_shortener/internal/lib/api/response"
	"url_shortener/internal/lib/random"
)

const (
	aliasLength = 6
)

type URLSaver interface {
	SaveURL(SaveURL string, alias string) (*int64, error)
	SearchByURL(LoadURL string) (*string, error)
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func SaveURL(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.urls.urlSave"

		log = log.With(
			slog.String("op", op),
			slog.String("request_url", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			// Логируем к себе
			log.Error("failed to deserialize request", slog.String("error", err.Error()))

			// отдаем ответ пользователю
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(&req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, validationError(validateErr))
			return
		}

		matchAlias, err := urlSaver.SearchByURL(req.Url)
		if err != nil {
			log.Error("url not found", slog.String("error", err.Error()))
		}

		if matchAlias != nil {
			responseSuccess(w, r, *matchAlias)
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
			if err != nil {
				log.Error("failed to generate alias", slog.String("error", err.Error()))
			}
		}

		id, err := urlSaver.SaveURL(req.Url, alias)
		if err != nil {
			log.Error("failed to save url", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved", slog.Int64("id", *id))

		responseSuccess(w, r, alias)
	}
}

func responseSuccess(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.Success(),
		Alias:    alias,
	})
}

func validationError(errs validator.ValidationErrors) resp.Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return resp.Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
