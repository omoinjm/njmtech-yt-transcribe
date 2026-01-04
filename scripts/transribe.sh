#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$(dirname "$SCRIPT_DIR")"

YT_LINK="https://www.youtube.com/watch?v=rdWZo5PD9Ek"
OUTPUT_DIR="/tmp/njmtech-yt-transcribe/"

./yt-transcribe -url $YT_LINK -output $OUTPUT_DIR
