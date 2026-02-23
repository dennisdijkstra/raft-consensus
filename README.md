# Go Raft Project

Distributed key-value store implementation using the Raft consensus algorithm.

## Tech Stack

- Go
- HashiCorp Raft
- Gorilla Mux
- MDB Store

## Project Structure

- `main.go` – application entrypoint with Raft setup and HTTP server

## Getting Started

1. Install Go.
2. Install dependencies:

```bash
go mod download
```

3. Run a node:

```bash
go run . <node-id>
```

Example:
```bash
go run . node1
```

## Development

- The application uses Raft consensus to maintain a distributed key-value store
- Data is persisted in MDB format with snapshot support
- HTTP API provides endpoints for key-value operations and cluster management

## API Endpoints

- `GET /{key}` – retrieve value for a key
- `PUT /` – set a key-value pair
- `GET /nodes/list` – list cluster nodes
- `DELETE /nodes/{id}` – remove a node from the cluster
