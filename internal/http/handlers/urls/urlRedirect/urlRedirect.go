package urlRedirect

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
	resp "url_shortener/internal/lib/api/response"
)

type GetURL interface {
	GetURL(alias string) (string, error)
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func RedirectUrl(log *slog.Logger, getUrl GetURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := strings.TrimLeft(r.RequestURI, "/")

		urlForRedirect, err := getUrl.GetURL(alias)

		if err != nil {
			log.Error("Failed to get url", slog.String("error", err.Error()))
			render.JSON(w, r, resp.Response{
				Status: StatusError,
				Error:  "Failed to get url",
			})
			return
		}

		http.Redirect(w, r, urlForRedirect, http.StatusTemporaryRedirect)
	}
}
