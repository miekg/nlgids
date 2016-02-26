all: customCaddy

customCaddy: *.go
	caddyext install nlgids:github.com/miekg/nlgids
	caddyext stack
	caddyext build

.PHONY:
install: customCaddy
	cp -f customCaddy /opt/bin

.PHONY:
restart:
	sudo systemctl restart caddy

.PHONY:
clean:
	rm -f customCaddy
