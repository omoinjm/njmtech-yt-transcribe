# Project: yt-transcribe - YouTube Video Transcriber CLI

## Project Overview

This project is a Command Line Interface (CLI) tool written in Go that transcribes YouTube videos into text. It operates in two main stages: first, it extracts the audio stream from a specified YouTube video, and then it processes this audio to generate a text transcription. The tool is designed with extensibility and maintainability in mind, adhering to modern software design principles.

## Main Technologies and Architecture

*   **Go (Golang):** The core programming language for the CLI tool.
*   **`flag` package:** Utilized for parsing command-line arguments, allowing users to specify the YouTube video URL and an optional output directory.
*   **`yt-dlp` (External Tool):** An essential dependency, `yt-dlp` is an external command-line program responsible for reliably downloading and extracting the audio track from YouTube videos. The Go application executes `yt-dlp` as a subprocess.
*   **Ollama Whisper:** The architecture is designed to integrate with a transcription service, using Ollama's Whisper model.

The project follows a modular architecture based on SOLID principles:
*   The `main` package orchestrates the overall workflow, handles CLI input, and performs dependency injection.
*   The `pkg/downloader` package defines the `YouTubeDownloader` interface and provides a concrete `YTDLPAudioDownloader` implementation that interfaces with `yt-dlp`.
*   The `pkg/transcriber` package defines the `Transcriber` interface and includes a `OllamaTranscriber` that would integrate with a real transcription service.

This design promotes loose coupling, making it straightforward to swap out different audio downloading mechanisms or transcription services without altering the core logic.

## Building and Running

To build and run this CLI tool, follow these steps:

1.  **Install Go:**
    Ensure you have Go version `1.25.5` or newer installed on your system. You can download it from [https://golang.org/doc/install](https://golang.org/doc/install).

2.  **Install `yt-dlp`:**
    The `yt-dlp` command-line tool is a prerequisite for audio downloading. Install it and ensure it's accessible in your system's PATH.
    *   **Linux/macOS Example:**
        ```bash
        sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
        sudo chmod a+rx /usr/local/bin/yt-dlp # Give execute permissions
        ```
    *   **Windows:** Download `yt-dlp.exe` from [https://github.com/yt-dlp/yt-dlp/releases](https://github.com/yt-dlp/yt-dlp/releases) and add its directory to your system's PATH.

3.  **Install `ffmpeg`:**
    `ffmpeg` is required by `yt-dlp` for post-processing audio, such as converting to MP3 format. Install it and ensure it's accessible in your system's PATH.
    *   **Linux (Debian/Ubuntu) Example:**
        ```bash
        sudo apt update
        sudo apt install ffmpeg
        ```
    *   **macOS (using Homebrew):**
        ```bash
        brew install ffmpeg
        ```
    *   **Windows:** Download `ffmpeg` from [https://ffmpeg.org/download.html](https://ffmpeg.org/download.html) and add its binaries directory to your system's PATH.

4.  **Install and run Ollama**
    Follow the instruction on [https://ollama.com/](https://ollama.com/) to install and run Ollama. Once Ollama is running, pull the `whisper` model:
    ```bash
    ollama pull whisper
    ```

5.  **Build the CLI Tool:**
    Navigate to the project's root directory (`/home/user/dev/github/projects/njmtech-yt-transcribe`) in your terminal and compile the application:
    ```bash
    go build -o yt-transcribe
    ```
    This will create an executable named `yt-transcribe` (or `yt-transcribe.exe` on Windows) in the current directory.

6.  **Run the Tool:**
    Execute the tool by providing a YouTube video URL using the `-url` flag. You can also specify a custom output directory with the `-output` flag.
    ```bash
    ./yt-transcribe -url <YOUTUBE_VIDEO_URL> [-output <OUTPUT_DIRECTORY>]
    ```
    *   **Example Usage:**
        ```bash
        ./yt-transcribe -url https://www.youtube.com/watch?v=dQw4w9WgXcQ
        ```
    *   The transcription will be saved as a `.txt` file. By default, it uses your system's temporary directory for output.

## Development Conventions

*   **Go Standard Formatting:** The codebase adheres to standard Go formatting practices, ensuring readability and consistency.
*   **SOLID Principles:** The design heavily utilizes interfaces (`YouTubeDownloader`, `Transcriber`) to promote modularity, extensibility, and testability, in line with SOLID principles (e.g., Dependency Inversion, Interface Segregation).
*   **Clear Comments:** Functions, interfaces, and complex logic sections are well-commented, explaining their purpose, usage, and any external dependencies or conceptual aspects (e.g., the mock transcriber).
*   **Robust Error Handling:** Errors are managed using Go's idiomatic error return values, and critical failures are handled with `log.Fatalf`.
*   **Temporary File Management:** Downloaded audio files are treated as temporary and are automatically cleaned up after transcription.

---
**Note on Transcription Implementation:** The current `OllamaTranscriber` uses the `whisper` model. For actual transcription, you would need to have Ollama installed and running with the `whisper` model pulled. The `OLLAMA_HOST` environment variable must be set to the URL of your Ollama server. By default it is `http://localhost:11434`.
