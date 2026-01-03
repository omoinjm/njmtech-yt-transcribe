#!/bin/bash
set -e

AUDIO_FILE="/tmp/njmtech-yt-transcribe/audio.wav"
OUTPUT_PREFIX="/tmp/njmtech-yt-transcribe/transcript"

whisper-cli \
  -f "$AUDIO_FILE" \
  --output-txt \
  --output-file "$OUTPUT_PREFIX" \
  --no-prints
