package cache_provider

type CacheProvider interface {
	CacheLoader
	CacheChecker
	CacheSaver
	CacheRetriever
	CacheAdder
}

type CacheLoader interface {
	loadCache()
}

type CacheSaver interface {
	SaveCache()
}

type CacheChecker interface {
	ExistsInCache(key string) bool
}

type CacheRetriever interface {
	RetrieveFromCache(key string) string
}

type CacheAdder interface {
	AddToCache(key string, content string)
}
