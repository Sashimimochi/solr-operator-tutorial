```bash
kind create cluster
kubectl apply -f init-pod.yaml
kubectl get pods
kubectl apply -f init-deployment.yaml
kubectl get deployments,pods
kubectl delete pod {podname}
kubectl get pods
kubectl delete -f init-deployment.yaml -f init-pod.yaml
```

```bash
kubectl apply -f nginx-dep.yaml
kubectl apply -f nginx-dep-svc.yaml
kubectl exec -it {podname} --bash # or sh
curl nginx-dep-svc -v --head
kubectl delete -f nginx-dep-svc.yaml -f nginx-dep.yaml
```
