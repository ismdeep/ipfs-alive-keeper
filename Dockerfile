FROM golang as builder

RUN mkdir -p /src/ipfs-alive-keeper
WORKDIR      /src/ipfs-alive-keeper
ADD .  .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ipfs-alive-keeper

FROM alpine:latest
MAINTAINER "L. Jiang <l.jiang.1024@gmail.com>"
COPY --from=builder /ipfs-alive-keeper /
RUN apk add --no-cache tzdata
RUN chmod +x /ipfs-alive-keeper
CMD ["/ipfs-alive-keeper", "-c", "/config.toml"]