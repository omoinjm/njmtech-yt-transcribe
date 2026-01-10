# Video Transcriber CLI (yt-transcribe)

`yt-transcribe` is a command-line interface (CLI) tool written in Go that allows you to download and transcribe audio from various video platforms. It leverages the powerful `yt-dlp` tool to extract audio and `ffmpeg` to convert it to WAV format.

## ✨ Features

- **Video Audio Extraction:** Downloads the audio stream from any valid video URL supported by `yt-dlp` (including YouTube, Instagram, and more).
- **WAV Format Output:** Converts downloaded audio to high-quality WAV format.
- **Date-Based Naming:** Saves audio files with a `video_YYYY-MM-DD.wav` naming convention.
- **Customizable Output Directory:** Allows specifying an output directory for the downloaded audio file.
- **Dependency Checks:** Automatically checks for `yt-dlp` and `ffmpeg` and prompts the user to install them if missing.
- **Modular Design:** Built with SOLID principles, making it easy to swap out different audio downloaders or transcription services (though transcription is currently disabled).

## 🚀 Getting Started

### Prerequisites

Before you can build and run `yt-transcribe`, you'll need the following:

1.  **Go:**
    - Ensure you have Go version `1.22` or newer installed on your system.
    - Download and installation instructions can be found at: [https://golang.org/doc/install](https://golang.org/doc/install)

2.  **`yt-dlp`:**
    - `yt-dlp` is an essential external tool that `yt-transcribe` uses to download audio from YouTube videos. The application explicitly checks for its presence. It must be installed and available in your system's `PATH`.
    - **Installation on Linux/macOS:**
      ```bash
      sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
      sudo chmod a+rx /usr/local/bin/yt-dlp # Give execute permissions
      ```
    - **Installation on Windows:**
      - Download `yt-dlp.exe` from the [official GitHub releases page](https://github.com/yt-dlp/yt-dlp/releases).
      - Place the `yt-dlp.exe` file in a directory that is included in your system's `PATH` environment variable (e.g., `C:\Windows`, or a custom directory you've added to PATH).

3.  **`ffmpeg`:**
    - `ffmpeg` is required by `yt-dlp` for post-processing audio, such as converting to WAV format. The application explicitly checks for its presence. It must be installed and available in your system's `PATH`.
    - **Installation on Linux (Debian/Ubuntu):**
      ```bash
      sudo apt update
      sudo apt install ffmpeg
      ```
    - **Installation on macOS (using Homebrew):**
      ```bash
      brew install ffmpeg
      ```
    - **Installation on Windows:**
      - Download `ffmpeg` from [https://ffmpeg.org/download.html](https://ffmpeg.org/download.html).
      - Add its binaries directory (e.g., `C:\ffmpeg\bin`) to your system's `PATH`.

### Building the Tool

1.  **Clone the repository** (if you haven't already):

    ```bash
    git clone https://github.com/your-username/yt-transcribe.git # Replace with actual repo if applicable
    cd yt-transcribe
    ```

    _(Note: Assuming the current directory is the project root)_

2.  **Build the executable:**
    Navigate to the project's root directory in your terminal and run:

    ```bash
    go build -a -v -o yt-transcribe
    ```

    This command will compile your Go code and create an executable file named `yt-transcribe` (or `yt-transcribe.exe` on Windows) in the current directory.

3.  **Run tests (optional):**
    To ensure everything is working correctly, you can run the unit tests with:
    ```bash
    go test -v ./
    ```

## 💡 Usage

To run the `yt-transcribe` tool, you can optionally provide a video URL using the `-url` flag. If no URL is provided, a default video will be used. You can also specify a custom output directory using the `-output` flag.

```bash
./yt-transcribe [-url <VIDEO_URL>] [-output <OUTPUT_DIRECTORY>]
```

- Replace `<VIDEO_URL>` with the actual link to the video you want to download audio from. If omitted, the default URL `https://www.youtube.com/watch?v=rdWZo5PD9Ek` will be used.
- Replace `<OUTPUT_DIRECTORY>` with the path where you want the `.wav` audio file to be saved. If omitted, the audio will be saved in your system's temporary directory.

### Examples:

1.  **Download audio from the default video and save to the default temporary directory:**

    ```bash
    ./yt-transcribe
    ```

    _(This will download audio from `https://www.youtube.com/watch?v=rdWZo5PD9Ek` and save a file like `video_2025-12-31.wav` in your temporary directory)_

2.  **Download audio from a specified video and save to the default temporary directory:**

    ```bash
    ./yt-transcribe -url https://www.youtube.com/watch?v=dQw4w9WgXcQ
    ```

    _(This will save a file like `video_2025-12-31.wav` in your temporary directory)_

3.  **Download audio from a specified video and save to a specific directory:**
    ```bash
    ./yt-transcribe -url https://www.youtube.com/watch?v=your_video_id -output ~/youtube_audio
    ```
    _(This will save a file like `video_2025-12-31.wav` in the `~/youtube_audio` directory)_

## 🤝 Contributing

(If this were an open-source project, you'd find contribution guidelines here.)

## 📄 License

(If this were an open-source project, you'd find license information here.)
