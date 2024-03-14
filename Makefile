#!make

CONSUL_VERSION      ?= 1.5.3

.PHONY: consul-deploy
consul-deploy:
	kubectl apply -n default -f ./manifests/consul.$(CONSUL_VERSION).yaml
	kubectl wait --all --for=condition=ready pod -n default -l app=consul --timeout=180s

.PHONY: consul-reboot
consul-reboot:
	kubectl rollout restart deployment -n default consul

.PHONY: eureka-deploy
eureka-deploy:
	kubectl apply -n default -f ./manifests/eureka.yaml
	kubectl wait --all --for=condition=ready pod -n default -l app=eureka --timeout=180s

.PHONY: eureka-reboot
eureka-reboot:
	kubectl rollout restart deployment -n default eureka

.PHONY: nacos-deploy
nacos-deploy:
	kubectl apply -n default -f ./manifests/nacos.yaml
	kubectl wait --all --for=condition=ready pod -n default -l app=nacos --timeout=180s

.PHONY: nacos-reboot
nacos-reboot:
	kubectl rollout restart deployment -n default nacos

.PHONY: zookeeper-deploy
zookeeper-deploy:
	kubectl apply -n default -f ./manifests/zookeeper.yaml
	kubectl wait --all --for=condition=ready pod -n default -l app=zookeeper --timeout=180s

.PHONY: zookeeper-reboot
zookeeper-reboot:
	kubectl rollout restart deployment -n default zookeeper

.PHONY: deploy-curl
deploy-curl: undeploy-curl
	kubectl create namespace curl
	fsm namespace add curl
	kubectl apply -n curl -f ./manifests/curl.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n curl -l app=curl --timeout=180s

.PHONY: undeploy-curl
undeploy-curl:
	kubectl delete deployments.apps -n curl curl --ignore-not-found
	kubectl delete namespace curl --ignore-not-found

.PHONY: deploy-httpbin
deploy-httpbin: undeploy-httpbin
	kubectl create namespace httpbin
	fsm namespace add httpbin
	kubectl apply -n httpbin -f ./manifests/httpbin.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n httpbin -l app=httpbin --timeout=180s

.PHONY: undeploy-httpbin
undeploy-httpbin:
	kubectl delete deployments.apps -n httpbin httpbin --ignore-not-found
	kubectl delete namespace httpbin --ignore-not-found

.PHONY: deploy-bookwarehouse-3k
deploy-bookwarehouse-3k:
	scripts/deploy-bookwarehouse-3k.sh

.PHONY: undeploy-bookwarehouse-3k
undeploy-bookwarehouse-3k:
	scripts/undeploy-bookwarehouse-3k.sh

.PHONY: deploy-bookwarehouse
deploy-bookwarehouse: undeploy-bookwarehouse
	kubectl delete namespace bookwarehouse --ignore-not-found
	kubectl create namespace bookwarehouse
	#fsm namespace add bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/bookwarehouse.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-bookwarehouse
undeploy-bookwarehouse:
	kubectl delete -n bookwarehouse -f ./manifests/bookwarehouse.yaml --ignore-not-found

.PHONY: deploy-fsm-bookwarehouse
deploy-fsm-bookwarehouse: undeploy-fsm-bookwarehouse
	kubectl delete namespace bookwarehouse --ignore-not-found
	kubectl create namespace bookwarehouse
	fsm namespace add bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/bookwarehouse.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-fsm-bookwarehouse
undeploy-fsm-bookwarehouse:
	kubectl delete -n bookwarehouse -f ./manifests/bookwarehouse.yaml --ignore-not-found
	fsm namespace remove bookwarehouse || true

.PHONY: deploy-consul-bookwarehouse
deploy-consul-bookwarehouse:
	#kubectl delete namespace bookwarehouse --ignore-not-found
	#kubectl create namespace bookwarehouse
	#fsm namespace add bookwarehouse
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
	kubectl apply -n bookwarehouse -f ./manifests/eureka/bookwarehouse.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-eureka-bookwarehouse
undeploy-eureka-bookwarehouse:
	kubectl delete -n bookwarehouse -f ./manifests/eureka/bookwarehouse.yaml --ignore-not-found

.PHONY: deploy-eureka-bookstore
deploy-eureka-bookstore: undeploy-eureka-bookstore
	kubectl delete namespace bookstore --ignore-not-found
	kubectl create namespace bookstore
	kubectl apply -n bookstore -f ./manifests/eureka/bookstore.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookstore -l app=bookstore --timeout=180s

.PHONY: undeploy-eureka-bookstore
undeploy-eureka-bookstore:
	kubectl delete -n bookstore -f ./manifests/eureka/bookstore.yaml --ignore-not-found

.PHONY: deploy-eureka-bookbuyer
deploy-eureka-bookbuyer: undeploy-eureka-bookbuyer
	kubectl delete namespace bookbuyer --ignore-not-found
	kubectl create namespace bookbuyer
	kubectl apply -n bookbuyer -f ./manifests/eureka/bookbuyer.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookbuyer -l app=bookbuyer --timeout=180s

.PHONY: undeploy-eureka-bookbuyer
undeploy-eureka-bookbuyer:
	kubectl delete -n bookbuyer -f ./manifests/eureka/bookbuyer.yaml --ignore-not-found

.PHONY: deploy-eureka-bookthief
deploy-eureka-bookthief: undeploy-eureka-bookthief
	kubectl delete namespace bookthief --ignore-not-found
	kubectl create namespace bookthief
	kubectl apply -n bookthief -f ./manifests/eureka/bookthief.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookthief -l app=bookthief --timeout=180s

.PHONY: undeploy-eureka-bookthief
undeploy-eureka-bookthief:
	kubectl delete -n bookthief -f ./manifests/eureka/bookthief.yaml --ignore-not-found

.PHONY: deploy-nacos-bookwarehouse
deploy-nacos-bookwarehouse: undeploy-nacos-bookwarehouse
	kubectl delete namespace bookwarehouse --ignore-not-found
	kubectl create namespace bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/nacos/bookwarehouse.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-nacos-bookwarehouse
undeploy-nacos-bookwarehouse:
	kubectl delete -n bookwarehouse -f ./manifests/nacos/bookwarehouse.yaml --ignore-not-found

.PHONY: deploy-nacos-bookstore
deploy-nacos-bookstore: undeploy-nacos-bookstore
	kubectl delete namespace bookstore --ignore-not-found
	kubectl create namespace bookstore
	kubectl apply -n bookstore -f ./manifests/nacos/bookstore.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookstore -l app=bookstore --timeout=180s

.PHONY: undeploy-nacos-bookstore
undeploy-nacos-bookstore:
	kubectl delete -n bookstore -f ./manifests/nacos/bookstore.yaml --ignore-not-found

.PHONY: deploy-nacos-bookbuyer
deploy-nacos-bookbuyer: undeploy-nacos-bookbuyer
	kubectl delete namespace bookbuyer --ignore-not-found
	kubectl create namespace bookbuyer
	kubectl apply -n bookbuyer -f ./manifests/nacos/bookbuyer.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookbuyer -l app=bookbuyer --timeout=180s

.PHONY: undeploy-nacos-bookbuyer
undeploy-nacos-bookbuyer:
	kubectl delete -n bookbuyer -f ./manifests/nacos/bookbuyer.yaml --ignore-not-found

.PHONY: deploy-nacos-bookthief
deploy-nacos-bookthief: undeploy-nacos-bookthief
	kubectl delete namespace bookthief --ignore-not-found
	kubectl create namespace bookthief
	kubectl apply -n bookthief -f ./manifests/nacos/bookthief.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookthief -l app=bookthief --timeout=180s

.PHONY: undeploy-nacos-bookthief
undeploy-nacos-bookthief:
	kubectl delete -n bookthief -f ./manifests/nacos/bookthief.yaml --ignore-not-found

.PHONY: deploy-dubbo-bookwarehouse
deploy-dubbo-bookwarehouse: undeploy-dubbo-bookwarehouse
	kubectl delete namespace bookwarehouse --ignore-not-found
	kubectl create namespace bookwarehouse
	kubectl apply -n bookwarehouse -f ./manifests/dubbo/bookwarehouse.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookwarehouse -l app=bookwarehouse --timeout=180s

.PHONY: undeploy-dubbo-bookwarehouse
undeploy-dubbo-bookwarehouse:
	kubectl delete -n bookwarehouse -f ./manifests/dubbo/bookwarehouse.yaml --ignore-not-found

.PHONY: deploy-dubbo-bookstore
deploy-dubbo-bookstore: undeploy-dubbo-bookstore
	kubectl delete namespace bookstore --ignore-not-found
	kubectl create namespace bookstore
	kubectl apply -n bookstore -f ./manifests/dubbo/bookstore.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookstore -l app=bookstore --timeout=180s

.PHONY: undeploy-dubbo-bookstore
undeploy-dubbo-bookstore:
	kubectl delete -n bookstore -f ./manifests/dubbo/bookstore.yaml --ignore-not-found

.PHONY: deploy-dubbo-bookbuyer
deploy-dubbo-bookbuyer: undeploy-dubbo-bookbuyer
	kubectl delete namespace bookbuyer --ignore-not-found
	kubectl create namespace bookbuyer
	kubectl apply -n bookbuyer -f ./manifests/dubbo/bookbuyer.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookbuyer -l app=bookbuyer --timeout=180s

.PHONY: undeploy-dubbo-bookbuyer
undeploy-dubbo-bookbuyer:
	kubectl delete -n bookbuyer -f ./manifests/dubbo/bookbuyer.yaml --ignore-not-found

.PHONY: deploy-dubbo-bookthief
deploy-dubbo-bookthief: undeploy-dubbo-bookthief
	kubectl delete namespace bookthief --ignore-not-found
	kubectl create namespace bookthief
	kubectl apply -n bookthief -f ./manifests/dubbo/bookthief.yaml
	sleep 2
	kubectl wait --all --for=condition=ready pod -n bookthief -l app=bookthief --timeout=180s

.PHONY: undeploy-dubbo-bookthief
undeploy-dubbo-bookthief:
	kubectl delete -n bookthief -f ./manifests/dubbo/bookthief.yaml --ignore-not-found

.PHONY: deploy-nacos
deploy-nacos: deploy-nacos-bookwarehouse deploy-nacos-bookstore deploy-nacos-bookbuyer

.PHONY: undeploy-nacos
undeploy-nacos: undeploy-nacos-bookbuyer undeploy-nacos-bookstore undeploy-nacos-bookwarehouse

.PHONY: deploy-dubbo
deploy-dubbo: deploy-dubbo-bookwarehouse deploy-dubbo-bookstore deploy-dubbo-bookbuyer

.PHONY: undeploy-dubbo
undeploy-dubbo: undeploy-dubbo-bookbuyer undeploy-dubbo-bookstore undeploy-dubbo-bookwarehouse

.PHONY: consul-port-forward
consul-port-forward:
	export POD=$$(kubectl get pods --selector app=consul -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 8500:8500 --address 0.0.0.0 &

.PHONY: eureka-port-forward
eureka-port-forward:
	export POD=$$(kubectl get pods --selector app=eureka -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 8761:8761 --address 0.0.0.0 &

.PHONY: nacos-port-forward
nacos-port-forward:
	export POD=$$(kubectl get pods --selector app=nacos -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 8848:8848 --address 0.0.0.0 &

.PHONY: zookeeper-port-forward
zookeeper-port-forward:
	export POD=$$(kubectl get pods --selector app=zookeeper -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 2181:2181 --address 0.0.0.0 &

.PHONY: zookeeper-ui-port-forward
zookeeper-ui-port-forward:
	export POD=$$(kubectl get pods --selector app=zookeeper -n default --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n default 8080:8081 --address 0.0.0.0 &

.PHONY: bookbuyer-port-forward
bookbuyer-port-forward:
	export POD=$$(kubectl get pods --selector app=bookbuyer -n bookbuyer --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl port-forward "$$POD" -n bookbuyer 14001:14001 --address 0.0.0.0 &

.PHONY: bookbuyer-logs
bookbuyer-logs:
	export POD=$$(kubectl get pods --selector app=bookbuyer -n bookbuyer --no-headers | grep 'Running' | awk 'NR==1{print $$1}');\
	kubectl logs "$$POD" -n bookbuyer -c bookbuyer -f

.PHONY: to-consul
to-consul:
	cd cmd/to-consul;go run .

.PHONY: eureka-reg-services
eureka-reg-services:
	go run eureka/reg/main.go

.PHONY: eureka-list-services
eureka-list-services:
	go run eureka/list/main.go | jq

.PHONY: deploy-fsm-eureka
deploy-fsm-eureka:
	scripts/deploy-fsm-eureka.sh

.PHONY: deploy-fsm-eureka-c2
deploy-fsm-eureka-c2:
	scripts/deploy-fsm-eureka-c2.sh

.PHONY: deploy-fsm-eureka-c3
deploy-fsm-eureka-c3:
	scripts/deploy-fsm-eureka-c3.sh

.PHONY: deploy-fsm-nacos
deploy-fsm-nacos:
	scripts/deploy-fsm-nacos.sh

.PHONY: deploy-fsm-consul
deploy-fsm-consul:
	scripts/deploy-fsm-consul.sh

.PHONY: deploy-fsm
deploy-fsm:
	scripts/deploy-fsm.sh

.PHONY: deploy-fsm.consul
deploy-fsm.consul:
	scripts/deploy-fsm.consul.sh

.PHONY: deploy-fsm.eureka
deploy-fsm.eureka:
	scripts/deploy-fsm.eureka.sh

.PHONY: deploy-fsm.nacos
deploy-fsm.nacos:
	scripts/deploy-fsm.nacos.sh

.PHONY: deploy-fsm-min
deploy-fsm-min:
	scripts/deploy-fsm-min.sh

.PHONY: undeploy-fsm
undeploy-fsm:
	fsm uninstall mesh --delete-cluster-wide-resources --delete-namespace|| true
	kubectl delete namespace derive-vm || true
	kubectl delete namespace derive-vm1 || true
	kubectl delete namespace derive-vm2 || true
	kubectl delete namespace derive-eureka || true
	kubectl delete namespace derive-eureka1 || true
	kubectl delete namespace derive-eureka2 || true
	kubectl delete namespace derive-consul || true
	kubectl delete namespace derive-consul1 || true
	kubectl delete namespace derive-consul2 || true
	kubectl delete namespace derive-nacos || true
	kubectl delete namespace derive-nacos1 || true
	kubectl delete namespace derive-nacos2 || true
