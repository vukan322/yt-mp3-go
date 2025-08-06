# YouTube Downloader

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Live App](https://img.shields.io/badge/Live_App-Online-brightgreen)](https://lukavukanovic.xyz/yt-downloader/)
[![Build and Deploy Go App](https://github.com/vukan322/yt-mp3-go/actions/workflows/deploy.yml/badge.svg)](https://github.com/vukan322/yt-mp3-go/actions/workflows/deploy.yml)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker&logoColor=white)](https://github.com/vukan322/yt-mp3-go#getting-started)


A modern, feature-rich web application for downloading YouTube videos as MP3 files.
Built with a robust Go backend and a dynamic, responsive frontend.

|                                     Initial Page                                      |                                         Download Options                                          |
|:-------------------------------------------------------------------------------------:|:-------------------------------------------------------------------------------------------------:|
| <img src="/screenshots/app-screenshot.png" alt="Initial Page Screenshot" width="450"> | <img src="/screenshots/app-screenshot-options.png" alt="Download Options Screenshot" width="450"> |

## Key Features

- **Asynchronous Downloads:** Utilizes Go's concurrency with goroutines to process downloads in the background, allowing the server to handle multiple users simultaneously without blocking.
- **Live Status Updates:** Employs Server-Sent Events (SSE) to push real-time progress updates (e.g., "Processing", "Complete", "Failed") to the user's browser.
- **Dynamic Frontend:** The UI is built with vanilla JavaScript, using AJAX for form submissions and dynamically updating the page content without requiring a full reload.
- **Internationalization (i18n):** Fully localized interface supporting both English and Serbian, with language detection and a manual switcher.
- **Modern UI/UX:** Features a clean, responsive design with light/dark themes, video thumbnail previews, file size display, and state persistence across page reloads using `sessionStorage`.
- **Audio Quality Control:** Allows users to select from multiple quality presets (Low, Medium, High) to balance file size and audio quality.
- **Custom Filenames:** Lets users easily edit the MP3 filename in the UI before downloading for better file organization.
- **Audio Normalization:** Provides an option to normalize audio to a standard volume level, ensuring a consistent listening experience across different tracks.
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

- **Go** (version 1.22 or later)
- **`yt-dlp`** & **`ffmpeg`**:

  ### macOS
  ```sh
  brew install yt-dlp/taps/yt-dlp ffmpeg
  ```

  ### Linux
  ```sh
  sudo apt update && sudo apt install yt-dlp ffmpeg
  ```
  
<hr>

## Getting Started

### Docker (Recommended)

The easiest way to run this application is using Docker. No need to install Go, yt-dlp, or ffmpeg.

```sh
  git clone https://github.com/vukan322/yt-mp3-go.git
  cd yt-mp3-go
  docker-compose up --build
```

Or manually with Docker:
```sh
  docker build -t yt-mp3-downloader .
  docker run -p 8080:8080 -v ./downloads:/app/downloads -v ./cookies.txt:/app/cookies.txt:ro yt-mp3-downloader
```
<hr>

### Local Development

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/vukan322/yt-mp3-go.git
    cd yt-mp3-go
    ```

2.  **Set up your environment:**
    Copy the example environment file and edit it if needed.
    ```sh
    cp .env.example .env
    ```

3.  **Tidy dependencies:**
    This will download the required Go modules.
    ```sh
    go mod tidy
    ```

4.  **Make build script executable:**
    You only need to do this once.
    ```sh
    chmod +x build-assets.sh
    ```

5.  **Run the application:**
    ```sh
    make run
    ```
    OR
    ```sh
    ./build-assets.sh && go run cmd/server/main.go
    ```

### Access the Application

Open your browser and navigate to `http://localhost:8080/yt-downloader`.

## Configuration

The application is configured using environment variables, which are loaded from an .env file during local development.

| Variable                      | Description                                                           | Default                       |
|:------------------------------|:----------------------------------------------------------------------|:------------------------------|
| **APP_ENV**                   | The application environment (**development** or **production**).      | **development**               |
| **DOMAIN**                    | The domain the application is running on.                             | **localhost**                 |
| **PORT**                      | The port the application will listen on.                              | **8080**                      |
| **BASE_PATH**                 | The URL path prefix for the application.                              | **/yt-downloader**            |

## Deployment

This project is configured for continuous deployment using GitHub Actions.

 - **Trigger**: A push to the `main` branch automatically triggers the build and deployment workflow.
 - **Strategy**: The deployment is **atomic**, meaning it uses a timestamped release directory and a current symlink to ensure zero downtime.
 - **Versioning**: The application version is automatically injected at build time. The format is `MAJOR.MINOR.PATCH`, 
 where the `PATCH` number is the GitHub Actions run number.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
