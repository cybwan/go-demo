#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091

export fsm_namespace=fsm-system
export fsm_mesh_name=fsm
export dns_svc_ip="$(kubectl get svc -n kube-system -l k8s-app=kube-dns -o jsonpath='{.items[0].spec.clusterIP}')"
echo $dns_svc_ip
export eureka_svc_addr="$(kubectl get svc -n default --field-selector metadata.name=eureka -o jsonpath='{.items[0].spec.clusterIP}')"
echo $eureka_svc_addr

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
    --set=fsm.localDNSProxy.wildcard.enable=true \
    --set=fsm.localDNSProxy.primaryUpstreamDNSServerIPAddr=$dns_svc_ip \
    --set fsm.featureFlags.enableValidateHTTPRouteHostnames=false \
    --set fsm.featureFlags.enableValidateGRPCRouteHostnames=false \
    --set fsm.featureFlags.enableValidateTLSRouteHostnames=false \
    --set fsm.featureFlags.enableValidateGatewayListenerHostname=false \
    --set fsm.featureFlags.enableGatewayProxyTag=true \
    --set=fsm.cloudConnector.eureka.enable=true \
    --set=fsm.cloudConnector.eureka.deriveNamespace=derive-eureka \
    --set=fsm.cloudConnector.eureka.httpAddr=http://$eureka_svc_addr:8761/eureka \
    --set=fsm.cloudConnector.eureka.syncToK8S.enable=true \
    --set=fsm.cloudConnector.eureka.syncToK8S.passingOnly=false \
    --set=fsm.cloudConnector.eureka.syncToK8S.suffixMetadata=version \
    --set=fsm.cloudConnector.eureka.syncToK8S.withGateway.enable=true \
    --set=fsm.cloudConnector.eureka.syncFromK8S.enable=true \
    --set "fsm.cloudConnector.eureka.syncFromK8S.denyK8sNamespaces={default,kube-system,fsm-system}" \
    --set=fsm.cloudConnector.eureka.syncFromK8S.withGateway.enable=true \
    --set=fsm.cloudConnector.machine.enable=true \
    --set=fsm.cloudConnector.machine.asInternalServices=false \
    --set=fsm.cloudConnector.machine.deriveNamespace=derive-vm \
    --set=fsm.cloudConnector.machine.syncToK8S.enable=true \
    --set=fsm.cloudConnector.machine.syncToK8S.withGatewayEgress.enable=true \
    --set=fsm.cloudConnector.gateway.ingress.ipSelector=ClusterIP \
    --set=fsm.cloudConnector.gateway.ingress.httpPort=10080 \
    --set=fsm.cloudConnector.gateway.ingress.grpcPort=10180 \
    --set=fsm.cloudConnector.gateway.egress.ipSelector=ClusterIP \
    --set=fsm.cloudConnector.gateway.egress.httpPort=10090 \
    --set=fsm.cloudConnector.gateway.egress.grpcPort=10190 \
    --set=fsm.cloudConnector.gateway.syncToFgw.enable=true \
    --set "fsm.cloudConnector.gateway.syncToFgw.denyK8sNamespaces={default,kube-system,fsm-system}" \
    --timeout=900s

kubectl create namespace derive-eureka
fsm namespace add derive-eureka
kubectl patch namespace derive-eureka -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"eureka"}}}'  --type=merge

kubectl create namespace derive-vm
fsm namespace add derive-vm
kubectl patch namespace derive-vm -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"machine"}}}'  --type=merge

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
      name: ingress-proxy
    - protocol: HTTP
      port: 10090
      name: egress-proxy
    - protocol: HTTP
      port: 10180
      name: ingress-grpc-proxy
    - protocol: HTTP
      port: 10190
      name: egress-grpc-proxy
EOF

#kubectl apply -n derive-vm -f - <<EOF
#apiVersion: v1
#kind: ServiceAccount
#metadata:
#  name: vm1
#---
#kind: VirtualMachine
#apiVersion: machine.flomesh.io/v1alpha1
#metadata:
#  name: vm6
#spec:
#  serviceAccountName: vm1
#  machineIP: 192.168.127.8
#  services:
#  - serviceName: bookwarehouse
#    port: 14001
#  - serviceName: bookdemo
#    port: 10011
#EOF