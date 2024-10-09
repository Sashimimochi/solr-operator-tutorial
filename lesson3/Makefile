.PHONY: solr

MONITOR_NS=monitoring
PROM_OPER_REL=mon

k8s:
#	helm install fluentd-operator fluent/fluent-operator --version 3.2.0 \
	--namespace fluent \
	--create-namespace \
	--set installCRDs=true
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	helm install solr-operator apache-solr/solr-operator --version 0.8.1 \
  --namespace solr \
  --create-namespace \
  --set zookeeper-operator.crd.create=true
# see: https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
	helm upgrade --install ${PROM_OPER_REL} prometheus-community/kube-prometheus-stack \
	--namespace ${MONITOR_NS} \
	--create-namespace \
	--set kubeStateMetrics.enabled=true \
	--set nodeExporter.enabled=true
	helm install fluentd fluent/fluentd -f fluent-values.yaml
#	sleep 30

solr:
	curl "http://solr.local/solr/admin/collections?action=CREATE&name=books&numShards=2&replicationFactor=3&maxShardsPerNode=2"
	curl -X POST -H "Content-Type: application/json" \
  -d @books.json \
  "http://solr.local/solr/books/update/?commit=true"

init:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add stable https://charts.helm.sh/stable
	helm repo update
#	helm repo add fluent https://fluent.github.io/helm-charts

del:
#	helm uninstall fluentd
	kubectl delete -f manifest/

create:
#	helm install fluentd fluent/fluentd
	kubectl apply -f manifest/
#	helm upgrade fluentd fluent/fluentd -f fluent-values.yaml
	sleep 200
	@make solr

setup:
	kind create cluster --config cluster.yaml
	@make k8s
	@make create

delete:
	kind delete cluster
