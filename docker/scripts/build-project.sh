#!/usr/bin/env ash

set -euo pipefail

ROOT_DIR="/satelit-scraper"
TARGET_DIR="docker/satelit-scraper"
ARTIFACTS_DIR="$ROOT_DIR/$TARGET_DIR"

make_install() {
  echo "Building project" >&2

  CGO_ENABLED=0 \
    go build -a -o satelit-scraper
  mv satelit-scraper "$ARTIFACTS_DIR"
}

copy_resources() {
  echo "Copying resources"

  mkdir -p "$ARTIFACTS_DIR/config"
  cp config/*.yml "$ARTIFACTS_DIR/config"
  cp docker/scripts/entry.sh "$ARTIFACTS_DIR"
}

archive() {
  echo "Packing artifacts"

  apk add tar
  find "$TARGET_DIR/" -type f -o -type l -o -type d \
    | sed s,^"$TARGET_DIR/",, \
    | tar -czf satelit-scraper.tar.gz \
    --no-recursion -C "$TARGET_DIR/" -T -
}

main() {
  mkdir "$ARTIFACTS_DIR"

  make_install
  copy_resources
  archive
}

main "$@"
