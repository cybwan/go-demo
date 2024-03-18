# 场景一 Consul 服务融合

## 安装步骤

```bash
make kind-up
make metallb-up

make eureka-deploy

make deploy-fsm.eureka

make deploy-bookwarehouse
make deploy-eureka-bookwarehouse

make eureka-port-forward
http://127.0.0.1:8761

fsm namespace remove bookwarehouse
kubectl rollout restart deployment -n bookwarehouse bookwarehouse

export LOOPS=100;make deploy-bookwarehouse-3k
export LOOPS=100;make undeploy-bookwarehouse-3k

make rebuild-fsm-connector restart-fsm-eureka-connector
make tail-fsm-eureka-connector-logs
```
