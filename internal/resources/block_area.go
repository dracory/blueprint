package resources

func BlockAreaJS() string {
	resource, err := Resource("js/blockarea_v0200.js")

	if err != nil {
		return ""
	}

	return resource
}
