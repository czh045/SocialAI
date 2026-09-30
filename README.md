# SocialAI

A full-stack media-sharing application with a React client, a Go API, JWT
authentication, Elasticsearch search, and local or Google Cloud Storage media
handling. Users can create accounts, publish images or videos, search posts,
manage their own uploads, and optionally generate an image through a
server-side AI integration.

## Portfolio Snapshot

- **Product:** media-first social application with protected user workflows.
- **Backend focus:** Go HTTP services, JWT authorization, Elasticsearch, cloud
  storage abstraction, and configurable CORS.
- **Frontend focus:** React routing, authentication-aware API requests, upload
  interactions, responsive media browsing, and search.
- **Deployment path:** Dockerized Elasticsearch for local development and
  Google App Engine / Google Cloud Storage configuration templates.

## Features

- Account registration, sign-in, and JWT-protected API requests
- Image and video upload with a 20 MB local safety limit
- Media feed, search, and responsive gallery/lightbox browsing
- Ownership checks for deleting a user's own post
- Elasticsearch indexing and keyword search
- Local file storage for development and optional Google Cloud Storage for
  deployment
- Optional server-side AI image generation; browser code never contains the
  OpenAI API key
- Configurable allowed origins for deployment outside localhost

## Architecture

```text
React client
    |
    v
Go API (Gorilla Mux, JWT, CORS)
    |                 |                 |
    v                 v                 v
Elasticsearch     Media storage     Optional image API
                 local / GCS
```

## Project Structure

```text
SocialAI/
|-- social-ai/                   # React client
|   |-- src/                     # Views, components, API client
|   `-- .env.example             # Browser-safe configuration template
|-- socialai/                    # Go API
|   |-- handler/                 # HTTP handlers and auth checks
|   |-- model/                   # Data models and Elasticsearch mapping
|   |-- service/                 # Search and storage integrations
|   `-- conf/                    # Local and GCP deployment templates
|-- docker-compose.yml           # Elasticsearch for local development
`-- docs/                        # Implementation and learning notes
```

## Run Locally

### Prerequisites

- Go 1.25+
- Node.js 20+
- Docker Desktop

Start Elasticsearch from the repository root:

```powershell
docker compose up -d
```

Start the Go API in one terminal:

```powershell
cd socialai
go mod download
go run .
```

Start the React client in a second terminal:

```powershell
cd social-ai
npm install
npm start
```

The client runs on `http://localhost:3000` and the API runs on
`http://localhost:8080`.

To use another API port, set both the backend `PORT` and the frontend
`REACT_APP_API_BASE_URL` before starting the services.

## Optional Cloud and AI Features

The application works locally without paid services. Cloud storage and image
generation are optional integrations:

```powershell
# Set only in the terminal session that starts the Go API.
$env:OPENAI_API_KEY = "your-key"
$env:OPENAI_IMAGE_MODEL = "gpt-image-2"
go run .
```

Use the example deployment files in `socialai/conf/` for Google App Engine and
Google Cloud Storage. Never commit real cloud credentials, `.env` files, local
uploads, or API keys.

## Engineering Notes

- Media deletion is handled as an ownership-protected operation and removes the
  corresponding local or cloud object.
- Local development accepts only the local React origin by default; production
  origins are configured explicitly with `SOCIALAI_ALLOWED_ORIGINS`.
- The external image-generation call is kept on the server so secret keys are
  not shipped in the React bundle.
- Elasticsearch is isolated behind the backend API instead of being accessed
  directly by the browser.
