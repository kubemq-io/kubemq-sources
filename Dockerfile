FROM kubemq/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq-hub/kubemq-sources
WORKDIR $GOPATH/github.com/kubemq-hub/kubemq-sources
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s -X main.version=$VERSION" -o kubemq-sources-run .
FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER KubeMQ info@kubemq.io
LABEL name="KubeMQ Target Connectors" \
      maintainer="info@kubemq.io" \
      vendor="" \
      version="" \
      release="" \
      summary="" \
      description=""
COPY licenses /licenses
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir /kubemq-sources
COPY --from=builder $GOPATH/github.com/kubemq-hub/kubemq-sources/kubemq-sources-run ./kubemq-sources
RUN chown -R 1001:root  /kubemq-sources && chmod g+rwX  /kubemq-sources
WORKDIR kubemq-sources
USER 1001
CMD ["./kubemq-sources-run"]
