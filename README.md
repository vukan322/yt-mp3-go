# YouTube Downloader

A modern, feature-rich web application for downloading YouTube videos as high-quality MP3 files.
Built with a robust Go backend and a dynamic, responsive frontend.

[![YouTube Downloader Screenshot](/screenshots/app-screenshot.png)](https://lukavukanovic.xyz/yt-downloader)
## Key Features

- **Asynchronous Downloads:** Utilizes Go's concurrency with goroutines to process downloads in the background, allowing the server to handle multiple users simultaneously without blocking.
- **Live Status Updates:** Employs Server-Sent Events (SSE) to push real-time progress updates (e.g., "Processing", "Complete", "Failed") to the user's browser.
- **Dynamic Frontend:** The UI is built with vanilla JavaScript, using AJAX for form submissions and dynamically updating the page content without requiring a full reload.
- **Internationalization (i18n):** Fully localized interface supporting both English and Serbian, with language detection and a manual switcher.
- **Modern UI/UX:** Features a clean, responsive design with light/dark themes, video thumbnail previews, file size display, and state persistence across page reloads using `sessionStorage`.
- **Playlist Safe:** Intelligently handles URLs that are part of a playlist, downloading only the single video specified by the user.
- **No Ads, No Fees:** A clean, user-focused experience with no advertisements or hidden costs.

## Built With

This project was brought to life with these modern technologies:

| Technology                    | Description                                                                                           |
|:------------------------------|:------------------------------------------------------------------------------------------------------|
| **Go (Golang)**               | The core language for the backend server, chosen for its performance and concurrency features.        |
| **Vanilla JavaScript (ES6+)** | Used for all client-side logic, including AJAX, DOM manipulation, and handling Server-Sent Events.    |
| **`yt-dlp`**                  | A powerful command-line tool used by the backend to handle the downloading from YouTube.              |
| **`ffmpeg`**                  | An essential external program used by `yt-dlp` to perform the audio conversion to MP3.                |
| **`go-i18n/v2`**              | The Go library used to manage the application's internationalization and translations.                |

---

## Prerequisites

To run this project locally, you must have the following installed and available in your system's `PATH`.

- **Go** (version 1.21 or later)
- **`yt-dlp`** & **`ffmpeg`**:

  ### macOS
  ```sh
  brew install yt-dlp/taps/yt-dlp ffmpeg
  ```

  ### Linux
  ```sh
  sudo apt update && sudo apt install yt-dlp ffmpeg
  ```

## Getting Started

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/vukan322/yt-mp3-go.git
    cd yt-mp3-go
    ```

2.  **Tidy dependencies:**
    This will download the required Go modules.
    ```sh
    go mod tidy
    ```

3.  **Run the application:**
    ```sh
    go run ./cmd/server/main.go
    ```

4.  **Open in your browser:**
    The application will be available at `http://localhost:8080/yt-downloader`.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
