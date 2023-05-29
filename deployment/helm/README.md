## Install Istio gateway

    $ helm repo add istio https://istio-release.storage.googleapis.com/charts
    $ helm repo update

    $ kubectl create namespace istio-system
    $ helm install istio-base istio/base -n istio-system --wait
    $ helm install istiod istio/istiod -n istio-system --wait
    $ kubectl label namespace istio-system istio-injection=enabled
    $ helm install istio-ingressgateway istio/gateway -n istio-system --wait

### Allow namespace pods automatically setup istio-proxy

    $ kubectl label namespace default istio-injection=enabled

### Debugging

    $ kubectl logs istiod-{pod_id} -n istio-system
    $ kubectl exec -it -n istio-system istio-ingressgateway-{pod_id} -- bash
    $ > netstat -anp | grep 3306

### Debug istio-proxy sidecar

    $ kubectl logs mysql-primary-0 -c istio-proxy

## Install Bitnami MySQL cluster

    $ helm repo add bitnami https://charts.bitnami.com/bitnami
    $ kubectl apply -f mysql/secrets.yaml
    $ helm install mysql -f mysql/values.yaml bitnami/mysql

### Configure MySQL Istio gateway and load balancer

    $ helm apply -f mysql/gateway.yaml
    $ helm apply -f mysql/virtualservice.yaml
    $ helm apply -f istio/istio-ingressgateway.yaml

### Expose MySQL port to local for testing

    $ kubectl port-forward svc/istio-ingressgateway 3306:3306 -n istio-system
    $ mycli -u root -p {root_password} -h 127.0.0.1 -P 3306

### Configure Database

    MySQL root@127.0.0.1:(none)> CREATE USER 'jarvis'@'%' IDENTIFIED BY '{password}';
    > select host, user from mysql.user;
    > CREATE DATABASE jarvis CHARACTER SET utf8 COLLATE utf8_general_ci;
    > GRANT ALL PRIVILEGES ON jarvis.* TO 'jarvis'@'%';
    > flush privileges;
    > USE jarvis;

    > GRANT REPLICATION SLAVE ON *.* TO 'jarvis'@'%';

### Migrate schema

    $ goose -dir internal/db/migration mysql "jarvis:{password}@tcp(127.0.0.1:3306)/jarvis?charset=utf8" up

## Install Bitnami Redis Sentinel

    $ helm install redis-sentinel bitnami/redis --values redis/values.yaml

### Configure Redis Istio gateway and load balancer

    $ helm apply -f redis/gateway.yaml
    $ helm apply -f redis/virtualservice.yaml

    # Add port exposure into loadbalancer
    $ helm apply -f istio/istio-ingressgateway.yaml

### Expose Redis port to local for testing

    $ kubectl port-forward svc/istio-ingressgateway 6379:6379 -n istio-system

    $ yarn global add redis-cli
    $ rdcli -h 127.0.0.1 -p 6379 -a "{password}"
