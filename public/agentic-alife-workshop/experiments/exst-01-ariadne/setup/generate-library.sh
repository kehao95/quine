#!/bin/bash
# generate-library.sh - Create the Library of Babel for The Borges Project
# Usage: ./generate-library.sh [output_dir] [size]
#   output_dir: where to create the library (default: ./library)
#   size: number of volumes to generate (default: 1000)

set -e

OUTPUT_DIR="${1:-./library}"
NUM_VOLUMES="${2:-1000}"

echo "=== Generating The Library of Babel ==="
echo "Output: $OUTPUT_DIR"
echo "Volumes: $NUM_VOLUMES"
echo ""

# Clean up if exists
if [ -d "$OUTPUT_DIR" ]; then
    echo "Cleaning existing library..."
    rm -rf "$OUTPUT_DIR"
fi

# Create structure: hex_XX/shelf_XX/volume_XXXXX.txt
# For 1000 volumes: 10 hexes × 10 shelves × 10 volumes = 1000

HEX_COUNT=10
SHELF_COUNT=10
VOL_PER_SHELF=10

echo "Structure: $HEX_COUNT hexes × $SHELF_COUNT shelves × $VOL_PER_SHELF volumes"
echo ""

generated=0
for ((h=0; h<HEX_COUNT; h++)); do
    hex_num=$(printf "%02d" $h)
    hex_dir="$OUTPUT_DIR/hex_$hex_num"
    mkdir -p "$hex_dir"

    for ((s=0; s<SHELF_COUNT; s++)); do
        shelf_num=$(printf "%02d" $s)
        shelf_dir="$hex_dir/shelf_$shelf_num"
        mkdir -p "$shelf_dir"

        for ((v=0; v<VOL_PER_SHELF; v++)); do
            vol_num=$(printf "%05d" $((h * 100 + s * 10 + v)))
            vol_path="$shelf_dir/volume_$vol_num.txt"

            # Each file is just a UUID — pure random, but text-safe
            uuidgen > "$vol_path"

            generated=$((generated + 1))

            if [ $((generated % 500)) -eq 0 ]; then
                echo "  Generated $generated volumes..."
            fi
        done
    done
done

echo ""
echo "=== Library Generation Complete ==="
echo "Total volumes: $generated"
echo "Library size: $(du -sh $OUTPUT_DIR | cut -f1)"
