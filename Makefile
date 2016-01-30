all: customCaddy

customCaddy: *.go
	caddyext install nlgids:github.com/miekg/nlgids
	caddyext stack
	caddyext build

.PHONY:
install: customCaddy
	cp -f customCaddy /opt/bin
	sudo setcap cap_net_bind_service=+ep /opt/bin/customCaddy

.PHONY:
clean:
	rm -f customCaddy
