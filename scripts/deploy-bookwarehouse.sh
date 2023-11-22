#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091
CTR_REGISTRY="${CTR_REGISTRY:-localhost:5000/flomesh}"
CTR_TAG="${CTR_TAG:-latest}"

BOOKWAREHOUSE_NAMESPACE="${BOOKWAREHOUSE_NAMESPACE:-bookwarehouse}"

KUBE_CONTEXT=$(kubectl config current-context)

kubectl delete namespace "$BOOKWAREHOUSE_NAMESPACE" --ignore-not-found

kubectl create namespace "$BOOKWAREHOUSE_NAMESPACE"

echo -e "Deploy Bookwarehouse Service Account"
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: bookwarehouse
  namespace: $BOOKWAREHOUSE_NAMESPACE
EOF

echo -e "Deploy Bookwarehouse Service"
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: bookwarehouse
  namespace: $BOOKWAREHOUSE_NAMESPACE
  labels:
    app: bookwarehouse
spec:
  ports:
  - port: 14001
    name: bookwarehouse-port

  selector:
    app: bookwarehouse
EOF

echo -e "Deploy Bookwarehouse Deployment"
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bookwarehouse
  namespace: "$BOOKWAREHOUSE_NAMESPACE"
spec:
  replicas: 1
  selector:
    matchLabels:
      app: bookwarehouse
  template:
    metadata:
      labels:
        app: bookwarehouse
        version: v1
    spec:
      serviceAccountName: bookwarehouse
      containers:
        # Main container with APP
        - name: bookwarehouse
          image: "${CTR_REGISTRY}/go-demo-bookwarehouse:${CTR_TAG}"
          imagePullPolicy: Always
          command: ["/bookwarehouse"]
          env:
            - name: IDENTITY
              value: bookwarehouse.${KUBE_CONTEXT}
EOF

kubectl wait --all --for=condition=ready pod -n "$BOOKWAREHOUSE_NAMESPACE" -l app=bookwarehouse --timeout=180s

