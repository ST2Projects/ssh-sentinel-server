install:
	cp -a samples/config.json .
	mkdir -p resources
	cp samples/ssh-sentinel.service /etc/systemd/system/ssh-sentinel.service
	systemctl daemon-reload

	echo "Config file copied to ${PWD}/config.json. Please edit it before starting the service"
	echo "To start the service run systemctl start ssh-sentinel - To enable start on boot: systemctl enable ssh-sentinel"
