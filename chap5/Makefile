CLUSTER_NAME=""
TAG=9.3
OPERATOR_VERSION=0.7.1
NAMESPACE=solr

install-repo:
	helm repo add apache-solr https://solr.apache.org/charts
	helm repo update

create-cluster:
	kind create cluster --name $(CLUSTER_NAME) --config ./kind-portmapping-config.yaml

deploy:
	kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
#	末尾に以下が追加されていること
#	127.0.0.1 example.solr.com
	cat /etc/hosts
	kubectl get all -n projectcontour
	kubectl create -f https://solr.apache.org/operator/downloads/crds/v0.7.1/all-with-dependencies.yaml
	helm install solr-operator apache-solr/solr-operator --version $(OPERATOR_VERSION) \
		--create-namespace --namespace $(NAMESPACE)
	kubectl create namespace kubernetes-dashboard
	kubectl apply -f manifests/ --recursive
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.4/aio/deploy/recommended.yaml
	kubectl -n kubernetes-dashboard create token admin-user
	open "http://example.solr.com"
	kubectl proxy &
	open "http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/."

quickstart:
	make create-cluster
	make deploy

create_solr:
	kind create cluster --name $(CLUSTER_NAME) --config ./kind-portmapping-config.yaml
	kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
	helm install solr-operator apache-solr/solr-operator --version $(OPERATOR_VERSION) \
		--create-namespace --namespace $(NAMESPACE) \
		--set zookeeper-operator.crd.create=true
	helm install example-solr apache-solr/solr --version $(OPERATOR_VERSION) \
		-n solr \
		--set image.tag=$(TAG) \
		--set solrOptions.javaMemory="-Xms300m -Xmx300m" \
		--set addressability.external.method=Ingress \
  		--set addressability.external.domainName="solr.com" \
  		--set addressability.external.useExternalAddress="true" \
  		--set ingressOptions.ingressClassName="nginx"
	kubectl apply -f example-solrcloud-ingress.yaml
#	末尾に以下が追加されていること
#	127.0.0.1 example.solr.com
	cat /etc/hosts
	kubectl get all -n projectcontour
	kubectl get all -n $(NAMESPACE)
	kubectl get solrclouds -n $(NAMESPACE)
	kubectl get pods -n $(NAMESPACE)

index:
	rm books.json
	curl "http://example.solr.com/solr/admin/collections?action=CREATE&name=books&numShards=1&replicationFactor=2&maxShardsPerNode=1"
	wget https://raw.githubusercontent.com/apache/solr/main/solr/example/exampledocs/books.json
	sleep 30
	curl -X POST -H "Content-Type: application/json" \
    	-d @books.json \
    	"http://example.solr.com/solr/books/update/?commit=true"

config:
	curl -u admin:admin -X POST -H "Content-Type:application/octet-stream" --data-binary @conf.zip "http://example.solr.com/solr/admin/configs?action=UPLOAD&name=mini.AUTOCREATED&overwrite=true"

clean:
#	kubectl delete -f manifests/ --recursive
#	helm uninstall solr-operator -n $(NAMESPACE)
#	kubectl delete -f https://solr.apache.org/operator/downloads/crds/v0.7.1/all-with-dependencies.yaml
	kind delete cluster --name $(CLUSTER_NAME)
	pkill -9 -f "kubectl proxy"

clean-helm:
	helm uninstall example-solr -n $(NAMESPACE)
	helm uninstall solr-operator -n $(NAMESPACE)
	kubectl delete -f example-solrcloud-ingress.yaml
	kind delete cluster --name $(CLUSTER_NAME)
