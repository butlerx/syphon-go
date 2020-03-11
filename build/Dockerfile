# Use a mutli-stage build pipeline to generate the executable
FROM golang:1.11

ARG VERSION="development"

ENV GO_PATH="/go"

ADD . /src
WORKDIR /src

RUN go test -v ./...

ENV GOOS=linux
RUN go build -i -v -o syphon.bin -a -ldflags "-X main.version=$VERSION" ./cmd/syphon

# Build the actual container
FROM scratch
LABEL maintainer="Cian Butler<butlerx@notthe.cloud>"

COPY --from=0 /src/syphon.bin /bin/syphon
ADD assets/config.toml /bin/config.toml

EXPOSE 2003

LABEL VERSION=$VERSION

WORKDIR /bin
ENTRYPOINT ["/bin/syphon"]
CMD ["--config", "config.toml"]
