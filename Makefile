latest:
	docker build -t run-agent .
	docker tag run-agent pikacloud/run-agent:latest
	docker push pikacloud/run-agent:latest
