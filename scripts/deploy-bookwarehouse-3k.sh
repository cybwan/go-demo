#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091

LOOPS="${LOOPS:-3000}"

number=$LOOPS
while [ "$number" -gt 0 ]; do
padded_number=$(printf "%04d" $number)
number=$((number - 1))
#echo "Number = $padded_number"
kubectl apply -n bookwarehouse -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: bookwarehouse$padded_number
  namespace: bookwarehouse
  labels:
    app: bookwarehouse
spec:
  ports:
    - port: 14001
      name: bookwarehouse-http-port
      appProtocol: HTTP
    - port: 9091
      name: bookwarehouse-grpc-port
      appProtocol: GRPC
  selector:
    app: bookwarehouse
EOF
done
