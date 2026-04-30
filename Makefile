LOG_LEVEL ?= INFO
KUBECONFIG ?= /etc/rancher/k3s/k3s.yaml
NODE_IP = $(shell kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./common/proto/**/*.proto

.PHONY: build
build: proto
	go build ./client/... ./server/... ./broker/...

.PHONY: test
test:
	go test -v -race ./common/... ./client/... ./server/... ./broker/...

.PHONY: tidy
tidy:
	cd common && go mod tidy
	cd client && go mod tidy
	cd server && go mod tidy
	cd broker && go mod tidy

.PHONY: docker-build
docker-build:
	docker build --build-arg MODULE=broker -t co-type/broker:latest .
	docker build --build-arg MODULE=server -t co-type/server:latest .
	docker save co-type/broker:latest | sudo k3s ctr images import -
	docker save co-type/server:latest | sudo k3s ctr images import -

.PHONY: deploy
deploy:
	helm upgrade --install cotype-release ./k8s \
		--namespace cotype \
		--create-namespace \
		--set global.nodeIp=$(NODE_IP) \
		--set global.logLevel=$(LOG_LEVEL)

.PHONY: delete
delete:
	helm uninstall cotype-release -n cotype
	kubectl delete namespace cotype

.PHONY: status
status:
	kubectl get pods,services -n cotype

.PHONY: logs-broker
logs-broker:
	kubectl logs -n cotype -l app=broker -f

.PHONY: logs-servers
logs-servers:
	kubectl logs -n cotype -l app=gameserver -f --max-log-requests 10

.PHONY: addresses
addresses:
	@echo "Broker GatewayService : $(NODE_IP):30051"
	@echo "GameServer-0          : $(NODE_IP):30100"
	@echo "GameServer-1          : $(NODE_IP):30101"
	@echo "GameServer-2          : $(NODE_IP):30102"

.PHONY: set-log-level
set-log-level:
	kubectl set env statefulset -n cotype --all LOG_LEVEL=$(LOG_LEVEL)
	kubectl set env deployment -n cotype --all LOG_LEVEL=$(LOG_LEVEL)

.PHONY: up
up: docker-build deploy status
