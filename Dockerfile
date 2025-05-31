FROM arm64v8/golang:1.24-bookworm AS builder

# Magic line, notice in use that the lib name is different!
RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu
# Add your app and do what you need to for dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -o app .

# Final stage - pick any old arm64 image you want
FROM shinsenter/scratch:latest

WORKDIR /root/
    
COPY --from=builder /go/app .

CMD ["./app", "--config.file=config.yaml"]