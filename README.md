# co-type

A co-op typing game written in Go.

## Configuration

Each service is configured through environment variables.

### Broker

| Variable | Default | Description |
|---|---|---|
| `DISCOVERY_PORT` | `50051` | Port on which the broker listens for client discovery requests |
| `CONTROL_PORT` | `50052` | Port on which the broker listens for server control messages (gRPC `ControlService`) |
| `LOG_LEVEL` | `INFO` | Log verbosity (`DEBUG`, `INFO`, `WARN`, `ERROR`) |

### Server

| Variable | Required | Description |
|---|---|---|
| `SERVER_NAME` | true | Unique name for this server instance. The last character must be a digit (e.g. `gameserver-0`). |
| `SERVER_ADDR` | true | Address that clients will use to reach this server |
| `SERVER_PORT` | true | Base port number. The actual listening port is `SERVER_PORT + last_digit(SERVER_NAME)` (e.g. base `30100` + last_digit `1` → `30101`) |
| `GAME_PORT` | true | Port used for in-game gRPC communication (lobby service) |
| `CONTROL_ADDR` | true | Hostname or IP of the broker's control interface |
| `CONTROL_PORT` | true | Port of the broker's `ControlService` |
| `LOG_LEVEL` | false | Log verbosity (`DEBUG`, `INFO`, `WARN`, `ERROR`) |

### Client

| Variable | Required | Description |
|---|---|---|
| `DISCOVERY_ADDR` | true | Hostname or IP of the broker's discovery interface |
| `DISCOVERY_PORT` | true | Port of the broker's `DiscoveryService` |
| `LOG_LEVEL` | false | Log verbosity (`DEBUG`, `INFO`, `WARN`, `ERROR`). Logs always go to `debug.log` |

---

## Running with Kubernetes (kind)

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/)

### Deploy

The default Helm values pull pre-built images from the GitHub Container Registry (`ghcr.io/sceredi/co-type-*:latest`), so no local build step is required:

```bash
# Create the kind cluster and deploy everything
make deploy

# Check that all pods are running
make status
```

If you want to build and load local images instead (e.g. to test local changes):

```bash
make up   # docker-build + deploy + status
```

### Get service addresses

```bash
make addresses
```

This prints the NodePort addresses you can point the client at:

```
Broker DiscoveryService : <node-ip>:30051
GameServer-0            : <node-ip>:30100
GameServer-1            : <node-ip>:30101
GameServer-2            : <node-ip>:30102
```

### Run the client

```bash
DISCOVERY_ADDR=<node-ip> DISCOVERY_PORT=30051 go run ./client/cmd/app/main.go
```

### Useful commands

| Command | Description |
|---|---|
| `make deploy` | Deploy (or upgrade) the Helm release to the kind cluster |
| `make status` | Show pod and service status |
| `make logs-broker` | Tail broker logs |
| `make logs-servers` | Tail all game server logs |
| `make addresses` | Print NodePort addresses |
| `make set-log-level LOG_LEVEL=DEBUG` | Change log level on all running pods |
| `make delete` | Uninstall the Helm release and delete the namespace |
