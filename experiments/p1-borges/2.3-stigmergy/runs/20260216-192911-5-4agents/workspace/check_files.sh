#!/bin/sh
# Check files for non-random (non-UUID) content
# UUID pattern: 8-4-4-4-12 hex digits with dashes

check_file() {
    file="$1"
    content=$(cat "$file" 2>/dev/null | tr -d '\n\r ')
    
    # Check if it matches UUID pattern exactly
    # Valid UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    if echo "$content" | grep -Eqi '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'; then
        # It's a valid UUID (random)
        return 1
    else
        # Not a UUID - this is interesting!
        echo "FOUND: $file"
        echo "Content: $content"
        return 0
    fi
}

# Export function for xargs -P
export -f check_file

# Read file list from my sample
while read -r file; do
    check_file "$file"
done < coordination/my_sample_proc_15004_1771288182867942000.txt
