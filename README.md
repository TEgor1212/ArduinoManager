# Arduino Manager GUI (Go + Fyne) — CI/CD и GitHub Releases

Проект **Arduino Manager GUI** — десктопное приложение на Go с интерфейсом Fyne для поиска/установки Arduino-библиотек и board core (через GitHub и локальный индекс), управления установленными компонентами и конфигурацией окружения Arduino.

В репозитории настроен **CI/CD на GitHub Actions**:
- **CI** (`test`): проверка форматирования, `go vet`, тесты и smoke-сборка на `ubuntu-latest` при push/PR.
- **CD** (`release`): сборка GUI-бинарников для Linux/macOS/Windows и публикация их в **GitHub Releases** при пуше тега вида `v*`.

> **Важно:** push и публикация тегов выполняются по указанию пользователя. Фактический запуск Actions/Releases и получение скриншотов происходят в GitHub/графической среде.

## Состав проекта

- `main.go` — основной исходный код GUI-приложения (Fyne) и логика менеджера.
- `logic.go` — вспомогательные функции/утилиты, используемые приложением и тестами.
- `logic_test.go` — unit-тесты вспомогательной логики.
- `go.mod`, `go.sum` — описание модуля и зависимостей Go/Fyne.
- `.github/workflows/ci.yml` — workflow GitHub Actions (CI + релизная публикация бинарников).
- `Dockerfile.test` — вспомогательный образ для тестов/сборки Fyne в Linux-окружении с GUI-зависимостями.
- `docs/screenshots/` — каталог для скриншотов документации (заполняется после запуска Actions/приложения).

## CI/CD пайплайн (`.github/workflows/ci.yml`)

### Триггеры
- `push` в ветку `main`
- `push` тегов `v*`
- `pull_request`

### Job `test` (CI, Ubuntu)
1. Checkout репозитория (`actions/checkout@v4`).
2. Установка Go 1.23 (`actions/setup-go@v5`, cache).
3. Установка системных зависимостей для сборки Fyne/CGO на Linux (`mesa`, `xorg`, `wayland` и др.).
4. Проверка форматирования: `gofmt -l .` (падает при наличии неотформатированных файлов).
5. Линтинг: `go vet ./...`.
6. Тесты: `go test ./... -v`.
7. Smoke-сборка: `go build -o /tmp/arduino-manager .`.

### Job `release` (CD, GitHub Releases)
- Запускается только для тегов `refs/tags/v*` и после успешного `test`.
- Требует `permissions: contents: write` для публикации релиза.
- Матрица сборки:
  - `ubuntu-latest` → `arduinomanager-linux-x64`
  - `macos-14` (arm64) → `arduinomanager-macos-arm64`
  - `windows-latest` → `arduinomanager-windows-x64.exe` (с `-H windowsgui`)
- На Windows настраивается MSYS2/MinGW64 (`gcc`, `pkg-config`) для CGO-сборки Fyne.
- Сборка: `CGO_ENABLED=1 go build -ldflags "-s -w ... -X main.version=<tag>" ...`
- Публикация артефактов в релиз: `softprops/action-gh-release@v2`.

## Поэтапный процесс: Workflow и публикация релиза (задача 2)

Ниже — пошаговая инструкция и места для демонстрации скриншотов. Скриншоты нужно добавить после фактического запуска в GitHub.

### Шаг 1. Подготовка и локальная проверка
```bash
# Форматирование (если установлен Go)
gofmt -w .

# Проверки, которые выполняет CI
go vet ./...
go test ./... -v
go build -o arduinomanager .
```
**Ожидаемый результат:** форматирование корректно, тесты проходят, бинарник собирается.

> В headless-среде без установленного Go локальный прогон команд CI может быть недоступен; в этом случае проверка выполняется в GitHub Actions.

### Шаг 2. Запуск CI (push/PR в `main`)
После push изменений в ветку `main` (или создания PR) запускается workflow CI (`test`).

**Скриншот (успешный Workflow/CI):**

![Успешный workflow CI](docs/screenshots/workflow-success.png)

_Если скриншот ещё не добавлен, сохраните снимок успешного job `test` в GitHub Actions как `docs/screenshots/workflow-success.png`._

### Шаг 3. Создание тега и публикация релиза
```bash
# Пример версии
git tag v1.0.0
git push origin v1.0.0
```
Пуш тега `v*` запускает job `release`, который собирает бинарники под Linux/macOS/Windows и прикрепляет их к GitHub Release.

**Скриншот (Actions после пуша тега):**

![Запуск Actions для релиза](docs/screenshots/actions-run.png)

**Скриншот (GitHub Release с бинарниками):**

![GitHub Release с бинарниками](docs/screenshots/release-assets.png)

### Шаг 4. Проверка результатов
- Вкладка **Actions**: jobs `test` и `release` завершены успешно.
- Вкладка **Releases**: создан релиз по тегу, в Assets доступны:
  - `arduinomanager-linux-x64`
  - `arduinomanager-macos-arm64`
  - `arduinomanager-windows-x64.exe`

## Этапы разработки проекта (задача 3)

Ниже описаны этапы разработки/подготовки проекта и точки для скриншотов (включая окно запущенного приложения).

### 1. Исходное приложение Arduino Manager GUI
- Реализован GUI на **Go + Fyne** (`main.go`) с поиском библиотек/плат через GitHub API и локальный индекс Arduino.
- Добавлены функции установки/обновления/удаления библиотек и board cores, работа с конфигурацией и логом операций.
- Вспомогательная логика поддерживается в `logic.go`, покрыта тестами в `logic_test.go`.

**Скриншот (окно запущенного приложения):**

![Окно запущенного приложения Arduino Manager](docs/screenshots/app-window.png)

_Для получения скриншота запустите приложение локально в графической среде (`go run .` или собранный бинарник), откройте главное окно и сохраните снимок в `docs/screenshots/app-window.png`._

### 2. Подготовка зависимостей и сборки
- Настроен Go-модуль (`go.mod`, `go.sum`) с зависимостью `fyne.io/fyne/v2`.
- Учтены требования CGO/системных GUI-библиотек для сборки Fyne на разных ОС.
- В `main.go` добавлена переменная `version` (`dev` по умолчанию) для инъекции версии релиза через `-ldflags -X main.version=...` в CI.

### 3. Настройка CI (качество и проверка сборки)
- В `.github/workflows/ci.yml` добавлен job `test`:
  - проверка форматирования (`gofmt`),
  - статический анализ (`go vet`),
  - запуск тестов (`go test ./...`),
  - smoke-сборка приложения.
- Для Linux в CI устанавливаются зависимости Fyne (OpenGL/X11/Wayland dev-пакеты).

### 4. Настройка CD и публикации в GitHub Releases
- Добавлен job `release`, зависящий от успешного `test` и запускаемый по тегам `v*`.
- Реализована кросс-платформенная сборка бинарников (Linux, macOS arm64, Windows x64) с CGO.
- Для Windows настроена сборочная среда MSYS2/MinGW64.
- Публикация собранных артефактов выполняется в GitHub Releases через `softprops/action-gh-release`.

### 5. Документирование и эксплуатация
- Подготовлен данный `README.md` с описанием пайплайна, пошаговым процессом запуска workflow/релиза и этапами разработки.
- Предусмотрены места для скриншотов успешных Actions/Releases и окна приложения (`docs/screenshots/`).

## Локальная сборка и запуск (для разработки)

Требования: Go (версия из `go.mod`, например 1.23+), C-компилятор и системные зависимости Fyne/CGO для вашей ОС.

```bash
# Зависимости модуля
go mod download

# Запуск
go run .

# Сборка
go build -o arduinomanager .
```

Linux (пример зависимостей для Fyne/CI):
```bash
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev libegl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev wayland-protocols
```

## Что нужно добавить для полной демонстрации (скриншоты)
- `docs/screenshots/workflow-success.png` — успешный CI workflow (`test`).
- `docs/screenshots/actions-run.png` — запуск Actions (включая релизный при теге).
- `docs/screenshots/release-assets.png` — страница GitHub Release с опубликованными бинарниками.
- `docs/screenshots/app-window.png` — окно запущенного приложения (GUI-скриншот, недоступен в headless-среде).

## Примечания об окружении и ограничениях
- GUI-скриншот окна приложения и фактические скриншоты GitHub Actions/Releases не могут быть получены автоматически в headless-окружении без графической сессии и без выполнения workflow в GitHub.
- Публикация релизных бинарников требует пуша тега `v*` в удалённый репозиторий; результат зависит от настроек GitHub Actions/permissions и доступности зависимостей в раннерах.
