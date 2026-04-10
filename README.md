# yt-transcribe

A CLI tool written in Go that downloads audio from YouTube, Instagram, and other platforms supported by `yt-dlp`, transcribes it using `whisper.cpp`, and uploads the transcript (SRT format with timestamps) to Vercel Blob storage. It can process a single URL, pull the next job from a Postgres database, or reprocess all existing records.

## Features

- Downloads audio via `yt-dlp` and converts to WAV with `ffmpeg`
- Transcribes using `whisper.cpp` — outputs SRT files with timestamps
- Uploads transcripts to Vercel Blob storage
- Three run modes: single URL, DB-driven, and reprocess-all
- Idle-safe DB connection (uses `pgxpool` — survives Neon's connection timeouts during long jobs)

---

## Environment Variables

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

| Variable | Required | Description |
|---|---|---|
| `WHISPER_MODEL_PATH` | ✅ | Path to the `ggml-*.bin` model file |
| `VERCEL_BLOB_API_URL` | ✅ | Upload endpoint for your Blob API |
| `VERCEL_BLOB_API_TOKEN` | ✅ | Auth token for the Blob API |
| `POSTGRES_URL` | `-db` / `-reprocess-all` only | Neon / Postgres connection string |
| `DOCKERHUB_USERNAME` | Docker Compose only | Your Docker Hub username (resolves the image name) |

---

## Usage

### Flags

```
-url <URL>        Transcribe a single video URL
-output <dir>     Directory for temporary audio files (default: /tmp)
-db               Fetch and process the next unprocessed URL from the database
-reprocess-all    Reprocess every record in the database (overwrites existing transcripts)
```

### Examples

**Single URL:**
```bash
./yt-transcribe -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

**Next unprocessed item from DB:**
```bash
./yt-transcribe -db
```

**Reprocess all records:**
```bash
./yt-transcribe -reprocess-all
```

---

## Running with Docker

The pre-built image is published to Docker Hub on every merge to `main`. All dependencies (`yt-dlp`, `ffmpeg`, `whisper.cpp`, model) are bundled inside the image.

### `docker run`

```bash
docker run --rm --env-file .env \
  your-dockerhub-username/njmtech-yt-transcribe:latest -db
```

Replace `-db` with any valid flag combination (e.g. `-url "https://..."`, `-reprocess-all`).

### `docker compose`

Ensure `DOCKERHUB_USERNAME` is set in your `.env` file, then:

```bash
# Process next unprocessed DB item (default command in docker-compose.yml)
docker compose run --rm yt-transcribe

# Override command
docker compose run --rm yt-transcribe -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
docker compose run --rm yt-transcribe -reprocess-all
```

---

## VPS Deployment (recommended for long videos)

For videos longer than a few minutes, run the worker on a dedicated VPS instead of CI. A helper script handles the full setup:

```bash
# On your VPS — clone the repo and run once
git clone https://github.com/omoinjm/njmtech-yt-transcribe.git
cd njmtech-yt-transcribe
bash scripts/setup-vps.sh
```

The script will:
1. Install Docker (if not already present)
2. Copy `docker-compose.yml` to `/opt/yt-transcribe/`
3. Create `/opt/yt-transcribe/.env` from `.env.example` (edit this file with your secrets)
4. Pull the latest image from Docker Hub
5. Register a cron job that runs `./yt-transcribe -db` every 30 minutes

Logs are written to `/var/log/yt-transcribe.log`.

---

## Building from Source

**Prerequisites:** Go 1.22+, `yt-dlp`, `ffmpeg`, `whisper-cli` (from [whisper.cpp](https://github.com/ggml-org/whisper.cpp))

```bash
git clone https://github.com/omoinjm/njmtech-yt-transcribe.git
cd njmtech-yt-transcribe
go build -o yt-transcribe .
```

**Run tests:**
```bash
go test ./...
```

---

## License

See [LICENSE](LICENSE).

