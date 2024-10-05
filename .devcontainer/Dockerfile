ARG DEBIAN_VERSION="bookworm"

FROM golang:${DEBIAN_VERSION} AS go

FROM mcr.microsoft.com/vscode/devcontainers/base:1-${DEBIAN_VERSION}

RUN mkdir /go \
    && groupadd -r golang \
    && usermod -a -G golang vscode \
    && chown -R vscode:golang /go \
    && chmod -R g+r+w+s /go

COPY --from=go --chown=vscode:golang --chmod=g+r+w+s /usr/local/go /usr/local/go

ENV GOPATH=/go \
    CGO_ENABLED=1 \
    GO111MODULE=on
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

USER vscode
ENV GOPATH=/go \
    CGO_ENABLED=1 \
    GO111MODULE=on
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

RUN echo "alias ll='ls -laF'" > ~/.bash_aliases