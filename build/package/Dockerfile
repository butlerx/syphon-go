FROM golang:1.14 as builder

ARG VERSION="development"
ENV GO_PATH="/go"

COPY cmd internal /src/
WORKDIR /src

ENV GOOS=linux
RUN go get -v -d ./... && \
    go test -v ./... && \
    go build -i -v -o syphon.bin \
      -a -ldflags "-X main.version=$VERSION" \
      ./cmd/syphon

# Build the actual container
FROM scratch
LABEL maintainer="Cian Butler<butlerx@notthe.cloud>"

COPY --from=builder /src/syphon.bin /bin/syphon
ADD config/config.toml /bin/config.toml

EXPOSE 2003
LABEL VERSION=$VERSION

WORKDIR /bin
ENTRYPOINT ["/bin/syphon"]
CMD ["--config", "config.toml"]
