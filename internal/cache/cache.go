package cache

import (
	"github.com/faabiosr/cachego"
	"github.com/jellydator/ttlcache/v3"
)

// Memory is the in-memory TTL cache used for transient runtime data.
var Memory *ttlcache.Cache[string, any]

// File is the filesystem-backed cache used for binary/string blobs like thumbnails.
var File cachego.Cache
