CLUSTER_NAME=port-mapping-cluster
TAG=0.1
CONFIG_DIR=configs
MANIFEST_DIR=k8s
SRC_DIR=src

create_cluster:
	kind create cluster --name $(CLUSTER_NAME) --config ./kind-portmapping-config.yaml

create_configmap:
	kubectl create configmap myconfig --from-file=app=./$(CONFIG_DIR)/my-app.yaml

create_secret:
	kubectl create secret generic my-credentials --from-file=config.json=./$(CONFIG_DIR)/config.json

deploy:
	docker rmi app:$(TAG)
	docker build -t app:$(TAG) ./$(SRC_DIR)
	kind load docker-image app:$(TAG) --name $(CLUSTER_NAME)
	kubectl apply -f ./$(MANIFEST_DIR)
	sleep 10
	kubectl get pods

clean:
	kubectl delete -f ./$(MANIFEST_DIR)
	kind delete cluster --name $(CLUSTER_NAME)
