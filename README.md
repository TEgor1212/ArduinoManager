# Arduino Manager

![CI/CD](https://github.com/TEgor1212/ArduinoManager/actions/workflows/ci.yml/badge.svg)

A **GUI Arduino library and board manager** built with **Go** and **Fyne**.

Arduino Manager searches for Arduino libraries via the **GitHub API**, browses the official Arduino library index, installs/updates/removes libraries and board cores, and manages the Arduino development environment — all from a desktop GUI.

---

## Run

```bash
go run .
```

## Build

```bash
go build -o arduinomanager .
```

> **Note:** Fyne requires CGO and system GUI libraries (OpenGL, X11, etc.). See the [Fyne prerequisites](https://docs.fyne.io/started/) for your OS.

### Linux dependencies (for Fyne)

```bash
sudo apt-get install -y libgl1-mesa-dev libegl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev wayland-protocols
```

---

## CI/CD Pipeline

The workflow in `.github/workflows/ci.yml` runs on pushes to `main`, PRs, and tag pushes (`v*`).

### CI job (push / PR)

Runs on `ubuntu-latest`:

1. **GUI dependencies** — installs Linux packages required by Fyne
2. **Format check** — `gofmt -l .` (fails on unformatted files)
3. **Lint** — `go vet ./...`
4. **Tests** — `go test ./... -v`
5. **Smoke build** — `go build -o /tmp/arduino-manager .`

### Release job (tag push `v*`)

After CI passes, the **release** job builds and publishes binaries to **GitHub Releases**:

| Platform | Binary |
|---|---|
| Linux (x86-64) | `arduinomanager-linux-x64` |
| macOS (arm64) | `arduinomanager-macos-arm64` |
| Windows (x86-64) | `arduinomanager-windows-x64.exe` |

Example:
```bash
git tag v1.0.0
git push origin v1.0.0
```

This creates a GitHub Release with the compiled binaries attached.