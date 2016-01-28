
#caddyext stack

customCaddy:
	caddyext build
	sudo setcap cap_net_bind_service=+ep customCaddy  
