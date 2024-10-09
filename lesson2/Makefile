deploy:
	docker build -t go-web-app:1.2 .
	kind load docker-image go-web-app:1.2
	kubectl apply -f manifest

delete:
	kind delete cluster

test:
	docker run --rm -it --add-host host.docker.internal:host-gateway -v `pwd`:/bzt-configs blazemeter/taurus:latest test.yaml

setup:
	@make delete
	kind create cluster
	@make deploy

