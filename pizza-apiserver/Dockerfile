FROM golang:1.20 as build
WORKDIR /go/src/github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pizza-apiserver cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver /
ENTRYPOINT ["/pizza-apiserver"]