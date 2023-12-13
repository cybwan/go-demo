## 部署业务 POD

```bash

POD=$(kubectl get pods --selector app=consul -n default --no-headers | grep 'Running' | awk 'NR==1{print $1}')
kubectl port-forward "$POD" -n default 8500:8500 --address 0.0.0.0

#模拟业务服务
export DEMO_HOME=https://raw.githubusercontent.com/flomesh-io/springboot-bookstore-demo/main

kubectl create namespace bookwarehouse
kubectl create namespace bookstore
kubectl create namespace bookbuyer

kubectl apply -n bookwarehouse -f $DEMO_HOME/manifests/consul/bookwarehouse-consul.yaml
kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

kubectl apply -n bookstore -f $DEMO_HOME/manifests/consul/bookstore-consul.yaml
kubectl apply -n bookstore -f $DEMO_HOME/manifests/consul/bookstore-v2-consul.yaml
kubectl wait --all --for=condition=ready pod -n bookstore -l app=bookstore --timeout=180s

kubectl apply -n bookbuyer -f $DEMO_HOME/manifests/consul/bookbuyer-consul.yaml
kubectl wait --all --for=condition=ready pod -n bookbuyer -l app=bookbuyer --timeout=180s





BUYER_V1_POD="$(kubectl get pods --selector app=bookbuyer,version=v1 -n bookbuyer --no-headers | grep 'Running' | awk 'NR==1{print $1}')"
kubectl port-forward $BUYER_V1_POD -n bookbuyer 8081:14001 --address 0.0.0.0 &

STORE_V1_POD="$(kubectl get pods --selector app=bookstore,version=v1 -n bookstore --no-headers | grep 'Running' | awk 'NR==1{print $1}')"
kubectl port-forward $STORE_V1_POD -n bookstore 8084:14001 --address 0.0.0.0 &

STORE_V2_POD="$(kubectl get pods --selector app=bookstore,version=v2 -n bookstore --no-headers | grep 'Running' | awk 'NR==1{print $1}')"
kubectl port-forward $STORE_V2_POD -n bookstore 8082:14001 --address 0.0.0.0 &

https://github.com/srvaroa/eurek8s/tree/master

https://srvaroa.github.io/kubernetes/eureka/paas/microservices/loadbalancing/2020/02/12/eureka-kubernetes.html

https://www.zeng.dev/post/20200428-eureka-multil-cluster-replica/

https://github.com/Haybu/spring-cloud-k8s-eureka-controller

https://www.python100.com/html/104745.html

https://medium.com/@fahimfarookme/the-mystery-of-eureka-health-monitoring-fd05fe757928
https://blog.csdn.net/xiao_jun_0820/article/details/77991963

https://kubernetes.io/docs/reference/using-api/health-checks/

https://www.reddit.com/r/kubernetes/comments/15l8ofh/how_to_check_health_of_api_server_endpoints_using/

https://blog.csdn.net/weixin_39128119/article/details/110875550

github.com/SimonWang00/goeureka
```

