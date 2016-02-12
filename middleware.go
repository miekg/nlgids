package nlgids

import (
	"net/http"

	"github.com/mholt/caddy/caddy/setup"
	"github.com/mholt/caddy/middleware"
)

const nlgids = "nlgids.london"

func Setup(c *setup.Controller) (middleware.Middleware, error) {
	return func(next middleware.Handler) middleware.Handler {
		return &handler{Next: next}
	}, nil
}

type handler struct {
	Next middleware.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	println("allo")
	if r.Method != "POST" {
		return h.Next.ServeHTTP(w, r)
	}

	path := middleware.Path(r.URL.Path)
	if !path.Matches("/api") {
		return h.Next.ServeHTTP(w, r)
	}

	switch {
	case path.Matches("/api/auth/test"):
		return WebInvoiceTest(w, r)

	case path.Matches("/api/auth/invoice"):

	case path.Matches("/api/open/contact"):
		return WebContact(w, r)

	case path.Matches("/api/open/booking"):

	case path.Matches("/api/open/calendar"):
	case path.Matches("/api/open/calendar/test"):
		return WebCalendarTest(w, r)

	}
	return http.StatusOK, nil
}
