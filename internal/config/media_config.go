package config

// ============================================================================
// Interface
// ============================================================================

// MediaConfigInterface defines media/storage configuration methods.
type MediaConfigInterface interface {
	SetMediaBucket(string)
	GetMediaBucket() string

	SetMediaDriver(string)
	GetMediaDriver() string

	SetMediaKey(string)
	GetMediaKey() string

	SetMediaEndpoint(string)
	GetMediaEndpoint() string

	SetMediaRegion(string)
	GetMediaRegion() string

	SetMediaRoot(string)
	GetMediaRoot() string

	SetMediaSecret(string)
	GetMediaSecret() string

	SetMediaUrl(string)
	GetMediaUrl() string
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetMediaBucket(v string) {
	c.mediaBucket = v
}

func (c *configImplementation) GetMediaBucket() string {
	return c.mediaBucket
}

func (c *configImplementation) SetMediaDriver(v string) {
	c.mediaDriver = v
}

func (c *configImplementation) GetMediaDriver() string {
	return c.mediaDriver
}

func (c *configImplementation) SetMediaKey(v string) {
	c.mediaKey = v
}

func (c *configImplementation) GetMediaKey() string {
	return c.mediaKey
}

func (c *configImplementation) SetMediaEndpoint(v string) {
	c.mediaEndpoint = v
}

func (c *configImplementation) GetMediaEndpoint() string {
	return c.mediaEndpoint
}

func (c *configImplementation) SetMediaRegion(v string) {
	c.mediaRegion = v
}

func (c *configImplementation) GetMediaRegion() string {
	return c.mediaRegion
}

func (c *configImplementation) SetMediaRoot(v string) {
	c.mediaRoot = v
}

func (c *configImplementation) GetMediaRoot() string {
	return c.mediaRoot
}

func (c *configImplementation) SetMediaSecret(v string) {
	c.mediaSecret = v
}

func (c *configImplementation) GetMediaSecret() string {
	return c.mediaSecret
}

func (c *configImplementation) SetMediaUrl(v string) {
	c.mediaUrl = v
}

func (c *configImplementation) GetMediaUrl() string {
	return c.mediaUrl
}
