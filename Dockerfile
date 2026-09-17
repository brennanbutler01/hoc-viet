FROM golang:1.25@sha256:699337d620559a59b4a2bb298ad59611e535d2ee755a34cf2d2a98f37578dc80 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /hoc-viet . && mkdir /data
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /hoc-viet /hoc-viet
COPY --from=build --chown=65532:65532 /data /data
USER 65532:65532
WORKDIR /data
ENV LISTEN_ADDR=0.0.0.0:8888 VOCABULARY_FILE=/data/words.json
EXPOSE 8888
ENTRYPOINT ["/hoc-viet"]
