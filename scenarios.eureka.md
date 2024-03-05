# 场景一 Consul 服务融合

## 安装步骤

```bash
make kind-up
make metallb-up

make eureka-deploy

make eureka-port-forward

http://127.0.0.1:8761

make deploy-fsm.eureka

make port-forward-fsm-repo
http://127.0.0.1:6060

make deploy-httpbin

make deploy-bookwarehouse

fsm namespace remove bookwarehouse
kubectl rollout restart deployment -n bookwarehouse bookwarehouse

export LOOPS=3000
make deploy-bookwarehouse-3k
make undeploy-bookwarehouse-3k

make build-fsm-cli

make rebuild-fsm-bootstrap restart-fsm-bootstrap

make rebuild-fsm-connector restart-fsm-eureka-connector
make tail-fsm-eureka-connector-logs
```
