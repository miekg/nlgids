all: Caddy

.PHONY: Caddy
Caddy: *.go
	@echo "edit caddy/caddymain/run.go and import this 'nlgids' middleware"
	@echo "edit caddyhttp/httpserver/plugin.go and add 'nlgids' to the directives"

.PHONY: restart
restart:
	sudo systemctl restart caddy
