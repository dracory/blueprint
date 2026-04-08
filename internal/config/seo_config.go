package config

// ============================================================================
// Interface
// ============================================================================

// SEOConfigInterface defines SEO configuration methods.
type SEOConfigInterface interface {
	SetIndexNowKey(string)
	GetIndexNowKey() string
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetIndexNowKey(v string) {
	c.indexNowKey = v
}

func (c *configImplementation) GetIndexNowKey() string {
	return c.indexNowKey
}
