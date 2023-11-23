#!make

.PHONY: deploy-consul
deploy-consul:
	kubectl apply -n default -f ./manifests/consul.yaml
	kubectl wait --all --for=condition=ready pod -n default -l app=consul --timeout=180s

.PHONY: reboot-consul
reboot-consul:
	kubectl rollout restart deployment -n default consul

.PHONY: deploy-eureka
deploy-eureka:
	kubectl apply -n default -f ./manifests/eureka.yaml
	kubectl wait --all --for=condition=ready pod -n default -l app=eureka --timeout=180s

.PHONY: reboot-eureka
reboot-eureka:
	kubectl rollout restart deployment -n default eureka

.PHONY: deploy-bookwarehouse
deploy-bookwarehouse: undeploy-bookwarehouse
	kubectl create namespace bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/bookwarehouse.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-bookwarehouse
undeploy-bookwarehouse:
	kubectl delete deployments.apps -n bookwarehouse bookwarehouse --ignore-not-found
	kubectl delete namespace bookwarehouse --ignore-not-found

.PHONY: deploy-consul-bookwarehouse
deploy-consul-bookwarehouse: undeploy-consul-bookwarehouse
	kubectl delete namespace bookwarehouse --ignore-not-found
	kubectl create namespace bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/consul/bookwarehouse-consul.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-consul-bookwarehouse
undeploy-consul-bookwarehouse:
	kubectl delete deployments.apps -n bookwarehouse bookwarehouse --ignore-not-found
	kubectl delete namespace bookwarehouse --ignore-not-found

.PHONY: deploy-consul-bookstore
deploy-consul-bookstore: undeploy-consul-bookstore
	kubectl delete namespace bookstore --ignore-not-found
	kubectl create namespace bookstore
	kubectl apply -n bookstore -f ./manifests/consul/bookstore-consul.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookstore -l app=bookstore --timeout=180s

.PHONY: undeploy-consul-bookstore
undeploy-consul-bookstore:
	kubectl delete deployments.apps -n bookstore bookstore --ignore-not-found
	kubectl delete namespace bookstore --ignore-not-found

.PHONY: deploy-consul-bookbuyer
deploy-consul-bookbuyer: undeploy-consul-bookbuyer
	kubectl delete namespace bookbuyer --ignore-not-found
	kubectl create namespace bookbuyer
	kubectl apply -n bookbuyer -f ./manifests/consul/bookbuyer-consul.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookbuyer -l app=bookbuyer --timeout=180s

.PHONY: undeploy-consul-bookbuyer
undeploy-consul-bookbuyer:
	kubectl delete deployments.apps -n bookbuyer bookbuyer --ignore-not-found
	kubectl delete namespace bookbuyer --ignore-not-found

.PHONY: deploy-eureka-bookwarehouse
deploy-eureka-bookwarehouse: undeploy-eureka-bookwarehouse
	kubectl delete namespace bookwarehouse --ignore-not-found
	kubectl create namespace bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/eureka/bookwarehouse-eureka.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-eureka-bookwarehouse
undeploy-eureka-bookwarehouse:
	kubectl delete deployments.apps -n bookwarehouse bookwarehouse --ignore-not-found
	kubectl delete namespace bookwarehouse --ignore-not-found

.PHONY: deploy-eureka-bookstore
deploy-eureka-bookstore: undeploy-eureka-bookstore
	kubectl delete namespace bookstore --ignore-not-found
	kubectl create namespace bookstore
	kubectl apply -n bookstore -f ./manifests/eureka/bookstore-eureka.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookstore -l app=bookstore --timeout=180s

.PHONY: undeploy-eureka-bookstore
undeploy-eureka-bookstore:
	kubectl delete deployments.apps -n bookstore bookstore --ignore-not-found
	kubectl delete namespace bookstore --ignore-not-found

.PHONY: deploy-eureka-bookbuyer
deploy-eureka-bookbuyer: undeploy-eureka-bookbuyer
	kubectl delete namespace bookbuyer --ignore-not-found
	kubectl create namespace bookbuyer
	kubectl apply -n bookbuyer -f ./manifests/eureka/bookbuyer-eureka.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookbuyer -l app=bookbuyer --timeout=180s

.PHONY: undeploy-eureka-bookbuyer
undeploy-eureka-bookbuyer:
	kubectl delete deployments.apps -n bookbuyer bookbuyer --ignore-not-found
	kubectl delete namespace bookbuyer --ignore-not-found

.PHONY: port-forward-consul
port-forward-consul:
	export POD=$$(kubectl get pods --selector app=consul -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 8500:8500 --address 0.0.0.0 &

.PHONY: port-forward-eureka
port-forward-eureka:
	export POD=$$(kubectl get pods --selector app=eureka -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 8761:8761 --address 0.0.0.0 &

.PHONY: to-consul
to-consul:
	cd cmd/to-consul;go run .