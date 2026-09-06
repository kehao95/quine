# Vision Tool Behavior Test

You are testing the `vision` tool. Your job is to verify that the tool correctly reads image files and delivers them as visual input.

## Setup

First, create a small test image using Python:

```
sh(command="python3 -c \"
import struct, zlib

def make_png(width, height, r, g, b):
    def chunk(name, data):
        c = struct.pack('>I', len(data)) + name + data
        return c + struct.pack('>I', zlib.crc32(name + data) & 0xffffffff)
    
    ihdr = struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0)
    raw = b''.join(b'\\x00' + bytes([r,g,b] * width) for _ in range(height))
    idat = zlib.compress(raw)
    
    return b'\\x89PNG\\r\\n\\x1a\\n' + chunk(b'IHDR', ihdr) + chunk(b'IDAT', idat) + chunk(b'IEND', b'')

# Red 4x4 image
with open('/tmp/red.png', 'wb') as f:
    f.write(make_png(4, 4, 255, 0, 0))

# Blue 4x4 image
with open('/tmp/blue.png', 'wb') as f:
    f.write(make_png(4, 4, 0, 0, 255))

print('Images created')
\"")
```

## Part 1: Basic Vision (color identification)

Use `vision(path="/tmp/red.png")` to analyze the red image.
- What color is the image?
- Write `VISION_OK` to fd 4 if you successfully identified it as red/reddish.

## Part 2: Color Discrimination

Use `vision(path="/tmp/blue.png")` to analyze the blue image.
- Can you distinguish it from the red image?
- Write `DISCRIMINATE_OK` to fd 4 if you can tell the two images are different colors.

## Part 3: Error Handling

Call `vision(path="/tmp/nonexistent.png")`.
- The tool should return an error (not crash).
- Write `ERROR_OK` to fd 4 if you received a proper error message.

## Part 4: Turn Counting

After completing Parts 1-3, report how many turns (sh calls) you used.
- Write `TURNS=N` to fd 4 where N is the number of sh calls.
- Note: `vision` calls do NOT count as turns.

## Completion

Write all markers to fd 4 and call exit(status="success").
