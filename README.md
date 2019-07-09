![k8s-image](https://user-images.githubusercontent.com/47022289/60791977-bf211880-a19f-11e9-87d2-9962450d0f8e.png)

# pod

kubernetesの最小単位は **pod** 。       
podは、**dockerコンテナ** の集合体。       
podが立つことは、dockerコンテナが立つことを意味する。      
podは、**node** の上に立つ。        
今回のnodeには、ローカル開発用のツールである **minikube** を用いている。       
nodeには、IPアドレスがある。

# deployment
podの定義を、**deployment**  という。        
deploymentは、yamlで書く。        

# service

podの通信を制御するヤツを **service** という。       
serviceも、yamlで書く。        
serviceで定義できる通信の仕方には4種類あるが、今回の構成で用いているのは2種類。        
**NodePort** と **ClusterIP** である。        

### ・Nodeport

NodePortは、**podとnodeの外の通信** に使われる。      
yamlの中で、ポート番号を30000-32767の間で定義する。(ポート番号を指定しないと、ランダムで使われていないポートが割り振られる。)

今回の構成で、nodeの外と通信するのは、**nginx** と **private-registry** のpod。        
NodePortでserviceを定義したpodには、nodeの外から      
`NodeIP:NodePort` でアクセスできる。
今回だと、      
`NodeIP:80` は nginxのpodに。       
`NodeIP:30000` はdocker-registryのpodのフォワーディングされる。

### ・ClusterIP

ClusterIPは、**同じnode上のpod間の通信** に使われる。     
yamlの中では少なくとも、**serviceName**、**port**、**targetPort** を定義する。       
ClusterIPはランダムで生成されるが、pod間の通信は、serviceNameで名前解決ができるので、ClusterIPをわざわざ指定する必要はないように思える。

serviceは、`serviceName:port` に来たリクエストを、     
`targetPort` をリッスンしているpodに渡す。    

今回は、nginxのpodからリバースプロキシで転送する際、      
転送先を `bluehub-gatewayのserviceName:port` と指定している。
こうすることで、bluehub-gatewayのpodと通信ができる。

bluehub-gatewayとbluehub-routerがgRPC通信する際も、
endpointに、     
`bluehub-routerのserviseName:port` を指定している。 

# docker-registryについて

venderによらない構成を考えた結果、    
docker-imageを保存するregistryをpodとして用意し、     
そこから他の3つのpod内コンテナのimageを持ってきている。


# ビルド

以下のシェルスクリプトで環境ができます。

```
$ sh apply.sh
```

### ・deploymentの適応

```
$ kubectl apply -f nginx/nginx-deployment.yaml
$ kubectl apply -f bluehub-gateway/bluehub-gateway-deployment.yaml
$ kubectl apply -f bluehub-router/bluehub-router-deployment.yaml
$ kubectl apply -f private-registry/private-registry-deployment.yaml
```
### ・podができたか確認

statusがRunningになっていればOK
```
$ kubectl get pods
NAME                                           READY   STATUS    RESTARTS   AGE
bluehub-gateway-deployment-9db874877-pqmwd     1/1     Running   0          42h
bluehub-router-deployment-5d479584c9-npc57     1/1     Running   0          45h
nginx-deployment-7cdf8cddcf-xgxlk              1/1     Running   0          42h
private-registry-deployment-794d874968-xwrwj   1/1     Running   0          22h
```

### ・serviceの適応

```
$ kubectl apply -f nginx/nginx-service.yaml
$ kubectl apply -f bluehub-gateway/bluehub-gateway-service.yaml
$ kubectl apply -f bluehub-router/bluehub-router-service.yaml
$ kubectl apply -f private-registry/private-registry-service.yaml
```

### ・serviceができたか確認

kubernetesというserviceは、自分で定義しなくても最初からある。
```
$ kubectl get service
NAME                       TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)           AGE
bluehub-gateway-service    ClusterIP   10.96.49.156     <none>        10000/TCP         45h
bluehub-router-service     ClusterIP   10.109.114.71    <none>        10001/TCP         45h
kubernetes                 ClusterIP   10.96.0.1        <none>        443/TCP           3d
nginx-service              NodePort    10.110.121.101   <none>        80:30005/TCP      2d21h
private-registry-service   NodePort    10.107.254.52    <none>        10002:30000/TCP   22h
```

### ・dashboard

```
$ minikube dashboard
```

# 使い方

まずは、nodeのIPを確かめる。

```
$ kubectl get node -o wide
NAME       STATUS   ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE            KERNEL-VERSION   CONTAINER-RUNTIME
minikube   Ready    master   3d    v1.14.0   10.0.2.15     <none>        Buildroot 2018.05   4.15.0           docker://18.6.2
```
本来は、ここのEXTERNAL-IPというのが、NodeIPにあたる。     
だが、minikubeは仕様上、このEXTERNAL-IPが存在しない。     
なので、minikubeであれば、以下のコマンド。      

```
$ minikube status
minikube status
host: Running
kubelet: Running
apiserver: Running
kubectl: Correctly Configured: pointing to minikube-vm at 192.168.99.101
```

このminikube-vmのIPアドレスを、NodeIPのように扱える。

ブラウザから、       
```
http://192.168.99.101:30005/v1/example/sayhello/hirasawa
```

にアクセスすると、
        
```
{"message":"Hello hirasawa"}
```

とレスポンスが返ってくる。

# 引用

gRPC部分はほとんど他人様のコードです。    
[https://qiita.com/ryu3/items/b2882d4f45c7f8485030#3grpc-stub-%E3%82%92%E7%94%9F%E6%88%90%E3%81%99%E3%82%8B](https://qiita.com/ryu3/items/b2882d4f45c7f8485030#3grpc-stub-%E3%82%92%E7%94%9F%E6%88%90%E3%81%99%E3%82%8B)

パッケージはこちら    
[https://github.com/grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
