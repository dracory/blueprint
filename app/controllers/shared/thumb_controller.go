package shared

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"project/app/links"
	"project/config"
	"project/resources"
	"time"

	"strings"

	"github.com/dracory/base/str"

	"github.com/disintegration/imaging"
	"github.com/dracory/base/img"
	"github.com/go-chi/chi/v5"
	"github.com/gouniverse/router"
	"github.com/spf13/cast"

	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
)

// == CONSTRUCTOR =============================================================

func NewThumbController() router.HTMLControllerInterface {
	return &thumbnailController{}
}

// == CONTROLLER ==============================================================

type thumbnailController struct{}

// ThumbnailHandler
// ================================================================
// Resizes local images to the specified width and height
// ================================================================
// Path
// /th/EXT/WIDTHxHEIGHT/QUALITY/path
// Example:
// /th/jpg/1920x0/70/images/backgrounds/pexels-pixabay-531756.jpg
// ================================================================
func (controller *thumbnailController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return errorMessage
	}

	cacheKey := str.MD5(fmt.Sprint(data.path, data.extension, data.width, "x", data.height, data.quality))

	if config.CacheFile != nil {
		if config.CacheFile.Contains(cacheKey) {
			thumb, err := config.CacheFile.Fetch(cacheKey)

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

	if config.CacheFile != nil {
		err := config.CacheFile.Save(cacheKey, thumb, 5*time.Minute) // cache for 5 minutes

		if err != nil {
			cfmt.Errorln("Error at thumbnailController > CacheFile.Save", "error", err.Error())
		}
	}

	controller.setHeaders(w, data.extension)
	return thumb
}

func (controller *thumbnailController) setHeaders(w http.ResponseWriter, fileExtension string) {
	w.Header().Set("Content-Type", lo.
		If(fileExtension == "jpg", "image/jpeg").
		ElseIf(fileExtension == "jpeg", "image/jpeg").
		ElseIf(fileExtension == "png", "image/png").
		ElseIf(fileExtension == "gif", "image/gif").
		Else(""))

	w.Header().Set("Cache-Control", "max-age=604800") // cache for SEO
}

func (controller *thumbnailController) prepareData(r *http.Request) (data thumbnailControllerData, errorMessage string) {
	data.extension = chi.URLParam(r, "extension")
	size := chi.URLParam(r, "size")
	quality := chi.URLParam(r, "quality")
	data.path = chi.URLParam(r, "*")
	data.isURL = false

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
			config.Logger.Error("Error at thumbnailController > generateThumb > from URL", "error", err.Error())
			return "", err.Error()
		}
	} else {
		var err error
		imgBytes, err = resources.ToBytes(data.path)

		if err != nil {
			config.Logger.Error("Error at thumbnailController > generateThumb > from RESOURCE", "error", err.Error())
			return "", err.Error()
		}
	}

	imgBytesResized, err := img.Resize(imgBytes, int(data.width), int(data.height), ext)

	if err != nil {
		config.Logger.Error("Error at thumbnailController > generateThumb", "error", err.Error())
		return "", err.Error()
	}

	return string(imgBytesResized), ""
}

func (controller *thumbnailController) urlToBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Url: " + url + " NOT FOUND")
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("no response")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Url: " + url + " NOT FOUND")
		return nil, err
	}

	return body, nil
}

func (controller *thumbnailController) toBytes(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Path: " + path + " NOT FOUND")
		return nil, err
	}
	return bytes, nil
}

type thumbnailControllerData struct {
	extension string
	width     int64
	height    int64
	quality   int64
	path      string
	isURL     bool
}
