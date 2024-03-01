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

export fsm_namespace=fsm-system
kubectl patch meshconfig fsm-mesh-config -n "$fsm_namespace" -p '{"spec":{"featureFlags":{"enableSidecarPrettyConfig":false}}}'  --type=merge

--set=fsm.featureFlags.enableSidecarPrettyConfig=false

export fsm_namespace=fsm-system
kubectl patch meshconfig fsm-mesh-config -n "$fsm_namespace" -p '{"spec":{"featureFlags":{"enableSidecarPrettyConfig":true}}}'  --type=merge


kubectl apply  -f - <<EOF
kind: ConsulConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: cluster1
spec:
  machineIP: 192.168.127.11
EOF


kubectl apply  -f - <<EOF
kind: ConsulConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: cluster2
spec:
  machineIP: 192.168.127.11
EOF

kubectl apply  -f - <<EOF
kind: ConsulConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: vm3
spec:
  machineIP: 192.168.127.11
EOF

kubectl create namespace derive-consul
fsm namespace add derive-consul
kubectl patch namespace derive-consul -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"consul"}}}'  --type=merge
```



