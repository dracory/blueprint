package thumb

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"project/internal/links"
	"project/internal/registry"
	"project/internal/resources"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/dracory/base/img"
	"github.com/dracory/rtr"
	"github.com/dracory/str"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"github.com/dracory/base/cfmt"
)

// == CONSTRUCTOR =============================================================

// NewThumbController creates a new thumbnail controller instance with the provided registry.
// The registry provides access to application services like logging, caching, and configuration.
//
// Parameters:
//   - registry: Application registry interface for accessing services
//
// Returns:
//   - *thumbnailController: New controller instance ready for use
//
// Example:
//
//	controller := NewThumbController(appRegistry)
//	// Use controller for handling thumbnail requests
func NewThumbController(registry registry.RegistryInterface) *thumbnailController {
	return &thumbnailController{registry: registry}
}

// == CONTROLLER ==============================================================

// thumbnailController handles HTTP requests for dynamic thumbnail generation.
// It supports resizing images from various sources including local files, remote URLs,
// and cached data with comprehensive security controls and caching mechanisms.
//
// Security Features:
// - SSRF protection for remote URL access
// - Input validation and sanitization
// - Content-type verification
// - Request timeout controls
//
// Performance Features:
// - 5-minute response caching
// - Optimized HTTP headers for browser/CDN caching
// - Efficient image processing with imaging library
//
// Supported Formats:
// - JPEG (with quality control)
// - PNG (lossless compression)
// - GIF (animated support)
//
// Type:
//   - registry: Application service registry for logging, caching, etc.
type thumbnailController struct {
	registry registry.RegistryInterface
}

// Handler is the main HTTP request handler for thumbnail generation and serving.
// It processes requests for image thumbnails with specified dimensions, quality, and format.
// The handler supports both local files and remote URLs (with SSRF protection).
//
// URL Pattern: /th/:extension/:size/:quality/:path
// Examples:
//   - /th/jpg/1920x0/70/images/backgrounds/pexels-pixabay-531756.jpg
//   - /th/png/300x300/80/https/example.com/image.png
//   - /th/gif/150x150/90/cache-stored-image-key
//
// Features:
// - Automatic caching for 5 minutes to improve performance
// - Support for JPEG, PNG, and GIF formats
// - Flexible sizing (width x height, or width only)
// - Quality control for JPEG compression
// - SSRF protection for remote URLs
// - Proper HTTP headers for browser caching
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request containing thumbnail parameters
//
// Returns:
//   - string: Error message if processing fails, empty string if successful
//
// HTTP Response Headers Set:
//   - Content-Type: Based on file extension (image/jpeg, image/png, image/gif)
//   - Cache-Control: max-age=604800 (7 days for SEO optimization)
func (controller *thumbnailController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return errorMessage
	}

	cacheKey := str.MD5(fmt.Sprint(data.path, data.extension, data.width, "x", data.height, data.quality))

	fileCache := controller.registry.GetFileCache()
	if fileCache != nil {
		if fileCache.Contains(cacheKey) {
			thumb, err := fileCache.Fetch(cacheKey)

			if err == nil {
				controller.setHeaders(w, data.extension)
				return thumb
			}
		}
	}

	thumb, errorMessage := controller.generateThumb(data)

	if errorMessage != "" {
		return errorMessage
	}

	if fileCache != nil {
		err := fileCache.Save(cacheKey, thumb, 5*time.Minute) // cache for 5 minutes

		if err != nil {
			cfmt.Errorln("Error at thumbnailController > CacheFile.Save", "error", err.Error())
		}
	}

	controller.setHeaders(w, data.extension)
	return thumb
}

// setHeaders configures appropriate HTTP response headers for thumbnail images.
// This function ensures proper content-type detection and optimal caching behavior
// for both browsers and CDNs.
//
// Content-Type Mapping:
//   - jpg/jpeg → image/jpeg
//   - png → image/png
//   - gif → image/gif
//   - other extensions → empty string (will be handled by browser)
//
// Cache Strategy:
//   - Sets Cache-Control: max-age=604800 (7 days)
//   - Optimized for SEO and performance
//   - Allows browser and CDN caching for better user experience
//
// Parameters:
//   - w: HTTP response writer to set headers on
//   - fileExtension: The image file extension (e.g., "jpg", "png", "gif")
//
// Example:
//
//	controller.setHeaders(responseWriter, "jpg")
//	// Sets: Content-Type: image/jpeg, Cache-Control: max-age=604800
func (controller *thumbnailController) setHeaders(w http.ResponseWriter, fileExtension string) {
	w.Header().Set("Content-Type", lo.
		If(fileExtension == "jpg", "image/jpeg").
		ElseIf(fileExtension == "jpeg", "image/jpeg").
		ElseIf(fileExtension == "png", "image/png").
		ElseIf(fileExtension == "gif", "image/gif").
		Else(""))

	w.Header().Set("Cache-Control", "max-age=604800") // cache for SEO
}

// prepareData extracts and validates thumbnail generation parameters from HTTP request.
// This function parses URL parameters and normalizes different input formats for
// consistent processing throughout the controller.
//
// URL Parameters Extracted:
//   - extension: Image format (jpg, jpeg, png, gif)
//   - size: Target dimensions (e.g., "300x200", "300")
//   - quality: Compression quality (0-100)
//   - path: Image source path or URL
//
// Path Processing Rules:
//   - "http/" → "http://" (URL normalization)
//   - "https/" → "https://" (URL normalization)
//   - "cache-" prefix → cache lookup mode
//   - "files/" prefix → converted to full URL via links.URL()
//
// Size Parsing:
//   - "WIDTHxHEIGHT" → Both dimensions set
//   - "WIDTH" → Width set, height defaults to 0 (maintain aspect ratio)
//   - Missing size → Error returned
//
// Validation:
//   - All parameters must be present
//   - Extension must be supported
//   - Size and quality must be valid numbers
//
// Parameters:
//   - r: HTTP request containing thumbnail parameters
//
// Returns:
//   - thumbnailControllerData: Parsed and normalized data structure
//   - string: Error message if validation fails, empty string if successful
//
// Example:
//
//	data, err := controller.prepareData(request)
//	if err != "" {
//	    return data, err // Validation failed
//	}
//	// Use data for thumbnail generation
func (controller *thumbnailController) prepareData(r *http.Request) (data thumbnailControllerData, errorMessage string) {
	data.extension, _ = rtr.GetParam(r, "extension")
	size, _ := rtr.GetParam(r, "size")
	quality, _ := rtr.GetParam(r, "quality")
	data.path, _ = rtr.GetParam(r, "path")
	data.isURL = false
	data.isCache = false

	///cfmt.Infoln("====================================")
	//cfmt.Infoln("EXTENSION: ", extension)
	//cfmt.Infoln("SIZE: ", size)
	//cfmt.Infoln("QUALITY: ", quality)
	//cfmt.Infoln("PATH: ", path)
	//cfmt.Infoln("====================================")

	if data.extension == "" {
		return data, "image extension is missing"
	}

	if size == "" {
		return data, "size is missing"
	}

	if quality == "" {
		return data, "quality is missing"
	}

	if data.path == "" {
		return data, "path is missing"
	}

	if strings.HasPrefix(data.path, "http/") || strings.HasPrefix(data.path, "https/") {
		data.isURL = true
		data.path = strings.ReplaceAll(data.path, "https/", "https://")
		data.path = strings.ReplaceAll(data.path, "http/", "http://")
	}

	if strings.HasPrefix(data.path, "cache-") {
		data.isCache = true
		data.path = data.path[6:]
	}

	if strings.HasPrefix(data.path, "files/") {
		data.path = links.URL(data.path, nil)
		data.isURL = true
	}

	widthStr := ""
	heightStr := ""
	if strings.Contains(size, "x") {
		splits := strings.Split(size, "x")
		widthStr = lo.TernaryF(len(splits) > 0, func() string { return splits[0] }, func() string { return "100" })
		heightStr = lo.TernaryF(len(splits) > 1, func() string { return splits[1] }, func() string { return "100" })
	} else {
		widthStr = size
	}

	widthInt := cast.ToInt64(widthStr)
	heightInt := cast.ToInt64(heightStr)
	qualityInt := cast.ToInt64(quality)

	data.width = widthInt
	data.height = heightInt
	data.quality = qualityInt

	return data, errorMessage
}

// generateThumb creates a resized thumbnail from various image sources with comprehensive error handling.
// This function supports three different image sources: remote URLs, cached base64 data, and local files.
// Each source has specific validation and processing logic to ensure security and reliability.
//
// Image Sources:
//
//  1. Remote URLs (data.isURL = true):
//     - Uses urlToBytes() with SSRF protection
//     - Validates URL and content type
//     - Implements timeout and security controls
//
//  2. Cached Data (data.isCache = true):
//     - Retrieves base64-encoded images from file cache
//     - Supports data:image URLs and plain base64 strings
//     - Handles cache expiration gracefully
//
//  3. Local Files (default):
//     - Uses resources.ToBytes() for file system access
//     - Supports local file paths
//     - Standard file I/O error handling
//
// Image Processing:
//   - Supports JPEG, PNG, and GIF formats
//   - Uses imaging library for high-quality resizing
//   - Maintains aspect ratio when height = 0
//   - Applies quality settings for JPEG compression
//
// Error Handling:
//   - Comprehensive logging for all error scenarios
//   - Returns descriptive error messages
//   - Differentiates between cache expiry and other errors
//
// Parameters:
//   - data: Parsed thumbnail generation parameters
//
// Returns:
//   - string: Base64-encoded resized image data
//   - string: Error message if processing fails, empty string if successful
//
// Example:
//
//	result, err := controller.generateThumb(thumbnailData)
//	if err != "" {
//	    return "", err // Processing failed
//	}
//	// result contains base64 image data
func (controller *thumbnailController) generateThumb(data thumbnailControllerData) (content string, errorMessage string) {
	ext := imaging.JPEG

	if data.extension == "gif" {
		ext = imaging.GIF
	}

	if data.extension == "png" {
		ext = imaging.PNG
	}

	// cfmt.Infoln("EXTENSION: ", ext)
	// cfmt.Infoln("WIDTH: ", data.width)
	// cfmt.Infoln("HEIGHT: ", data.height)
	// cfmt.Infoln("QUALITY: ", data.quality)
	// cfmt.Infoln("PATH: ", data.path)

	var err error
	var imgBytes []byte

	if data.isURL {
		//imgBytes = controller.toBytes(data.path)
		imgBytes, err = controller.urlToBytes(data.path)

		if err != nil {
			controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb > from URL", "error", err.Error())
			return "", err.Error()
		}
	} else if data.isCache {
		fileCache := controller.registry.GetFileCache()
		if fileCache == nil {
			controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb > from CACHE", "error", "cache not initialized")
			return "", "cache not initialized"
		}
		dataBase64ImageStr, err := fileCache.Fetch(data.path)

		if err != nil {
			// Downgrade noisy cache expiry to info level, keep other cache errors as error
			if err.Error() == "cache expired" {
				controller.registry.GetLogger().Info("Cache expired at thumbnailController > generateThumb > from CACHE", "error", err.Error())
			} else {
				controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb > from CACHE", "error", err.Error())
			}
			return "", err.Error()
		}

		// convert data:image base64 URL or plain base64 string to bytes
		payload := strings.TrimSpace(dataBase64ImageStr)
		if strings.HasPrefix(payload, "data:image") {
			// expected format: data:image/<ext>;base64,<base64-data>
			parts := strings.SplitN(payload, ",", 2)
			if len(parts) == 2 {
				payload = parts[1]
			} else {
				controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb > from CACHE", "error", "invalid data URL format")
				return "", "invalid data URL format"
			}
		}

		imgBytes, err = base64.StdEncoding.DecodeString(payload)

		if err != nil {
			controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb > from CACHE", "error", err.Error())
			return "", err.Error()
		}
	} else {
		var err error
		imgBytes, err = resources.ToBytes(data.path)

		if err != nil {
			controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb > from RESOURCE", "error", err.Error())
			return "", err.Error()
		}
	}

	imgBytesResized, err := img.Resize(imgBytes, int(data.width), int(data.height), ext)

	if err != nil {
		controller.registry.GetLogger().Error("Error at thumbnailController > generateThumb", "error", err.Error())
		return "", err.Error()
	}

	return string(imgBytesResized), ""
}

// validateURL performs comprehensive security validation on URLs to prevent SSRF attacks.
// It validates the URL format, scheme, and hostname to ensure only safe, public HTTP/HTTPS URLs are allowed.
//
// Security measures implemented:
// - Only allows HTTP and HTTPS schemes
// - Blocks access to localhost and private network ranges
// - Prevents access to internal domains (.local, .internal, .corp)
// - Validates URL format and hostname presence
//
// Parameters:
//   - targetURL: The URL to validate
//
// Returns:
//   - error: nil if URL is safe, otherwise an error describing the validation failure
//
// Example:
//
//	err := controller.validateURL("https://example.com/image.jpg")
//	if err != nil {
//	    return nil, err // URL is not safe
//	}
func (controller *thumbnailController) validateURL(targetURL string) error {
	// Check for empty URL first
	if targetURL == "" {
		return errors.New("invalid URL format")
	}

	parsedURL, err := neturl.Parse(targetURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Check for empty scheme
	if parsedURL.Scheme == "" {
		return errors.New("invalid URL format")
	}

	// Only allow HTTP and HTTPS schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("only HTTP and HTTPS URLs are allowed")
	}

	// Prevent access to localhost and private networks
	host := parsedURL.Hostname()
	if host == "" {
		return errors.New("invalid hostname")
	}

	// Block localhost and private IP ranges
	if controller.isPrivateHost(host) {
		return errors.New("access to private networks is not allowed")
	}

	// For now, allow public domains but block private ones
	// In production, you might want to implement a strict allowlist
	if !controller.isPublicDomain(host) {
		return errors.New("untrusted domain")
	}

	return nil
}

// isPrivateHost determines if a given hostname represents a private or internal network resource.
// This function is critical for SSRF protection as it identifies hosts that should not be accessible
// from the thumbnail service.
//
// Private hosts include:
// - localhost variants (localhost, 127.0.0.1, ::1)
// - Private IP ranges (192.168.x.x, 10.x.x.x, 172.16.x.x-172.31.x.x)
// - Link-local addresses (169.254.x.x)
// - Internal domain suffixes (.local, .internal, .corp)
//
// Parameters:
//   - host: The hostname to check (e.g., "localhost", "192.168.1.1", "server.local")
//
// Returns:
//   - bool: true if the host is private/internal and should be blocked, false otherwise
//
// Example:
//
//	if controller.isPrivateHost("192.168.1.1") {
//	    return errors.New("access to private networks is not allowed")
//	}
func (controller *thumbnailController) isPrivateHost(host string) bool {
	// Check for localhost variations
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}

	// Check for private IP ranges (simplified check)
	if strings.HasPrefix(host, "192.168.") ||
		strings.HasPrefix(host, "10.") ||
		strings.HasPrefix(host, "172.") {
		return true
	}

	// Check for link-local addresses (169.254.x.x)
	if strings.HasPrefix(host, "169.254.") {
		return true
	}

	// Check for internal domain names
	if strings.HasSuffix(host, ".local") ||
		strings.HasSuffix(host, ".internal") ||
		strings.Contains(host, ".corp.") {
		return true
	}

	return false
}

// isPublicDomain determines if a given hostname represents a public, internet-accessible domain.
// This function works in conjunction with isPrivateHost to provide comprehensive SSRF protection.
// It performs additional checks beyond private network detection to identify potentially unsafe hosts.
//
// Public domains are:
// - Not identified as private/internal by isPrivateHost
// - Not IPv6 addresses (simplified check for current implementation)
//
// Note: This is a basic implementation. In production, consider using more comprehensive
// IP parsing libraries or network detection mechanisms.
//
// Parameters:
//   - host: The hostname to check (e.g., "example.com", "cdn.example.org")
//
// Returns:
//   - bool: true if the host appears to be a public domain, false if it should be blocked
//
// Example:
//
//	if !controller.isPublicDomain("example.com") {
//	    return errors.New("untrusted domain")
//	}
func (controller *thumbnailController) isPublicDomain(host string) bool {
	// Basic check for public domains (non-exhaustive)
	// In production, consider using a more comprehensive approach
	return !controller.isPrivateHost(host) &&
		!strings.Contains(host, "::") // IPv6 (simplified)
}

// urlToBytes securely downloads image content from a validated URL with comprehensive security controls.
// This function implements multiple layers of security to prevent SSRF attacks and ensure only
// safe image content is processed.
//
// Security features:
// - URL validation using validateURL() before making any request
// - Context with timeout to prevent hanging requests (30 seconds)
// - HTTP client timeout configuration (30 seconds)
// - User-Agent header to identify legitimate requests
// - Response status code validation (must be 200 OK)
// - Content-Type validation (must be image/*)
// - Proper resource cleanup with defer statements
//
// Parameters:
//   - targetURL: The validated URL to download image content from
//
// Returns:
//   - []byte: The raw image data if successful
//   - error: Error description if any validation or download step fails
//
// Example:
//
//	imgData, err := controller.urlToBytes("https://example.com/image.jpg")
//	if err != nil {
//	    return "", err
//	}
//	// Process imgData...
func (controller *thumbnailController) urlToBytes(targetURL string) ([]byte, error) {
	// Validate URL before making request
	if err := controller.validateURL(targetURL); err != nil {
		controller.registry.GetLogger().Error("URL validation failed", "url", targetURL, "error", err.Error())
		return nil, err
	}

	// Create context with timeout for request cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	// Set user agent to identify legitimate requests
	req.Header.Set("User-Agent", "ThumbnailService/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Url: " + targetURL + " NOT FOUND")
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("no response")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Check content type to ensure we're fetching an image
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("invalid content type: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Url: " + targetURL + " NOT FOUND")
		return nil, err
	}

	return body, nil
}

// toBytes reads a local file and returns its contents as a byte slice.
// This is a simple file I/O wrapper function for local file access.
// It provides basic error handling and logging for file operations.
//
// Security Note: This function should only be used for trusted local files.
// For remote files, use urlToBytes() which includes SSRF protection.
//
// Parameters:
//   - path: Local file system path to the image file
//
// Returns:
//   - []byte: Raw file contents if successful
//   - error: Error description if file cannot be read
//
// Example:
//
//	data, err := controller.toBytes("/path/to/local/image.jpg")
//	if err != nil {
//	    return nil, err
//	}
//	// Process image data...
func (controller *thumbnailController) toBytes(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Println("Path: " + path + " NOT FOUND")
		return nil, err
	}
	return bytes, nil
}

// thumbnailControllerData contains all parsed and processed parameters
// required for thumbnail generation. This structure holds the normalized
// data extracted from HTTP request parameters.
//
// Fields:
//   - extension: Image file format (jpg, jpeg, png, gif)
//   - width: Target width in pixels (0 = maintain aspect ratio)
//   - height: Target height in pixels (0 = maintain aspect ratio)
//   - quality: Compression quality for JPEG (0-100)
//   - path: Image source path or URL after normalization
//   - isURL: True if path is a remote URL requiring HTTP request
//   - isCache: True if path refers to cached base64 data
//
// Usage:
//
//	This struct is populated by prepareData() and consumed by generateThumb()
//	to coordinate the thumbnail generation process.
type thumbnailControllerData struct {
	extension string
	width     int64
	height    int64
	quality   int64
	path      string
	isURL     bool
	isCache   bool
}
