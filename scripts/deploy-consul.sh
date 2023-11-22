#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091

export DEMO_HOME=https://raw.githubusercontent.com/flomesh-io/springboot-bookstore-demo/main
kubectl apply -n default -f $DEMO_HOME/manifests/consul.yaml
kubectl wait --all --for=condition=ready pod -n default -l app=consul --timeout=180s