# 场景一 Consul & Eureka & Nacos 服务融合

## 安装步骤

```bash
#make kind-up
#make metallb-up

make consul-deploy
make eureka-deploy
make nacos-deploy

make consul-port-forward
make eureka-port-forward
make nacos-port-forward

http://127.0.0.1:8500
http://127.0.0.1:8761
http://127.0.0.1:8848/nacos

# https://github.com/cybwan/go-demo/blob/main/scripts/deploy-fsm.sh
make deploy-fsm

make deploy-bookwarehouse

make deploy-bookwarehouse-3k

make undeploy-bookwarehouse-3k
```



