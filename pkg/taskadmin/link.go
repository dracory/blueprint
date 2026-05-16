package taskadmin

import (
	"net/http"
	"net/url"
)

func link(request http.Request, path string, params map[string]string) string {
	endpoint := request.Context().Value(keyEndpoint).(string)

	params["path"] = path

	return endpoint + query(params)
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := url.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}
