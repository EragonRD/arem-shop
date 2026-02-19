#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="$ROOT_DIR/mermaid"
OUT_DIR="$ROOT_DIR/schemas"

mkdir -p "$OUT_DIR"

DIAGRAMS=(
  "plan-explication"
  "architecture-globale"
  "pipeline-securite"
  "erd"
  "sequence-login"
  "sequence-transaction-sale"
  "sequence-public-catalog"
)

render_diagram() {
  local name="$1"
  local input_file="$SRC_DIR/${name}.md"
  local output_file="$OUT_DIR/${name}.png"
  local output_prefix="$OUT_DIR/${name}.png"
  local generated_file="$OUT_DIR/${name}-1.png"

  if [[ ! -f "$input_file" ]]; then
    echo "[error] missing input file: $input_file" >&2
    exit 1
  fi

  npx -y @mermaid-js/mermaid-cli \
    -i "$input_file" \
    -o "$output_prefix" \
    -e png \
    -w 1920 \
    -H 1080 \
    -s 2 \
    -b white \
    -q

  if [[ -f "$generated_file" ]]; then
    mv "$generated_file" "$output_file"
  elif [[ ! -f "$output_file" ]]; then
    echo "[error] diagram output not found for $name" >&2
    exit 1
  fi

  echo "[ok] $output_file"
}

for diagram in "${DIAGRAMS[@]}"; do
  render_diagram "$diagram"
done

echo "All diagrams generated in $OUT_DIR"
