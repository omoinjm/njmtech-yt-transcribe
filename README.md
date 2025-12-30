# YouTube Video Transcriber CLI (yt-transcribe)

`yt-transcribe` is a command-line interface (CLI) tool written in Go that allows you to transcribe the audio from YouTube videos into text. It leverages the powerful `yt-dlp` tool to extract audio and is designed with an extensible architecture to integrate with various transcription services.

## ✨ Features

*   **YouTube Audio Extraction:** Downloads the audio stream from any valid YouTube video URL.
*   **Text Transcription:** Processes the extracted audio to generate a text transcription (currently mocked, but extensible).
*   **Customizable Output:** Allows specifying an output directory for the generated transcription file.
*   **Temporary File Management:** Automatically cleans up downloaded audio files after transcription.
*   **Modular Design:** Built with SOLID principles, making it easy to swap out different audio downloaders or transcription services.

## 🚀 Getting Started

### Prerequisites

Before you can build and run `yt-transcribe`, you'll need the following:

1.  **Go:**
    *   Ensure you have Go version `1.25.5` or newer installed on your system.
    *   Download and installation instructions can be found at: [https://golang.org/doc/install](https://golang.org/doc/install)

2.  **`yt-dlp`:**
    *   `yt-dlp` is an essential external tool that `yt-transcribe` uses to download audio from YouTube videos. It must be installed and available in your system's `PATH`.
    *   **Installation on Linux/macOS:**
        ```bash
        sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
        sudo chmod a+rx /usr/local/bin/yt-dlp # Give execute permissions
        ```
    *   **Installation on Windows:**
        *   Download `yt-dlp.exe` from the [official GitHub releases page](https://github.com/yt-dlp/yt-dlp/releases).
        *   Place the `yt-dlp.exe` file in a directory that is included in your system's `PATH` environment variable (e.g., `C:\Windows`, or a custom directory you've added to PATH).

### Building the Tool

1.  **Clone the repository** (if you haven't already):
    ```bash
    git clone https://github.com/your-username/yt-transcribe.git # Replace with actual repo if applicable
    cd yt-transcribe
    ```
    *(Note: Assuming the current directory is the project root)*

2.  **Build the executable:**
    Navigate to the project's root directory in your terminal and run:
    ```bash
    go build -o yt-transcribe
    ```
    This command will compile your Go code and create an executable file named `yt-transcribe` (or `yt-transcribe.exe` on Windows) in the current directory.

## 💡 Usage

To run the `yt-transcribe` tool, you need to provide a YouTube video URL using the `-url` flag. You can optionally specify a custom output directory using the `-output` flag.

```bash
./yt-transcribe -url <YOUTUBE_VIDEO_URL> [-output <OUTPUT_DIRECTORY>]
```

*   Replace `<YOUTUBE_VIDEO_URL>` with the actual link to the YouTube video you want to transcribe.
*   Replace `<OUTPUT_DIRECTORY>` with the path where you want the transcription `.txt` file to be saved. If omitted, the transcription will be saved in your system's temporary directory.

### Examples:

1.  **Transcribe a video and save to the default temporary directory:**
    ```bash
    ./yt-transcribe -url https://www.youtube.com/watch?v=dQw4w9WgXcQ
    ```

2.  **Transcribe a video and save to a specific directory:**
    ```bash
    ./yt-transcribe -url https://www.youtube.com/watch?v=your_video_id -output ~/transcriptions
    ```

## 📝 Note on Transcription Implementation

The current `yt-transcribe` tool uses a **mock `OpenAITranscriber`**. This means that while the audio downloading and file handling logic is fully functional, the transcription itself will return placeholder text rather than actual content from a transcription service.

To enable real transcription, you would need to modify the `Transcribe` method within `pkg/transcriber.go`. This typically involves:

1.  **Obtaining an API Key:** Get an API key from your chosen transcription service (e.g., OpenAI's Whisper API, Google Cloud Speech-to-Text, AssemblyAI).
2.  **Integrating the API:** Use the service's Go SDK or make direct HTTP requests to their API endpoint.
3.  **Handling Audio Upload:** Implement the logic to upload the extracted audio file to the transcription service.
4.  **Parsing Response:** Process the service's response to extract the transcribed text.

This modular design allows you to easily switch to any transcription service by implementing the `transcriber.Transcriber` interface.

## 🤝 Contributing

(If this were an open-source project, you'd find contribution guidelines here.)

## 📄 License

(If this were an open-source project, you'd find license information here.)
