#!/bin/bash

set -aueo pipefail

# shellcheck disable=SC1091

number=0
while [ "$number" -lt 10 ]; do
number=$((number + 1))
padded_number=$(printf "%04d" $number)
echo "Number = $padded_number"
kubectl delete -n bookwarehouse -f - <<EOF
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
