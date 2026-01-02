ARG BUILD_ARG_GO_VERSION=1.25
ARG BUILD_ARG_ALPINE_VERSION=3.22

FROM golang:${BUILD_ARG_GO_VERSION}-alpine${BUILD_ARG_ALPINE_VERSION} AS builder

ENV GOLANGCI_LINT_VERSION=2.6.2
ENV MOCKERY_VERSION=3.6.1

RUN apk add --update --no-cache git gcc libc-dev ca-certificates tar

ADD https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz /tmp

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/vektra/mockery/v3@v${MOCKERY_VERSION}

RUN apk add --no-cache curl unzip

ENV PATH="/usr/local/bin:/usr/local/include:$PATH"

RUN mkdir -p /tmp/utils && \
    tar -xzvf /tmp/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz -C /go/bin && \
    chmod +x /go/bin/*

WORKDIR /src

FROM golang:${BUILD_ARG_GO_VERSION}-alpine${BUILD_ARG_ALPINE_VERSION}

RUN apk add git

# Set all directories as safe
RUN git config --global --add safe.directory '*'

COPY --from=builder /go/bin/* /go/bin/

WORKDIR /src
