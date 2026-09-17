# Học Việt: translation and vocabulary API

A small Go service for translating English words into Vietnamese and saving a personal vocabulary collection. Built with Huma and Chi, with generated OpenAPI documentation.

This is a local, single-user project. It has no authentication or multi-user isolation and is not intended to be exposed directly to the internet. The separate Tofu.Vocab portfolio demo is at https://tofu-vocab-demo.vercel.app.

## Run locally

Requires Go 1.25 or newer. The repository already contains its module definition; do not run `go mod init` again.

```sh
go mod download
go run .
```

Open http://127.0.0.1:8888/docs for interactive API documentation.

| Method | Endpoint | Purpose |
| --- | --- | --- |
| GET | `/translation/{word}` | Translate an English word or short phrase to Vietnamese |
| POST | `/words` | Save a word and its translation |
| GET | `/words` | Retrieve the saved vocabulary |

```sh
curl http://127.0.0.1:8888/words
curl -X POST http://127.0.0.1:8888/words \
  -H 'Content-Type: application/json' \
  -d '{"word":"hello","translation":"xin chào"}'
curl http://127.0.0.1:8888/translation/hello
```

Translation sends the requested text to the public MyMemory service. Availability and quotas depend on that provider. Requests have a ten-second timeout, inherit client cancellation, reject unsuccessful or malformed provider responses, and limit response bodies to one MiB. The test suite uses local HTTP servers and does not contact MyMemory.

## Docker

```sh
docker build -t hoc-viet:local .
docker run --rm --name hoc-viet-local -p 127.0.0.1:5195:8888 hoc-viet:local
```

Open http://127.0.0.1:5195/docs. The container runs as an unprivileged user. Its vocabulary is disposable unless you explicitly mount persistent storage; removing the container deletes its saved words.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `LISTEN_ADDR` | `127.0.0.1:8888` | Listener address; the Docker image uses `0.0.0.0:8888` internally |
| `VOCABULARY_FILE` | `words.json` | JSON storage file; the Docker image uses `/data/words.json` |

The vocabulary directory must exist and be writable. One repository instance serializes access within the process. Writes stage a private file beside the destination, sync it, and atomically rename it into place. Existing malformed data is preserved and reported as an error. Multiple processes writing the same file are not supported; use a database if that becomes a requirement.

Vocabulary data, compiled binaries, environment files, and temporary build artifacts are ignored by Git. Do not put personal data or credentials in commits.

## Verify

```sh
go test -race -count=1 ./...
go vet ./...
go build -o /tmp/hoc-viet .
```

Tests cover query encoding, provider failures and oversized responses, cancellation/timeouts, concurrent saves, persistence, malformed-file preservation, and the real vocabulary HTTP handlers. Storage failures return server errors rather than being mislabeled as invalid user input. Server startup failures are reported, request timeouts are configured, and shutdown drains active requests.

Without a local Go installation:

```sh
docker run --rm -v "$PWD:/src" -w /src golang:1.25 \
  go test -race -count=1 ./...
```

Remaining limitations: no authentication, no request rate limiting, no multi-process storage coordination, no migration to a database, and no verification against the live translation provider. This service is not deployed on Vercel; its file persistence requires a different hosting design before public use.

Recovery verification (September 17, 2026): all three Go packages passed race-enabled tests; `go vet` and compilation passed. The Docker image built and its running container passed empty-list, save, persisted-read, and documentation checks. The disposable verification container was stopped and removed.
