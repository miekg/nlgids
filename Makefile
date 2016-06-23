all: Caddy

.PHONY: Caddy
Caddy: *.go
	@echo "edit caddy/caddymain/run.go and import this 'nlgids' middleware"

.PHONY: restart
restart:
	sudo systemctl restart caddy
