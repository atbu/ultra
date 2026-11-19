#!/usr/bin/env bash
# Builds the static web UI into a folder (default: dist) that can be served by any static file host, such as GitHub
# Pages. Usage: web/build.sh [output-directory]
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
out="${1:-$root/dist}"

rm -rf "$out"
mkdir -p "$out"

cp "$root/web/index.html" "$root/web/style.css" "$root/web/app.js" "$out/"

# The JavaScript glue that loads Go WebAssembly must match the Go version used to compile it, so copy it from the
# toolchain rather than keeping a copy in the repository.
goroot="$(go env GOROOT)"
if [ -f "$goroot/lib/wasm/wasm_exec.js" ]; then
  cp "$goroot/lib/wasm/wasm_exec.js" "$out/"
else
  cp "$goroot/misc/wasm/wasm_exec.js" "$out/"
fi

GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$out/ultra.wasm" "$root/cmd/wasm"

echo "Built web UI in $out"
