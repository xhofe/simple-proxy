FROM alpine:3.18 as builder
LABEL stage=go-builder
WORKDIR /app/
COPY ./ ./
RUN apk add --no-cache go; \
    go mod download; \
    go build -o bin/simple-proxy -ldflags "-s -w" .

FROM alpine:3.18
LABEL MAINTAINER="i@nn.ci"
VOLUME /opt/simple-proxy/data/
WORKDIR /opt/simple-proxy/
RUN apk add --no-cache ca-certificates; \
    mkdir -p data
COPY --from=builder /app/bin/simple-proxy ./
COPY config.json.example ./data/config.json
EXPOSE 3000
CMD [ "simple-proxy", "-conf", "data/config.json" ]