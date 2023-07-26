FROM gcr.io/distroless/static-debian11

LABEL maintainer="zhangguanzhang@qq.com"

WORKDIR /opt/harbor-cleaner

COPY ./bin/harbor-cleaner cleaner
COPY ./cmd/config.yaml config.yaml

ENTRYPOINT ["./cleaner"]
