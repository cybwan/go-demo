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

kubectl apply -n derive-consul -f - <<EOF
kind: AccessControl
apiVersion: policy.flomesh.io/v1alpha1
metadata:
  name: consul
spec:
  sources:
  - kind: IPRange
    name: 10.244.1.2/32
EOF

make deploy-httpbin

make deploy-bookwarehouse

fsm namespace remove bookwarehouse
kubectl rollout restart deployment -n bookwarehouse bookwarehouse

export LOOPS=100
make deploy-bookwarehouse-3k
make undeploy-bookwarehouse-3k

kubectl apply  -f - <<EOF
kind: ConsulConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: cluster1
spec:
  httpAddr: 127.0.0.1:8500
  deriveNamespace: derive-consul
  syncToK8S:
    enable: true
  syncFromK8S:
    enable: true
EOF

make build-fsm-cli

make rebuild-fsm-bootstrap restart-fsm-bootstrap

make rebuild-fsm-connector restart-fsm-consul-connector

make tail-fsm-consul-connector-logs
```



