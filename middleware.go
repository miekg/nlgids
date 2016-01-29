package nlgids

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddy/setup"
	"github.com/mholt/caddy/middleware"
)

const nlgids = "nlgids.london"

func Setup(c *setup.Controller) (middleware.Middleware, error) {
	return func(next middleware.Handler) middleware.Handler {
		return &handler{}
	}, nil
}

type handler struct {
	Next middleware.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// only POST from nlgids.london are handled by this middleware
	if r.Method != "POST" || !strings.HasPrefix(r.URL.Host, nlgids) {
		return h.Next.ServeHTTP(w, r)
	}
	// switch on path /api/path and call the correct function(s)

	r.ParseForm()

	name := r.PostFormValue("name")
	fmt.Fprintf(w, "Hello, %s!", name)
	return http.StatusOK, nil
}

/*
	w.Write([]byte("Hello, I'm a caddy middleware" + r.URL.Path))
	w.Write([]byte(r.Method))
*/
