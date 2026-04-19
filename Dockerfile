# Use a multi-stage build to create a lean final image.
# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

# Install necessary dependencies: git for fetching Go modules,
# and build-base for cgo if any dependencies require it.
RUN apk add --no-cache git build-base

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first,
# leveraging Docker's layer caching. NOTE: when enabling build-tags that pull
# additional imports, it's safer to run 'go mod download' after copying the full
# source so the module graph includes all build-tagged files.
COPY go.mod go.sum ./

# Copy the rest of the application's source code.
COPY . .

# Ensure module dependencies are downloaded and go.sum is populated inside the
# build context. This avoids missing go.sum entries when building with tag
# variants (e.g., 'infisical').
RUN go mod download

# Build the Go application.
# Build with the `infisical` tag when INFISICAL_ENABLED is true at build time.
# -ldflags="-w -s" strips debug information, reducing the binary size.
# CGO_ENABLED=0 disables cgo, creating a static binary.
ARG INFISICAL_ENABLED=false
RUN if [ "$INFISICAL_ENABLED" = "true" ]; then \
      CGO_ENABLED=0 go build -tags=infisical -ldflags="-w -s" -o /yt-transcribe . ; \
    else \
      CGO_ENABLED=0 go build -ldflags="-w -s" -o /yt-transcribe . ; \
    fi

# Stage 2: Create the final, minimal image
FROM alpine:3.21

# Install runtime dependencies: ffmpeg, curl, cmake, build-base.
# cmake and build-base are needed for building whisper.cpp.
RUN apk add --no-cache ffmpeg curl cmake build-base git

# Install yt-dlp from a pinned release for reproducible builds.
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/download/2025.10.22/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp

# Install and build whisper.cpp from a pinned release tag.
RUN git clone --depth 1 --branch v1.8.4 https://github.com/ggml-org/whisper.cpp.git /whisper.cpp && \
    cd /whisper.cpp && \
    sh ./models/download-ggml-model.sh base.en && \
    cmake -B build -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build -j && \
    cp build/bin/whisper-cli /usr/local/bin/whisper-cli

# Copy the built Go binary from the builder stage.
COPY --from=builder /yt-transcribe /usr/local/bin/

# Set the entrypoint for the container.
ENTRYPOINT ["/usr/local/bin/yt-transcribe"]
