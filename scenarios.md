# 场景一 Consul & Eureka & Nacos 服务融合

## 场景描述

同一 k8s 下部署 consul & eureka & nacos 服务, 无须使用 fgw

consul 下部署 bookwarehouse 服务

eureka 下部署 bookstore 服务

nacos 下部署 bookbuyer 服务

K8s 下部署 httpbin 服务

## 安装步骤

```bash
make kind-up
make metallb-up

make consul-deploy
make eureka-deploy
make nacos-deploy
#make zookeeper-deploy

make consul-port-forward
make eureka-port-forward
make nacos-port-forward
#make zookeeper-port-forward
#make zookeeper-ui-port-forward

http://127.0.0.1:8500
http://127.0.0.1:8761
http://127.0.0.1:8848/nacos
#http://127.0.0.1:2181
#http://127.0.0.1:8080

make deploy-consul-bookwarehouse
make deploy-eureka-bookstore
make deploy-nacos-bookbuyer

make deploy-fsm
make deploy-httpbin

make port-forward-fsm-repo
http://127.0.0.1:6060

make bookbuyer-port-forward
http://127.0.0.1:14001
```
