package cache_provider

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type LocalCacheProvider struct {
	logger          *log.Logger
	cacheFile       string
	cache           map[string]string
	loadedCacheSize int
}

func NewLocalCacheProvider(logger *log.Logger, cacheFile string) (CacheProvider, error) {
	c := &LocalCacheProvider{
		logger:    logger,
		cacheFile: cacheFile,
		cache:     make(map[string]string),
	}
	c.loadCache()
	return c, nil
}

type Schemas struct {
	Schemas []Schema `json:"schemas"`
}

type Schema struct {
	Id      string          `json:"id"`
	Content json.RawMessage `json:"content"`
}

func (c *LocalCacheProvider) SaveCache() {
	//only save cache if it has gained new entries
	if len(c.cache) > c.loadedCacheSize {
		var schemas Schemas
		for k, v := range c.cache {
			{
				schema := Schema{
					Id:      k,
					Content: json.RawMessage(v),
				}
				schemas.Schemas = append(schemas.Schemas, schema)
			}
		}
		j, err := json.MarshalIndent(schemas, "", "  ")
		if err != nil {
			c.logger.Printf("error marshalling cache: %s", err.Error())
		}
		c.logger.Printf("saving cache")
		err = ioutil.WriteFile(c.cacheFile, j, 0644)
		if err != nil {
			c.logger.Printf("error saving cache: %s", err.Error())
		}
	} else {
		c.logger.Printf("cache unchanged, not saving")
	}
}

func (c *LocalCacheProvider) loadCache() {
	c.cache = make(map[string]string)
	jsonFile, err := os.Open(c.cacheFile)
	if err != nil {
		c.logger.Printf("error loading cache file %s", c.cacheFile)
	} else {
		var schemas Schemas
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			c.logger.Printf("error reading schema file: %s", err.Error())
		}
		err = json.Unmarshal(byteValue, &schemas)
		if err != nil {
			c.logger.Printf("error unmarshalling json: %s", err.Error())
		}
		for _, s := range schemas.Schemas {
			if s.Id != "" {
				content, err := json.Marshal(s.Content)
				if err != nil {
					c.logger.Printf("error marshalling content: %s", err.Error())
				}
				c.cache[s.Id] = string(content)
				c.loadedCacheSize++
				c.logger.Printf("loaded schema %s into cache", s.Id)
			} else {
				c.logger.Printf("failed to map from %s", c.cacheFile)
			}
		}
	}
	defer jsonFile.Close()
}

func (c *LocalCacheProvider) ExistsInCache(key string) bool {
	_, exists := c.cache[key]
	return exists
}

func (c *LocalCacheProvider) RetrieveFromCache(key string) string {
	c.logger.Printf("loading %s from cache", key)
	return string(c.cache[key])
}

func (c *LocalCacheProvider) AddToCache(key string, content string) {
	c.cache[key] = content
	c.logger.Printf("added %s to cache", key)
}
