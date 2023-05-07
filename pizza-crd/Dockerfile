FROM golang:1.20 as build
WORKDIR /go/src/github.com/ZhengHe-MD/programming-kubernetes/pizza-crd
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/pizza-crd-webhook

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pizza-crd-webhook /
ENTRYPOINT ["/pizza-crd-webhook"]