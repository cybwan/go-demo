# 场景一 Consul 服务融合

## 安装步骤

```bash
make kind-up
make metallb-up

make nacos-deploy

make nacos-port-forward

http://127.0.0.1:8848/nacos

make deploy-fsm.nacos

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

make rebuild-fsm-connector restart-fsm-nacos-connector
make tail-fsm-nacos-connector-logs
```
