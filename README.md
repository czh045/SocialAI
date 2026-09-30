# Social AI

This project follows the Social AI Coding Pad lessons:

- React frontend: sign up, sign in, JWT-based requests, image/video collection, search, upload, deletion, and AI image generation.
- Go backend: Gorilla Mux routes, Elasticsearch indexes/search, JWT authentication, Google Cloud Storage integration, and Google App Engine deployment files.

## Local development

1. Start Elasticsearch from this directory:

```powershell
docker compose up -d
```

2. Install Go 1.25 or newer, then start the backend:

```powershell
cd Z:\Program\java\Social\socialai
go mod tidy
go run .
```

3. Install and start the React app:

```powershell
cd Z:\Program\java\Social\social-ai
npm.cmd install
npm.cmd start
```

The frontend runs on `http://localhost:3000` and the backend runs on `http://localhost:8080`.

If OnlineOrder is already using port `8080`, keep it running and start Social AI on `8081` instead:

```powershell
cd Z:\Program\java\Social\socialai
$env:PORT = "8081"
$env:SOCIALAI_PUBLIC_BASE_URL = "http://localhost:8081/uploads"
go run .

cd Z:\Program\java\Social\social-ai
$env:REACT_APP_API_BASE_URL = "http://localhost:8081"
npm.cmd start
```

## Teacher cloud configuration

The course deploys the Go backend to Google App Engine, connects it to Elasticsearch on GCE, and stores uploaded media in Google Cloud Storage. The cloud template is at `socialai/conf/deploy.gcp.example.yml`.

For safety, do not place real GCP credentials or OpenAI keys in Git. Copy the examples into ignored local configuration files and replace only the placeholder values.

## AI Image Generation

The image-generation button is available after signing in. The browser sends the
prompt to `POST /generate-image`; the Go backend calls OpenAI, so the API key is
never embedded in the React build.

Set these variables in the PowerShell window that starts the Go backend:

```powershell
$env:OPENAI_API_KEY = "your-key"
$env:OPENAI_IMAGE_MODEL = "gpt-image-2"
go run .
```

Without `OPENAI_API_KEY`, all Social AI features except AI image generation keep
working. Image generation is a paid cloud request, so leave the variable unset
when you are not actively testing it.

## Local Safety Defaults

- Local uploads are limited to 20 MB each.
- Deleting a post also removes its local upload or GCS media object.
- The backend accepts browser requests only from `http://localhost:3000` by default.
- Set `SOCIALAI_ALLOWED_ORIGINS` to a comma-separated list when deploying a frontend elsewhere.
- `socialai/conf/deploy.yml`, `socialai/uploads/`, and frontend `.env` files are ignored by Git.
