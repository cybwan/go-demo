# 场景一 Consul 服务融合

## 安装步骤

```bash
make kind-up
make metallb-up

make consul-deploy

make consul-port-forward

http://127.0.0.1:8500

make deploy-fsm.consul

make port-forward-fsm-repo
http://127.0.0.1:6060

make deploy-httpbin

make deploy-bookwarehouse

fsm namespace remove bookwarehouse
kubectl rollout restart deployment -n bookwarehouse bookwarehouse

export LOOPS=100
make deploy-bookwarehouse-3k
make undeploy-bookwarehouse-3k

make build-fsm-cli

make rebuild-fsm-bootstrap restart-fsm-bootstrap

make rebuild-fsm-connector restart-fsm-consul-connector

make tail-fsm-consul-connector-logs
```
