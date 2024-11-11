.PHONY: solr

MONITOR_NAMESPACE=monitoring
PROM_OPER_REL=mon
SOLR_OPERATOR_VERSION=0.8.1
SOLR_NAMESPACE=solr

k8s:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	helm install solr-operator apache-solr/solr-operator --version ${SOLR_OPERATOR_VERSION} \
  --namespace ${SOLR_NAMESPACE} \
  --create-namespace \
  --set zookeeper-operator.crd.create=true
	helm upgrade --install ${PROM_OPER_REL} prometheus-community/kube-prometheus-stack \
	--namespace ${MONITOR_NAMESPACE} \
	--create-namespace \
	--set kubeStateMetrics.enabled=true \
	--set nodeExporter.enabled=true
	helm install fluentd fluent/fluentd -f fluent-values.yaml

solr:
	curl "http://solr.local/solr/admin/collections?action=CREATE&name=books&numShards=2&replicationFactor=3&maxShardsPerNode=2"
	wget https://raw.githubusercontent.com/apache/solr/main/solr/example/exampledocs/books.json
	curl -X POST -H "Content-Type: application/json" \
  -d @books.json \
  "http://solr.local/solr/books/update/?commit=true"
	curl "http://solr.local/solr/books/select?q=*:*"

init:
	helm repo add apache-solr https://solr.apache.org//charts
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add stable https://charts.helm.sh/stable
	helm repo add fluent https://fluent.github.io/helm-charts
	helm repo update

del:
	helm uninstall fluentd
	kubectl delete -f manifest/

create:
	kubectl apply -f manifest/
	sleep 60
	kubectl wait --namespace ${SOLR_NAMESPACE} --for=condition=ready pod/example-solrcloud-0 --timeout=600s
	@make solr

delete:
	kind delete cluster

clean:
	@make delete
	helm repo remove apache-solr prometheus-community stable fluent
	rm books.json

setup:
	@make delete
	kind create cluster --config cluster.yaml
	@make k8s
	@make create
