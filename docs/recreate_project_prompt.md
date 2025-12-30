# Prompt to Recreate YouTube Transcriber CLI Project

This prompt outlines the steps, tools, and principles required to create the `yt-transcribe` CLI tool from scratch. It is designed to be comprehensive enough for an AI agent or a developer to understand the project's requirements and structure.

---

## Project Goal

Create a Command Line Interface (CLI) tool in Go that transcribes YouTube videos into text.

## Core Features

1.  **Accept YouTube Link:** The CLI tool must accept a valid YouTube video URL as a command-line argument.
2.  **Audio Download:** Download the audio stream from the provided YouTube video.
3.  **Audio Transcription:** Transcribe the downloaded audio into text.
4.  **Output Text:** Save the transcribed text to a file (defaulting to the system's temporary directory, with an option for a user-specified directory).
5.  **Temporary File Cleanup:** Automatically remove temporary audio files after transcription.

## Technology Stack & Tools

*   **Programming Language:** Go (Golang)
*   **CLI Argument Parsing:** Go's standard `flag` package.
*   **Audio Downloading:** Integration with the external `yt-dlp` command-line tool via `os/exec`.
*   **Audio Transcription:** Conceptual (mocked) integration with an external transcription service (e.g., OpenAI Whisper API).

## Principles & Architecture

The project should be developed following these principles:

1.  **SOLID Principles:**
    *   **Interface-Oriented Design:** Emphasize interfaces (`YouTubeDownloader`, `Transcriber`) for defining contracts and promoting extensibility.
    *   **Dependency Inversion:** `main.go` should depend on abstractions (interfaces) rather than concrete implementations.
    *   **Single Responsibility Principle (SRP):** Each package, interface, and struct should have a single, well-defined responsibility.
2.  **Modular Structure:**
    *   Organize code into logical packages (e.g., `pkg/downloader`, `pkg/transcriber`).
    *   **Explicit File Separation:** Separate interface definitions from their concrete implementations into distinct files within their respective packages for clarity and adherence to Go idioms (e.g., `pkg/downloader/interface.go` for the interface, `pkg/downloader/yt_dlp.go` for the `yt-dlp` implementation).
3.  **Robust Error Handling:** Implement clear and informative error handling throughout the application.
4.  **Comprehensive Comments:** Add extensive comments explaining the purpose of packages, interfaces, structs, functions, and complex logic sections.

## Project Setup & Conventions

1.  **Go Module:** Initialize a standard Go module (e.g., `yt-transcribe`).
2.  **`.gitignore`:** Create a `.gitignore` file specific to Go projects, including common Go build artifacts, IDE files, and the generated executable.
3.  **Documentation & Project Info:**
    *   `README.md`: Provide a comprehensive `README.md` for tool usage, prerequisites, building, and running instructions.
    *   `LICENSE`: Include an appropriate open-source license (e.g., MIT).
    *   `CONTRIBUTING.md`: Provide guidelines for potential contributors.
    *   `SECURITY.md`: Outline the process for reporting security vulnerabilities.

## Testing

1.  **Unit Tests:** Develop thorough unit tests for all core components.
2.  **Mocking External Dependencies:**
    *   For the `downloader` package, implement robust mocking for `os/exec.Command`, `os.LookPath`, and `os.Stat` to simulate various scenarios (e.g., `yt-dlp` not found, successful download, command failures, missing output). Use global variables for these functions to allow them to be overridden in tests.
    *   For the `transcriber` package, test the mock implementation's expected behavior.
3.  **Helper Function Tests:** Test any helper functions (e.g., `sanitizeFilename`).
4.  **Test Isolation:** Ensure tests are isolated and do not interfere with each other (e.g., using `t.TempDir()` and mutexes for global variable mocks).

---

## Development Steps (Chronological Hint for Implementation)

1.  Initialize Go module and create `main.go`.
2.  Define CLI argument parsing in `main.go`.
3.  Create `pkg/downloader` package:
    *   Define `YouTubeDownloader` interface in `pkg/downloader/interface.go`.
    *   Implement `YTDLPAudioDownloader` in `pkg/downloader/yt_dlp.go`, wrapping `yt-dlp`.
    *   Introduce global mock variables (`commandExecutor`, `osLookPath`, `osStat`, `cmdCombinedOutput`) in `pkg/downloader/yt_dlp.go` for testability.
4.  Create `pkg/transcriber` package:
    *   Define `Transcriber` interface in `pkg/transcriber/interface.go`.
    *   Implement `OpenAITranscriber` (mock) in `pkg/transcriber/openai.go`.
5.  Integrate `downloader` and `transcriber` in `main.go` using dependency injection.
6.  Implement `sanitizeFilename` helper in `main.go`.
7.  Add `go.mod` for dependencies.
8.  Create `.gitignore`, `README.md`, `LICENSE`, `CONTRIBUTING.md`, `SECURITY.md`.
9.  Implement unit tests:
    *   `main_test.go` for `sanitizeFilename`.
    *   `pkg/transcriber/openai_test.go` for `OpenAITranscriber`.
    *   `pkg/downloader/yt_dlp_test.go` for `YTDLPAudioDownloader`, with extensive mocking.

---
**Note on "Auto-Update":** This prompt reflects the project's state as of its generation. If the project undergoes further development, this prompt would need to be regenerated by analyzing the new codebase to capture the changes.
