FROM golang:1.15 as builder

ARG VERSION="development"
ENV GO_PATH="/go"

COPY . /src/
WORKDIR /src

ENV GOOS=linux
RUN touch .image-id && \
      make build

# Build the actual container
FROM scratch
LABEL maintainer="Cian Butler<butlerx@notthe.cloud>"

COPY --from=builder /src/syphon.bin /bin/syphon
ADD configs/config.toml /bin/config.toml

EXPOSE 2003
LABEL VERSION=$VERSION

WORKDIR /bin
ENTRYPOINT ["/bin/syphon"]
CMD ["--config", "config.toml"]
