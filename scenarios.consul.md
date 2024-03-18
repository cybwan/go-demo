# 场景一 Consul 服务融合

## 安装步骤

```bash
make kind-up
make metallb-up

make consul-deploy
CONSUL_VERSION=1.15.4 make consul-deploy

make deploy-fsm.consul

make deploy-bookwarehouse
make deploy-consul-bookwarehouse

make consul-port-forward
http://127.0.0.1:8500

fsm namespace remove bookwarehouse
kubectl rollout restart deployment -n bookwarehouse bookwarehouse

export LOOPS=3000;make deploy-bookwarehouse-3k
export LOOPS=3000;make undeploy-bookwarehouse-3k

make rebuild-fsm-connector restart-fsm-consul-connector
make tail-fsm-consul-connector-logs
```
