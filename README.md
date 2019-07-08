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

### Nodeport

NodePortは、**podとnodeの外の通信** に使われる。      
yamlの中で、ポート番号を30000-32767の間で定義する。(ポート番号を指定しないと、ランダムで使われていないポートが割り振られる。)

今回の構成で、nodeの外と通信するのは、**nginx** と **private-registry** のpod。        
NodePortでserviceを定義したpodには、nodeの外から      
`NodeIP:NodePort` でアクセスできる。
今回だと、      
`NodeIP:80` は nginxのpodに。       
`NodeIP:30000` はdocker-registryのpodのフォワーディングされる。

### ClusterIP

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
