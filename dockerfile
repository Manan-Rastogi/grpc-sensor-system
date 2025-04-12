# 1. Base Image
# 2. Set Working Directory
# 3. Copy Go Modules
# 4. Download Dependencies
# 5. Copy Source Code
# 6. Build the Go App
# 7. Command to Run

FROM golang:1.24 AS builder

WORKDIR /app

COPY ./go.mod ./ 
COPY ./go.sum ./ 

RUN go mod download

COPY . .

RUN go build -o main ./server

# Multistage
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/main . 
EXPOSE 50051

CMD ["./main"]


# docker build -t grpc-server-build .
# docker run -p 50051:50051 grpc-server-build

# docker login
# docker tag grpc-server-build udmop/grpc-server:latest
# docker push udmop/grpc-server:latest