# gRPC Sensor Server

A lightweight gRPC server for simulating and handling sensor data communication. Built with Golang.

## 🚀 Features
- Unary, Server Streaming, Client Streaming, and Bidirectional Streaming RPCs
- Handles simulated sensor data
- Metadata-based authentication
- Dockerized for easy deployment

## 🐳 Docker Usage

### Pull from Docker Hub
```bash
docker pull udmop/grpc-sensor-server
```

### Run the container
```bash
docker run -p 50051:50051 udmop/grpc-sensor-server
```

## 🧠 Requirements for Client
- `.proto` file to generate client stubs
- gRPC client (Go / Python / Node / etc.)
- Use metadata headers like:
  - `authorization: Bearer secret123`
  - `client-id: grpc-client-007`

## 📁 Proto File
Share or include the `sensor.proto` file so clients can generate gRPC code.

## 📦 Build Locally
```bash
docker build -t grpc-sensor-server .
docker run -p 50051:50051 grpc-sensor-server
```

