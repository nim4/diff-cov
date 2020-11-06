#!/usr/bin/env sh
set -e

rm diff-cov diff-cov.tar.gz 2>/dev/null || true

GOOS=linux go build ./cmd/diff-cov

chmod +x diff-cov
tar -czvf diff-cov.tar.gz diff-cov