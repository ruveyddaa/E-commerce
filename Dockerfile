# ----- Build Stage -----
    FROM golang:1.24.2-alpine AS builder

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    

    RUN go build -o main main.go
    
    # ----- Run Stage -----
    FROM alpine:latest
    WORKDIR /
    
    RUN apk --no-cache add ca-certificates
    
    COPY --from=builder /app/main .
    

    ENTRYPOINT ["/main"]