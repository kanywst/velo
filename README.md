# Velo

Velo is a network speed measurement tool.
It allows you to monitor your internet connection speed (Download, Upload, Latency) automatically over time and visualize the results in a graph.

## Features

- **Speed Test**: Measures download speed, upload speed, and latency using `speedtest-go`.
- **Automatic Monitoring**: Automatically runs a speed test every hour while the application is running.
- **Visualization**: Displays historical data on an interactive chart (Time vs Speed).
- **Cross-Platform**: Built on Wails v2, targeting desktop environments (macOS, Windows, Linux).

## Prerequisites

- **Go** (v1.25 or later)
- **Node.js** & **npm**
- **Wails CLI**

  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

## Project Setup

1. **Clone/Navigate to the directory:**

   ```bash
   cd velo
   ```

2. **Install Backend Dependencies:**

   ```bash
   go mod tidy
   ```

3. **Install Frontend Dependencies:**

   ```bash
   cd frontend
   npm install
   cd ..
   ```

## Development

To start the application in development mode with hot-reloading:

```bash
wails dev
```

This command will:

1. Build the frontend assets.
2. Compile the Go backend.
3. Launch the application window.
4. Watch for file changes in both frontend and backend.

## Building for Production

To create a production-ready binary:

```bash
wails build
```

The output binary will be located in the `build/bin` directory.
