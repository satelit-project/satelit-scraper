#!/usr/bin/env ash

set -euo pipefail

main() {
  echo "Running service" >&2
  ST_LOG=prod \
    ST_IMPORTER_URL="$ST_IMPORTER_URL" \
    DO_SPACES_KEY="$DO_SPACES_KEY" \
    DO_SPACES_SECRET="$DO_SPACES_SECRET" \
    DO_SPACES_HOST="$DO_SPACES_HOST" \
    DO_BUCKET="$DO_BUCKET" \
    exec ./satelit-scraper
}

main "$@"
