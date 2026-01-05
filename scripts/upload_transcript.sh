#!/bin/bash
set -e

AUDIO_FILE="/tmp/njmtech-yt-transcribe/audio.wav"
OUTPUT_PREFIX="/tmp/njmtech-yt-transcribe/transcript"
WHISPER_CONTEXT="/whisper.cpp/models/ggml-base.en.bin"

whisper-cli \
  -m "$WHISPER_CONTEXT" \
  -f "$AUDIO_FILE" \
  --output-txt \
  --output-file "$OUTPUT_PREFIX" \
  --no-prints
