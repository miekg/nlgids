package nlgids

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("nlgids", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// Config holds the Caddy file directives.
// Typically these will look like:
//
// nlgids {
//	   recipients ans@nlgids.london miek@miek.nl
//	   subject ans@nlgids.london
//	   secret /etc/caddy/NLgids-fcbeb7928cdb.json
//	   template /etc/caddy/tmpl
//	   tours /var/www/nlgids.london/tours.json
// }
type Config struct {
	Recipients []string // Who gets nlgids email.
	Subject    string   // Calendar auth subject.
	Secret     string   // File containing the google service account secret.
	Template   string   // Directory where the templates live.
	Tours      string   // tours.json location, defaults to /var/www/nlgids.london/tours.json
}

func setup(c *caddy.Controller) error {
	config, err := nlgidsParse(c)
	if err != nil {
		return err
	}
	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return &NLgids{Next: next, Config: config}
	})

	return nil
}

// nlgidsParse will parse the following directives.
// recipients ans@nlgids.london miek@miek.nl
// subject ans@nlgids.london
// secret /etc/caddy/NLgids-fcbeb7928cdb.json
func nlgidsParse(c *caddy.Controller) (*Config, error) {
	config := new(Config)
	config.Tours = "/var/www/nlgids.london/tours.json"
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "recipients":
				rcpts := c.RemainingArgs()
				if len(rcpts) == 0 {
					return nil, c.ArgErr()
				}
				config.Recipients = append(config.Recipients, rcpts...)
			case "subject":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				config.Subject = c.Val()
				if !strings.Contains(config.Subject, "@") {
					return nil, fmt.Errorf("nlgids: subject must contain @-sign: %s", c.Val())
				}
			case "secret":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				config.Secret = c.Val()
				_, err := os.Open(config.Secret)
				if err != nil {
					return nil, fmt.Errorf("nlgids: secret file must be readable: %s", err)
				}
			case "template":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				config.Template = c.Val()
			case "tours":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				config.Tours = c.Val()
			}
		}
	}
	return config, nil
}

// NLgids is the NLgids handler.
type NLgids struct {
	Next   httpserver.Handler
	Config *Config
}

func (n *NLgids) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "POST" {
		return n.Next.ServeHTTP(w, r)
	}
	path := httpserver.Path(r.URL.Path)
	if !path.Matches("/api") {
		return n.Next.ServeHTTP(w, r)
	}

	switch {
	case path.Matches("/api/auth/invoice"):
		return n.WebInvoice(w, r)
	case path.Matches("/api/auth/conform"):
		return n.WebConform(w, r)
	case path.Matches("/api/open/contact"):
		return n.WebContact(w, r)
	case path.Matches("/api/open/booking"):
		return n.WebBooking(w, r)
	case path.Matches("/api/open/calendar"):
		return n.WebCalendar(w, r)
	}
	return http.StatusOK, nil
}
