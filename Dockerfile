FROM kubemq/gobuilder-ubuntu as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq-io/kubemq-sources
WORKDIR $GOPATH/github.com/kubemq-io/kubemq-sources
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags container -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.version=$VERSION" -o kubemq-sources-run .
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER KubeMQ info@kubemq.io
LABEL name="KubeMQ Sources Connectors" \
      maintainer="info@kubemq.io" \
      vendor="kubemq.io" \
      version="0.5.12" \
      release="stable" \
      summary="KubeMQ Sources connects external systems and cloud services to KubeMQ Message Broker" \
      description=""
COPY licenses /licenses
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir -p /opt/mqm/lib64
COPY --from=builder /opt/mqm/lib64 /opt/mqm/lib64
RUN mkdir /kubemq-connector
COPY --from=builder $GOPATH/github.com/kubemq-io/kubemq-sources/kubemq-sources-run ./kubemq-connector
COPY --from=builder $GOPATH/github.com/kubemq-io/kubemq-sources/default_config.yaml ./kubemq-connector/config.yaml
RUN chown -R 1001:root  /kubemq-connector && chmod g+rwX  /kubemq-connector
WORKDIR kubemq-connector
USER 1001
CMD ["./kubemq-sources-run"]
