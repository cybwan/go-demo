# 场景一 Nacos 服务融合

## 安装步骤

```bash
make kind-up
make metallb-up

make nacos-deploy

make deploy-fsm.nacos

make deploy-bookwarehouse

make nacos-port-forward
http://127.0.0.1:8848/nacos

fsm namespace remove bookwarehouse
kubectl rollout restart deployment -n bookwarehouse bookwarehouse

export LOOPS=3000;make deploy-bookwarehouse-3k
export LOOPS=3000;make undeploy-bookwarehouse-3k

make rebuild-fsm-connector restart-fsm-nacos-connector
make tail-fsm-nacos-connector-logs
```
