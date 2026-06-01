LOG_LEVEL ?= INFO
KIND_CLUSTER_NAME ?= cotype
KIND_CONFIG ?= k8s/cluster.yaml
KUBECONFIG ?= $(HOME)/.kube/config
export KUBECONFIG

.PHONY: kind-cluster
kind-cluster:
	@if ! kind get clusters | grep -qx "$(KIND_CLUSTER_NAME)"; then \
		kind create cluster --name $(KIND_CLUSTER_NAME) --config $(KIND_CONFIG); \
	fi

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
docker-build: kind-cluster
	docker build --build-arg MODULE=broker -t co-type/broker:latest .
	docker build --build-arg MODULE=server -t co-type/server:latest .
	kind load docker-image --name $(KIND_CLUSTER_NAME) co-type/broker:latest
	kind load docker-image --name $(KIND_CLUSTER_NAME) co-type/server:latest

.PHONY: deploy
deploy: kind-cluster
	@node_ip="$$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')" && \
	helm upgrade --install cotype-release ./k8s \
		--namespace cotype \
		--create-namespace \
		--set global.nodeIp="$$node_ip" \
		--set global.logLevel=$(LOG_LEVEL)

.PHONY: delete
delete:
	helm uninstall cotype-release -n cotype
	kubectl delete namespace cotype

.PHONY: status
status: kind-cluster
	kubectl get pods,services -n cotype

.PHONY: logs-broker
logs-broker: kind-cluster
	kubectl logs -n cotype -l app=broker -f

.PHONY: logs-servers
logs-servers: kind-cluster
	kubectl logs -n cotype -l app=gameserver -f --max-log-requests 10

.PHONY: addresses
addresses: kind-cluster
	@node_ip="$$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')" && \
	echo "Broker GatewayService : $$node_ip:30051" && \
	echo "GameServer-0          : $$node_ip:30100" && \
	echo "GameServer-1          : $$node_ip:30101" && \
	echo "GameServer-2          : $$node_ip:30102"

.PHONY: set-log-level
set-log-level: kind-cluster
	kubectl set env statefulset -n cotype --all LOG_LEVEL=$(LOG_LEVEL)
	kubectl set env deployment -n cotype --all LOG_LEVEL=$(LOG_LEVEL)

.PHONY: up
up:
	$(MAKE) docker-build
	$(MAKE) deploy
	$(MAKE) status
