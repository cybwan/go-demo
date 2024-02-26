#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091

export fsm_namespace=fsm-system
export fsm_mesh_name=fsm
export dns_svc_ip="$(kubectl get svc -n kube-system -l k8s-app=kube-dns -o jsonpath='{.items[0].spec.clusterIP}')"
echo $dns_svc_ip
export consul_svc_addr="$(kubectl get svc -n default --field-selector metadata.name=consul -o jsonpath='{.items[0].spec.clusterIP}')"
echo $consul_svc_addr

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
    --set=fsm.cloudConnector.consul.enable=true \
    --set=fsm.cloudConnector.consul.deriveNamespace=derive-consul \
    --set=fsm.cloudConnector.consul.httpAddr=$consul_svc_addr:8500 \
    --set=fsm.cloudConnector.consul.syncToK8S.enable=false \
    --set=fsm.cloudConnector.consul.syncToK8S.clusterId=consul_cluster_1 \
    --set=fsm.cloudConnector.consul.syncToK8S.suffixTag=version \
    --set=fsm.cloudConnector.consul.syncToK8S.withGateway.enable=true \
    --set=fsm.cloudConnector.consul.syncFromK8S.enable=true \
    --set fsm.cloudConnector.consul.syncFromK8S.appendTag[0]=tag0 \
    --set fsm.cloudConnector.consul.syncFromK8S.appendTag[1]=tag1 \
    --set "fsm.cloudConnector.consul.syncFromK8S.allowK8sNamespaces={derive-eureka,derive-nacos,bookwarehouse}" \
    --set=fsm.cloudConnector.consul.syncFromK8S.withGateway.enable=true \
    --set=fsm.cloudConnector.gateway.ingress.ipSelector=ExternalIP \
    --set=fsm.cloudConnector.gateway.ingress.httpPort=10080 \
    --set=fsm.cloudConnector.gateway.ingress.grpcPort=10180 \
    --set=fsm.cloudConnector.gateway.egress.ipSelector=ClusterIP \
    --set=fsm.cloudConnector.gateway.egress.httpPort=10090 \
    --set=fsm.cloudConnector.gateway.egress.grpcPort=10190 \
    --set=fsm.cloudConnector.gateway.syncToFgw.enable=true \
    --set "fsm.cloudConnector.gateway.syncToFgw.denyK8sNamespaces={default,kube-system,fsm-system}" \
    --timeout=900s

kubectl create namespace derive-consul
fsm namespace add derive-consul
kubectl patch namespace derive-consul -p '{"metadata":{"annotations":{"flomesh.io/mesh-service-sync":"consul"}}}'  --type=merge

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
