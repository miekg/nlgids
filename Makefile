all: customCaddy

customCaddy:
	caddyext install nlgids:github.com/miekg/nlgids
	caddyext stack
	caddyext build

install: customCaddy setcap
	cp -f customCaddy /opt/bin

.PHONY:
setcap:
	sudo setcap cap_net_bind_service=+ep customCaddy

.PHONY:
clean:
	rm -f customCaddy
