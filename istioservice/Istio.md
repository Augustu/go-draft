### Istio

```bash
# ref: https://istio.io/latest/docs/setup/getting-started/
# ref: https://istio.io/latest/docs/examples/bookinfo/

# ref: https://istio.io/latest/docs/setup/additional-setup/config-profiles/

istioctl install --set profile=demo -y

# Add a namespace label to instruct Istio to automatically inject Envoy sidecar proxies when you deploy your application later
kubectl label namespace default istio-injection=enabled

```



**Deploy the sample application**

```bash
# git clone https://github.com/istio/istio

kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml

kubectl get services

kubectl get pods

kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"

```



**Determine the ingress IP and port**

```bash
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml

kubectl get gateway

# istio-ingressgateway -> NodePort

kubectl get svc istio-ingressgateway -n istio-system

export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
# export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
# export TCP_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="tcp")].nodePort}')

export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

curl -s "http://${GATEWAY_URL}/productpage" | grep -o "<title>.*</title>"

http://$GATEWAY_URL/productpage

```



**Apply default destination rules**

```bash
kubectl get destinationrules -o yaml

kubectl apply -f samples/bookinfo/networking/destination-rule-all.yaml

kubectl get destinationrules -o yaml

```



**View the dashboard**

```bash
kubectl apply -f samples/addons

kubectl rollout status deployment/kiali -n istio-system

istioctl dashboard kiali

curl http://$GATEWAY_URL/productpage

watch -n 1 curl -o /dev/null -s -w %{http_code} $GATEWAY_URL/productpage

for i in $(seq 1 100); do curl -s -o /dev/null "http://$GATEWAY_URL/productpage"; done

```





**Cleanup**

```bash
samples/bookinfo/platform/kube/cleanup.sh

kubectl get virtualservices   #-- there should be no virtual services
kubectl get destinationrules  #-- there should be no destination rules
kubectl get gateway           #-- there should be no gateway
kubectl get pods              #-- the Bookinfo pods should be deleted

```





