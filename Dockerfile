# Use a multi-stage build to create a lean final image.
# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Install necessary dependencies: git for fetching Go modules,
# and build-base for cgo if any dependencies require it.
RUN apk add --no-cache git build-base

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first,
# leveraging Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Build the Go application.
# -ldflags="-w -s" strips debug information, reducing the binary size.
# CGO_ENABLED=0 disables cgo, creating a static binary.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /yt-transcribe .

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Install runtime dependencies: ffmpeg, curl, cmake, build-base.
# cmake and build-base are needed for building whisper.cpp.
RUN apk add --no-cache ffmpeg curl cmake build-base git

# Install yt-dlp from the latest release on GitHub.
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp

# Install and build whisper.cpp, and download the base.en model.
RUN git clone https://github.com/ggerganov/whisper.cpp.git /whisper.cpp && \
    cd /whisper.cpp && \
    sh ./models/download-ggml-model.sh base.en && \
    cmake -B build -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build -j && \
    cp build/bin/whisper-cli /usr/local/bin/whisper-cli

# Copy the built Go binary from the builder stage.
COPY --from=builder /yt-transcribe /usr/local/bin/

# Set the entrypoint for the container.
ENTRYPOINT ["/usr/local/bin/yt-transcribe"]
