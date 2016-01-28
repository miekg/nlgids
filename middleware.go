package nlgids

import (
	"net/http"

	"github.com/mholt/caddy/caddy/setup"
	"github.com/mholt/caddy/middleware"
)

func Setup(c *setup.Controller) (middleware.Middleware, error) {
	return func(next middleware.Handler) middleware.Handler {
		return &handler{}
	}, nil
}

type handler struct{
	Next middleware.Handler
//	h.Next.ServeHTTP(w, r)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	//Call to ParseForm makes form fields available.
    r.ParseForm()

    name := r.PostFormValue("name")
    fmt.Fprintf(w, "Hello, %s!", name)
}
	w.Write([]byte("Hello, I'm a caddy middleware" + r.URL.Path))
	w.Write([]byte(r.Method))
	return 200, nil
}
