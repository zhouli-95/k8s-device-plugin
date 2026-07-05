FROM golang:1.24 AS builder

WORKDIR /workspace

ARG GOPROXY

COPY . .

RUN make build


FROM ubuntu:24.04

WORKDIR /

ARG VERSION="unknown"
ARG GIT_COMMIT="unknown"

LABEL name="Device Plugin"
LABEL version=${VERSION}
LABEL commit=${GIT_COMMIT}
LABEL summary="Device plugin for Kubernetes"
LABEL description="Device plugin for Kubernetes"

COPY --from=builder /workspace/out/ .

CMD ["/device-plugin"]
