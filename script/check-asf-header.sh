if command -v docker >/dev/null 2>&1; then
    CONTAINER_RUNTIME="docker"
elif command -v podman >/dev/null 2>&1; then
    CONTAINER_RUNTIME="podman"
else
    echo "Neither Docker nor Podman is installed. Please install either Docker or Podman."
    exit 1
fi

$CONTAINER_RUNTIME run --rm -v "$(pwd)":/github/workspace ghcr.io/korandoru/hawkeye-native format

gofmt -w -l .