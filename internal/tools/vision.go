package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/draw"

	"github.com/kehao95/quine/internal/tape"
)

// maxImageDimension is the maximum allowed width or height for images.
// Keeps images reasonably sized to reduce token cost and API latency.
const maxImageDimension = 1096

type visionStructuredResult struct {
	Tool        string `json:"tool"`
	Status      string `json:"status"`
	Path        string `json:"path,omitempty"`
	Instruction string `json:"instruction,omitempty"`
	MIMEType    string `json:"mime_type,omitempty"`
	Scaled      bool   `json:"scaled,omitempty"`
	Error       string `json:"error,omitempty"`
}

// HandleVision reads the image at args["path"], base64-encodes it, detects its
// MIME type, and returns a ToolResult with the Image field populated.
// Images larger than maxImageDimension are automatically scaled down.
// It does NOT consume a turn.
func HandleVision(toolID string, args map[string]any) tape.ToolResult {
	return HandleVisionWithReader(toolID, args, os.ReadFile)
}

// HandleVisionWithReader is HandleVision with an injectable file reader. The
// runtime uses this to read paths through the active workspace backend instead
// of only through the runtime process' host-visible filesystem.
func HandleVisionWithReader(toolID string, args map[string]any, readFile func(string) ([]byte, error)) tape.ToolResult {
	path, _ := args["path"].(string)
	if path == "" {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(visionStructuredResult{
				Tool:   "vision",
				Status: "error",
				Error:  "[VISION ERROR] missing required parameter: path",
			}),
			IsError: true,
		}
	}

	instruction, _ := args["instruction"].(string)
	if readFile == nil {
		readFile = os.ReadFile
	}

	data, err := readFile(path)
	if err != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(visionStructuredResult{
				Tool:   "vision",
				Status: "error",
				Path:   path,
				Error:  fmt.Sprintf("[VISION ERROR] cannot read file %q: %v", path, err),
			}),
			IsError: true,
		}
	}

	mimeType := detectMIMEType(data, path)
	if mimeType == "" {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(visionStructuredResult{
				Tool:   "vision",
				Status: "error",
				Path:   path,
				Error:  fmt.Sprintf("[VISION ERROR] unsupported image format for file %q (detected by magic bytes and extension)", path),
			}),
			IsError: true,
		}
	}

	// Check if image needs scaling
	scaledData, scaledMIME, wasScaled, scaleErr := maybeScaleImage(data, mimeType)
	if scaleErr != nil {
		return tape.ToolResult{
			ToolID: toolID,
			Content: tape.MarshalToolResultContent(visionStructuredResult{
				Tool:   "vision",
				Status: "error",
				Path:   path,
				Error:  fmt.Sprintf("[VISION ERROR] failed to process image %q: %v", path, scaleErr),
			}),
			IsError: true,
		}
	}

	encoded := base64.StdEncoding.EncodeToString(scaledData)

	return tape.ToolResult{
		ToolID: toolID,
		Content: tape.MarshalToolResultContent(visionStructuredResult{
			Tool:        "vision",
			Status:      "completed",
			Path:        path,
			Instruction: instruction,
			MIMEType:    scaledMIME,
			Scaled:      wasScaled,
		}),
		Image: &tape.ImagePart{
			MIMEType: scaledMIME,
			Data:     encoded,
		},
	}
}

// maybeScaleImage checks if the image exceeds maxImageDimension and scales it down if needed.
// Returns the (possibly modified) image data, MIME type, whether scaling occurred, and any error.
func maybeScaleImage(data []byte, mimeType string) ([]byte, string, bool, error) {
	if isNetpbmMIME(mimeType) {
		img, err := decodeNetpbm(data)
		if err != nil {
			return data, mimeType, false, err
		}
		return encodePNGMaybeScaled(img)
	}

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

func encodePNGMaybeScaled(img image.Image) ([]byte, string, bool, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	wasScaled := false

	if width > maxImageDimension || height > maxImageDimension {
		newWidth, newHeight := width, height
		if width > height {
			newWidth = maxImageDimension
			newHeight = height * maxImageDimension / width
		} else {
			newHeight = maxImageDimension
			newWidth = width * maxImageDimension / height
		}
		dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
		img = dst
		wasScaled = true
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, "", false, fmt.Errorf("encoding converted image: %w", err)
	}
	return buf.Bytes(), "image/png", wasScaled, nil
}

func isNetpbmMIME(mimeType string) bool {
	switch mimeType {
	case "image/x-portable-pixmap", "image/x-portable-graymap", "image/x-portable-anymap":
		return true
	default:
		return false
	}
}

func decodeNetpbm(data []byte) (image.Image, error) {
	pos := 0
	magic, err := nextNetpbmToken(data, &pos)
	if err != nil {
		return nil, err
	}
	if magic != "P2" && magic != "P3" && magic != "P5" && magic != "P6" {
		return nil, fmt.Errorf("unsupported Netpbm magic %q", magic)
	}

	width, err := nextNetpbmInt(data, &pos, "width")
	if err != nil {
		return nil, err
	}
	height, err := nextNetpbmInt(data, &pos, "height")
	if err != nil {
		return nil, err
	}
	maxValue, err := nextNetpbmInt(data, &pos, "max value")
	if err != nil {
		return nil, err
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid Netpbm dimensions %dx%d", width, height)
	}
	if maxValue <= 0 || maxValue > 65535 {
		return nil, fmt.Errorf("invalid Netpbm max value %d", maxValue)
	}

	switch magic {
	case "P2", "P3":
		return decodeASCIIGrid(data, pos, magic == "P3", width, height, maxValue)
	case "P5", "P6":
		return decodeBinaryGrid(data, pos, magic == "P6", width, height, maxValue)
	default:
		panic("unreachable Netpbm magic")
	}
}

func nextNetpbmInt(data []byte, pos *int, name string) (int, error) {
	token, err := nextNetpbmToken(data, pos)
	if err != nil {
		return 0, fmt.Errorf("read Netpbm %s: %w", name, err)
	}
	value, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("invalid Netpbm %s %q", name, token)
	}
	return value, nil
}

func nextNetpbmToken(data []byte, pos *int) (string, error) {
	if err := skipNetpbmWhitespaceAndComments(data, pos); err != nil {
		return "", err
	}
	if *pos >= len(data) {
		return "", fmt.Errorf("unexpected end of Netpbm header")
	}
	start := *pos
	for *pos < len(data) && !isNetpbmWhitespace(data[*pos]) {
		(*pos)++
	}
	return string(data[start:*pos]), nil
}

func skipNetpbmWhitespaceAndComments(data []byte, pos *int) error {
	for *pos < len(data) {
		if isNetpbmWhitespace(data[*pos]) {
			(*pos)++
			continue
		}
		if data[*pos] == '#' {
			for *pos < len(data) && data[*pos] != '\n' {
				(*pos)++
			}
			continue
		}
		return nil
	}
	return nil
}

func isNetpbmWhitespace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	default:
		return false
	}
}

func decodeBinaryGrid(data []byte, pos int, colorImage bool, width, height, maxValue int) (image.Image, error) {
	if pos >= len(data) || !isNetpbmWhitespace(data[pos]) {
		return nil, fmt.Errorf("Netpbm binary raster missing header separator")
	}
	pos++
	if pos < len(data) && data[pos-1] == '\r' && data[pos] == '\n' {
		pos++
	}

	channels := 1
	if colorImage {
		channels = 3
	}
	bytesPerSample := 1
	if maxValue > 255 {
		bytesPerSample = 2
	}
	expected := int64(width) * int64(height) * int64(channels) * int64(bytesPerSample)
	if expected < 0 || int64(len(data)-pos) < expected {
		return nil, fmt.Errorf("Netpbm raster is truncated")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	offset := pos
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if colorImage {
				r := readNetpbmSample(data, &offset, bytesPerSample, maxValue)
				g := readNetpbmSample(data, &offset, bytesPerSample, maxValue)
				b := readNetpbmSample(data, &offset, bytesPerSample, maxValue)
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
				continue
			}
			v := readNetpbmSample(data, &offset, bytesPerSample, maxValue)
			img.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	return img, nil
}

func decodeASCIIGrid(data []byte, pos int, colorImage bool, width, height, maxValue int) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if colorImage {
				r, err := nextNetpbmInt(data, &pos, "red sample")
				if err != nil {
					return nil, err
				}
				g, err := nextNetpbmInt(data, &pos, "green sample")
				if err != nil {
					return nil, err
				}
				b, err := nextNetpbmInt(data, &pos, "blue sample")
				if err != nil {
					return nil, err
				}
				img.SetRGBA(x, y, color.RGBA{
					R: scaleNetpbmSample(r, maxValue),
					G: scaleNetpbmSample(g, maxValue),
					B: scaleNetpbmSample(b, maxValue),
					A: 255,
				})
				continue
			}
			v, err := nextNetpbmInt(data, &pos, "gray sample")
			if err != nil {
				return nil, err
			}
			gray := scaleNetpbmSample(v, maxValue)
			img.SetRGBA(x, y, color.RGBA{R: gray, G: gray, B: gray, A: 255})
		}
	}
	return img, nil
}

func readNetpbmSample(data []byte, offset *int, bytesPerSample, maxValue int) uint8 {
	if bytesPerSample == 1 {
		value := int(data[*offset])
		(*offset)++
		return scaleNetpbmSample(value, maxValue)
	}
	value := int(data[*offset])<<8 | int(data[*offset+1])
	*offset += 2
	return scaleNetpbmSample(value, maxValue)
}

func scaleNetpbmSample(value, maxValue int) uint8 {
	if value < 0 {
		value = 0
	}
	if value > maxValue {
		value = maxValue
	}
	return uint8((value*255 + maxValue/2) / maxValue)
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
	if len(data) >= 2 && data[0] == 'P' {
		switch data[1] {
		case '2', '5':
			return "image/x-portable-graymap"
		case '3', '6':
			return "image/x-portable-pixmap"
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
	case ".pgm":
		return "image/x-portable-graymap"
	case ".ppm":
		return "image/x-portable-pixmap"
	case ".pnm":
		return "image/x-portable-anymap"
	}

	return ""
}
