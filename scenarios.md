# 场景一 Consul & Eureka & Nacos

## 场景描述

同一 k8s 下部署 consul & eureka & nacos 服务, 无须使用 fgw

consul 下部署 bookwarehouse 服务

eureka 下部署 bookstore 服务

nacos 下部署 bookbuyer 服务

## 安装步骤

```bash
make consul-deploy
make consul-port-forward
make eureka-deploy
make eureka-port-forward
make nacos-deploy
make nacos-port-forward

make deploy-consul-bookwarehouse
make deploy-eureka-bookstore
make deploy-nacos-bookbuyer

make deploy-fsm
```

