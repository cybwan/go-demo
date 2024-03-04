# 场景一 Consul & Eureka & Nacos 服务融合

## 安装步骤

```bash
#make kind-up
#make metallb-up

make consul-deploy
make eureka-deploy
make nacos-deploy

make consul-port-forward
make eureka-port-forward
make nacos-port-forward

http://127.0.0.1:8500
http://127.0.0.1:8761
http://127.0.0.1:8848/nacos

# https://github.com/cybwan/go-demo/blob/main/scripts/deploy-fsm.sh
make deploy-fsm

make deploy-bookwarehouse

export LOOPS=100
make deploy-bookwarehouse-3k

export LOOPS=100
make undeploy-bookwarehouse-3k


make build-fsm-cli

make rebuild-fsm-bootstrap restart-fsm-bootstrap

make rebuild-fsm-connector restart-fsm-consul-connector

make tail-fsm-consul-connector-logs
```



