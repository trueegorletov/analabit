package registry

import (
	"github.com/trueegorletov/analabit/core/source"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CrawlOptions struct {
	VarsitiesList    []string
	VarsitiesExclude []string
	CacheDir         string
	CacheTTLMinutes  int
	DrainStages      []int
	DrainIterations  int
}

type CrawlResult struct {
	LoadedVarsities []*source.Varsity
	CacheUsed       bool
	CacheFile       string
}

// CrawlWithOptions performs crawling and cache lookup, given a set of definitions.
// If cacheTTLMinutes == -1, disables cache lookup and always crawls.
func CrawlWithOptions(defs []source.VarsityDefinition, params CrawlOptions) (*CrawlResult, error) {

	var filteredDefs []source.VarsityDefinition
	varsitiesToUse := make(map[string]bool)
	if len(params.VarsitiesList) == 1 && params.VarsitiesList[0] == "all" {
		for _, def := range defs {
			varsitiesToUse[def.Code] = true
		}
	} else {
		for _, code := range params.VarsitiesList {
			varsitiesToUse[code] = true
		}
	}
	for _, code := range params.VarsitiesExclude {
		delete(varsitiesToUse, code)
	}
	for _, def := range defs {
		if varsitiesToUse[def.Code] {
			filteredDefs = append(filteredDefs, def)
		}
	}
	if len(filteredDefs) == 0 {
		return &CrawlResult{LoadedVarsities: []*source.Varsity{}}, nil
	}

	cacheDir := params.CacheDir
	cacheTTL := params.CacheTTLMinutes
	var validCacheFile string
	var latestTimestamp int64 = -1
	cacheUsed := false

	if cacheTTL != -1 {
		ttlSeconds := int64(cacheTTL * 60)
		if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
			err := filepath.WalkDir(cacheDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && strings.HasSuffix(d.Name(), ".gob") {
					name := strings.TrimSuffix(d.Name(), ".gob")
					ts, err := strconv.ParseInt(name, 10, 64)
					if err == nil && time.Now().Unix()-ts < ttlSeconds && ts > latestTimestamp {
						latestTimestamp = ts
						validCacheFile = path
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("error walking cache dir: %w", err)
			}
		}
	}

	var loadedVarsities []*source.Varsity
	if validCacheFile != "" {
		file, err := os.Open(validCacheFile)
		if err != nil {
			log.Printf("Failed to open cache file: %v", err)
		} else {
			defer file.Close()
			caches, err := source.DeserializeList(file)
			if err != nil {
				log.Printf("Failed to deserialize cache file: %v", err)
			} else {
				loadedVarsities = source.LoadWithCaches(filteredDefs, caches)
				cacheUsed = true
			}
		}
	}
	if len(loadedVarsities) == 0 {
		loadedVarsities = source.LoadFromDefinitions(filteredDefs)
		if len(loadedVarsities) > 0 && cacheTTL != -1 {
			_ = os.MkdirAll(cacheDir, 0755)
			newFile := filepath.Join(cacheDir, fmt.Sprintf("%d.gob", time.Now().Unix()))
			if file, err := os.Create(newFile); err != nil {
				log.Printf("Failed to create cache file: %v", err)
			} else {
				defer file.Close()
				var cachesToSave []*source.VarsityDataCache
				for _, v := range loadedVarsities {
					if v.VarsityDataCache != nil {
						cachesToSave = append(cachesToSave, v.VarsityDataCache)
					}
				}
				if err := source.SerializeList(cachesToSave, file); err != nil {
					log.Printf("Failed to serialize cache list: %v", err)
				}
			}
		}
	}

	sort.Slice(loadedVarsities, func(i, j int) bool {
		return loadedVarsities[i].Name < loadedVarsities[j].Name
	})
	return &CrawlResult{
		LoadedVarsities: loadedVarsities,
		CacheUsed:       cacheUsed,
		CacheFile:       validCacheFile,
	}, nil

}
