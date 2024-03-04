#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091

export fsm_namespace=fsm-system
export fsm_mesh_name=fsm
export dns_svc_ip="$(kubectl get svc -n kube-system -l k8s-app=kube-dns -o jsonpath='{.items[0].spec.clusterIP}')"
echo $dns_svc_ip
export consul_svc_addr="$(kubectl get svc -n default --field-selector metadata.name=consul -o jsonpath='{.items[0].spec.clusterIP}')"
echo $consul_svc_addr
export eureka_svc_addr="$(kubectl get svc -n default --field-selector metadata.name=eureka -o jsonpath='{.items[0].spec.clusterIP}')"
echo $eureka_svc_addr
export nacos_svc_addr="$(kubectl get svc -n default --field-selector metadata.name=nacos -o jsonpath='{.items[0].spec.clusterIP}')"
echo $nacos_svc_addr

fsm install \
    --mesh-name "$fsm_mesh_name" \
    --fsm-namespace "$fsm_namespace" \
    --set=fsm.certificateProvider.kind=tresor \
    --set=fsm.image.registry=localhost:5000/flomesh \
    --set=fsm.image.tag=latest \
    --set=fsm.image.pullPolicy=Always \
    --set=fsm.sidecar.sidecarLogLevel=warn \
    --set=fsm.controllerLogLevel=warn \
    --set=fsm.serviceAccessMode=mixed \
    --set=fsm.featureFlags.enableAutoDefaultRoute=true \
    --set=clusterSet.region=LN \
    --set=clusterSet.zone=DL \
    --set=clusterSet.group=FLOMESH \
    --set=clusterSet.name=LAB \
    --set fsm.fsmIngress.enabled=false \
    --set fsm.fsmGateway.enabled=true \
    --set=fsm.localDNSProxy.enable=true \
    --set=fsm.localDNSProxy.wildcard.enable=false \
    --set=fsm.localDNSProxy.primaryUpstreamDNSServerIPAddr=$dns_svc_ip \
    --set fsm.featureFlags.enableValidateHTTPRouteHostnames=false \
    --set fsm.featureFlags.enableValidateGRPCRouteHostnames=false \
    --set fsm.featureFlags.enableValidateTLSRouteHostnames=false \
    --set fsm.featureFlags.enableValidateGatewayListenerHostname=false \
    --set fsm.featureFlags.enableGatewayProxyTag=true \
    --set=fsm.featureFlags.enableSidecarPrettyConfig=false \
    --timeout=900s

kubectl create namespace derive-consul
fsm namespace add derive-consul
kubectl patch namespace derive-consul -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"consul"}}}'  --type=merge

kubectl create namespace derive-eureka
fsm namespace add derive-eureka
kubectl patch namespace derive-eureka -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"eureka"}}}'  --type=merge

kubectl create namespace derive-nacos
fsm namespace add derive-nacos
kubectl patch namespace derive-nacos -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"nacos"}}}'  --type=merge

kubectl create namespace derive-vm1
fsm namespace add derive-vm1
kubectl patch namespace derive-vm1 -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"machine"}}}'  --type=merge

kubectl create namespace derive-vm2
fsm namespace add derive-vm2
kubectl patch namespace derive-vm2 -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"machine"}}}'  --type=merge

export fsm_namespace=fsm-system
kubectl apply -n "$fsm_namespace" -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: k8s-fgw
spec:
  gatewayClassName: fsm-gateway-cls
  listeners:
    - protocol: HTTP
      port: 10080
      name: igrs-http
    - protocol: HTTP
      port: 10090
      name: egrs-http
    - protocol: HTTP
      port: 10180
      name: igrs-grpc
    - protocol: HTTP
      port: 10190
      name: egrs-grpc
EOF

kubectl apply  -f - <<EOF
kind: GatewayConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: fgw-1
spec:
  ingress:
    ipSelector: ExternalIP
    httpPort: 10080
    grpcPort: 10180
  egress:
    ipSelector: ClusterIP
    httpPort: 10090
    grpcPort: 10190
  syncToFgw:
    enable: true
    denyK8sNamespaces:
      - default
      - kube-system
      - fsm-system
EOF

kubectl apply  -f - <<EOF
kind: ConsulConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: cluster-1
spec:
  httpAddr: $consul_svc_addr:8500
  deriveNamespace: derive-consul
  syncToK8S:
    enable: true
    clusterId: consul_cluster_1
    suffixTag: version
    withGateway: true
  syncFromK8S:
    enable: true
    consulK8STag: k8s
    consulNodeName: k8s-sync
    appendTags:
      - tag0
      - tag1
    allowK8sNamespaces:
      - derive-eureka
      - derive-nacos
      - bookwarehouse
    withGateway: true
EOF

kubectl apply  -f - <<EOF
kind: EurekaConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: cluster-2
spec:
  httpAddr: http://$eureka_svc_addr:8761/eureka
  deriveNamespace: derive-eureka
  syncToK8S:
    enable: true
    clusterId: eureka_cluster_1
    suffixMetadata: version
    withGateway: true
  syncFromK8S:
    enable: true
    appendMetadatas:
      - key: type
        value: smart-gateway
      - key: version
        value: release
      - key: zone
        value: yinzhou
    allowK8sNamespaces:
      - derive-consul
      - derive-nacos
      - bookwarehouse
    withGateway: true
EOF

kubectl apply  -f - <<EOF
kind: NacosConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: cluster-3
spec:
  httpAddr: $nacos_svc_addr:8848
  deriveNamespace: derive-nacos
  syncToK8S:
    enable: true
    clusterId: nacos_cluster_1
    suffixMetadata: version
    withGateway: true
  syncFromK8S:
    enable: true
    allowK8sNamespaces:
      - derive-consul
      - derive-eureka
      - bookwarehouse
    withGateway: true
EOF

kubectl apply  -f - <<EOF
kind: MachineConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: vm-cluster-1
spec:
  deriveNamespace: derive-vm1
  syncToK8S:
    enable: true
    clusterId: vm_cluster_1
    withGateway: true
EOF

kubectl apply  -f - <<EOF
kind: MachineConnector
apiVersion: connector.flomesh.io/v1alpha1
metadata:
  name: vm-cluster-2
spec:
  deriveNamespace: derive-vm2
  syncToK8S:
    enable: true
    clusterId: vm_cluster_2
    withGateway: true
EOF

kubectl apply -n derive-vm1 -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm11
---
kind: VirtualMachine
apiVersion: machine.flomesh.io/v1alpha1
metadata:
  name: vm11
spec:
  serviceAccountName: vm11
  machineIP: 192.168.127.11
  services:
  - serviceName: hello11
    port: 10011
EOF

kubectl apply -n derive-vm1 -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm12
---
kind: VirtualMachine
apiVersion: machine.flomesh.io/v1alpha1
metadata:
  name: vm12
spec:
  serviceAccountName: vm12
  machineIP: 192.168.127.12
  services:
  - serviceName: hello12
    port: 10011
EOF

kubectl apply -n derive-vm1 -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm13
---
kind: VirtualMachine
apiVersion: machine.flomesh.io/v1alpha1
metadata:
  name: vm13
spec:
  serviceAccountName: vm13
  machineIP: 192.168.226.21
EOF

kubectl apply -n derive-vm1 -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm14
---
kind: VirtualMachine
apiVersion: machine.flomesh.io/v1alpha1
metadata:
  name: vm14
spec:
  serviceAccountName: vm14
  machineIP: 192.168.226.22
EOF

kubectl apply -n derive-vm2 -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm21
---
kind: VirtualMachine
apiVersion: machine.flomesh.io/v1alpha1
metadata:
  name: vm21
spec:
  serviceAccountName: vm21
  machineIP: 192.168.127.21
  services:
  - serviceName: world21
    port: 10011
EOF

kubectl apply -n derive-vm2 -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vm22
---
kind: VirtualMachine
apiVersion: machine.flomesh.io/v1alpha1
metadata:
  name: vm22
spec:
  serviceAccountName: vm22
  machineIP: 192.168.127.22
  services:
  - serviceName: world22
    port: 10011
EOF