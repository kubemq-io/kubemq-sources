helm repo add bitnami https://charts.bitnami.com/bitnami
helm install rabbitmq bitnami/rabbitmq --set auth.username=rabbitmq --set auth.password=rabbitmq
helm delete rabbitmq
kubectl port-forward --namespace default svc/rabbitmq 5672:5672
kubectl port-forward --namespace default svc/rabbitmq 15672:15672
