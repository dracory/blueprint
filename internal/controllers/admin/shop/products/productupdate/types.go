package productupdate

import (
	"net/http"
	"strings"

	"github.com/dracory/req"
	"github.com/samber/lo"
)

// MetadataRequest represents the JSON request for saving metadata
type MetadataRequest struct {
	Metadata []MetadataItem `json:"metadata"`
}

// MetadataResponse represents the JSON response for metadata
type MetadataResponse struct {
	Metadata []MetadataItem `json:"metadata"`
	Message  string         `json:"message"`
}

// MetadataItem represents a single metadata entry
type MetadataItem struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TagItem represents a single tag entry
type TagItem struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

// TagRequest represents the JSON request for saving tags
type TagRequest struct {
	Tags []TagItem `json:"tags"`
}

// ReqArrayOfMaps extracts an array of maps from request parameters
func ReqArrayOfMaps(r *http.Request, key string, defaultValue []map[string]string) []map[string]string {
	all := req.GetAll(r)

	reqArrayOfMaps := []map[string]string{}

	if all == nil {
		return reqArrayOfMaps
	}

	mapIndexMap := map[string]map[string]string{}

	for k, v := range all {
		if !strings.HasPrefix(k, key+"[") {
			continue
		}
		if !strings.HasSuffix(k, "]") {
			continue
		}
		if !strings.Contains(k, "][") {
			continue
		}
		if len(v) != 1 {
			continue
		}

		mapValue := v[0]

		str := strings.TrimSuffix(strings.TrimPrefix(k, key+"["), "]")
		split := strings.Split(str, "][")

		if len(split) != 2 {
			continue
		}

		index, key := split[0], split[1]

		if lo.HasKey(mapIndexMap, index) {
			if mapIndexMap[index] == nil {
				mapIndexMap[index] = map[string]string{}
			}
			mapIndexMap[index][key] = mapValue
		} else {
			mapIndexMap[index] = map[string]string{
				key: mapValue,
			}
		}
	}

	for _, v := range mapIndexMap {
		if v == nil {
			continue
		}
		reqArrayOfMaps = append(reqArrayOfMaps, v)
	}

	return reqArrayOfMaps
}
