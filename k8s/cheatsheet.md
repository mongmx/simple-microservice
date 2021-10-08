Deploy ingress-nginx
-

> kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.0.0/deploy/static/provider/cloud/deploy.yaml

---

Manage namespace
-

> kubectl create namespace tutorial

> kubectl get namespace

> kubectl delete --all tutorial

---

Manage pod
-

> kubectl get svc

> kubectl get pod

> kubectl get pod

---

Skaffold usage
-
> skaffold dev

> skaffold build --default-repo=registry.gitlab.com/pay9 --filename='skaffold-prod.yaml' --tag='latest'

> docker login registry.gitlab.com -u xxx -p xxx

---

Deploy manual
-

> kubectl --context k8s-cluster-name apply -f k8s/infra/deployment.yaml

> kubectl --context k8s-cluster-name apply -f k8s/dev/deployment.yaml

> kubectl --context k8s-cluster-name apply -f k8s/prod/deployment.yaml

---

Update after deploy
-
> kubectl --context k8s-cluster-name get pod

> kubectl --context k8s-cluster-name delete pod pod-name

---