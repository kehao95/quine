package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"

	"github.com/kehao95/quine/internal/tape"
)

// maxImageDimension is the maximum allowed width or height for images.
// Keeps images reasonably sized to reduce token cost and API latency.
const maxImageDimension = 1096

// HandleVision reads the image at args["path"], base64-encodes it, detects its
// MIME type, and returns a ToolResult with the Image field populated.
// Images larger than maxImageDimension are automatically scaled down.
// It does NOT consume a turn.
func HandleVision(toolID string, args map[string]any) tape.ToolResult {
	path, _ := args["path"].(string)
	if path == "" {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[VISION ERROR] missing required parameter: path",
			IsError: true,
		}
	}

	instruction, _ := args["instruction"].(string)

	data, err := os.ReadFile(path)
	if err != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[VISION ERROR] cannot read file %q: %v", path, err),
			IsError: true,
		}
	}

	mimeType := detectMIMEType(data, path)
	if mimeType == "" {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[VISION ERROR] unsupported image format for file %q (detected by magic bytes and extension)", path),
			IsError: true,
		}
	}

	// Check if image needs scaling
	scaledData, scaledMIME, wasScaled, scaleErr := maybeScaleImage(data, mimeType)
	if scaleErr != nil {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: fmt.Sprintf("[VISION ERROR] failed to process image %q: %v", path, scaleErr),
			IsError: true,
		}
	}

	encoded := base64.StdEncoding.EncodeToString(scaledData)

	content := fmt.Sprintf("Image at %s:", path)
	if wasScaled {
		content = fmt.Sprintf("Image at %s (scaled down to fit %dpx limit):", path, maxImageDimension)
	}
	if instruction != "" {
		content = fmt.Sprintf("%s\nYour task: %s", content, instruction)
	}

	return tape.ToolResult{
		ToolID:  toolID,
		Content: content,
		Image: &tape.ImagePart{
			MIMEType: scaledMIME,
			Data:     encoded,
		},
	}
}

// maybeScaleImage checks if the image exceeds maxImageDimension and scales it down if needed.
// Returns the (possibly modified) image data, MIME type, whether scaling occurred, and any error.
func maybeScaleImage(data []byte, mimeType string) ([]byte, string, bool, error) {
	// Only process PNG and JPEG (the formats we can decode/encode)
	if mimeType != "image/png" && mimeType != "image/jpeg" {
		return data, mimeType, false, nil
	}

	// Decode the image to check dimensions
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		// If we can't decode, return original data and let the API handle it
		return data, mimeType, false, nil
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if scaling is needed
	if width <= maxImageDimension && height <= maxImageDimension {
		return data, mimeType, false, nil
	}

	// Calculate new dimensions maintaining aspect ratio
	newWidth, newHeight := width, height
	if width > height {
		if width > maxImageDimension {
			newWidth = maxImageDimension
			newHeight = height * maxImageDimension / width
		}
	} else {
		if height > maxImageDimension {
			newHeight = maxImageDimension
			newWidth = width * maxImageDimension / height
		}
	}

	// Create scaled image using high-quality resampling
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	// Encode the scaled image
	var buf bytes.Buffer
	switch mimeType {
	case "image/png":
		err = png.Encode(&buf, dst)
	case "image/jpeg":
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return data, mimeType, false, fmt.Errorf("encoding scaled image: %w", err)
	}

	return buf.Bytes(), mimeType, true, nil
}

// detectMIMEType identifies the image MIME type from magic bytes, with
// extension fallback. Returns "" for unsupported types.
func detectMIMEType(data []byte, path string) string {
	// Magic byte detection (preferred — extension can lie)
	if len(data) >= 8 {
		switch {
		case data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G':
			return "image/png"
		case data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF:
			return "image/jpeg"
		case data[0] == 'G' && data[1] == 'I' && data[2] == 'F' && data[3] == '8':
			return "image/gif"
		case len(data) >= 12 &&
			data[0] == 'R' && data[1] == 'I' && data[2] == 'F' && data[3] == 'F' &&
			data[8] == 'W' && data[9] == 'E' && data[10] == 'B' && data[11] == 'P':
			return "image/webp"
		}
	}

	// Extension fallback
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	}

	return ""
}
