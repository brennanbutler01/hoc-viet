# Translation & Vocabulary API

A REST API for translating words and managing vocabulary collections. Built with Go using the Huma v2 framework.

## Features

- **Translation Service** - Translate words from English to Vietnamese using the MyMemory API
- **Vocabulary Management** - Store and retrieve studied words locally
- **Auto-generated Documentation** - Interactive OpenAPI/Swagger documentation

## API Documentation

### Translation Endpoints

| Method | Endpoint              | Description                                           |
| ------ | --------------------- | ----------------------------------------------------- |
| `GET`  | `/translation/{word}` | Get translation for a word from English to Vietnamese |

### Vocabulary Endpoints

| Method | Endpoint | Description                         |
| ------ | -------- | ----------------------------------- |
| `POST` | `/words` | Add a new word to your vocabulary   |
| `GET`  | `/words` | Retrieve all saved vocabulary words |

## Quick Start

### Prerequisites

- Go 1.19 or higher
- Internet connection for translation services

### Installation and Setup

1. Clone the repository:

```bash
git clone <your-repo-url>
cd translation-vocabulary-api
```

2. Initialize the Go module:

```bash
go mod init your-project
go mod tidy
```

3. Start the server:

```bash
go run .
```

The API will be available at `http://localhost:8888`. Interactive documentation can be found at `http://localhost:8888/docs`.

## Usage Examples

### Translate a Word

```bash
curl http://localhost:8888/translation/hello
```

Response:

```json
{
  "responseData": {
    "translatedText": "xin chào",
    "match": 1.0
  },
  "quotaFinished": false,
  "responseStatus": 200,
  "matches": [...]
}
```

### Add Word to Vocabulary

```bash
curl -X POST http://localhost:8888/words \
  -H "Content-Type: application/json" \
  -d '{
    "word": "hello",
    "translation": "xin chào"
  }'
```

Response:

```json
{
	"word": "hello",
	"translation": "xin chào",
	"dateCreated": "2025-09-17T10:30:00Z"
}
```

### Get Vocabulary List

```bash
curl http://localhost:8888/words
```

Response:

```json
[
	{
		"word": "hello",
		"translation": "xin chào",
		"dateCreated": "2025-09-17T10:30:00Z"
	},
	{
		"word": "world",
		"translation": "thế giới",
		"dateCreated": "2025-09-17T10:35:00Z"
	}
]
```
