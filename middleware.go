package nlgids

import (
	"fmt"
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
	// only POST are handled by this middleware
	if r.Method != "POST" {
		return h.Next.ServeHTTP(w, r)
	}

	path := middleware.Path(r.URL.Path)
	if !path.Matches("/api") {
		return h.Next.ServeHTTP(w, r)
	}
	//	r.ParseForm() // Required if you don't call r.FormValue()
	//      fmt.Println(r.Form["new_data"])

	switch {
	case path.Matches("/api/auth/test"):
		return WebInvoiceTest(w, r)

	case path.Matches("/api/auth/invoice"):
		name := r.PostFormValue("name")
		fmt.Fprintf(w, "auth! %s", name)
	case path.Matches("/api/open/contact"):
		fallthrough
	case path.Matches("/api/open/booking"):
		name := r.PostFormValue("name")
		fmt.Fprintf(w, "Hello, %s!", name)
	}
	return http.StatusOK, nil
}
