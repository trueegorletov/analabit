package registry

import (
	"analabit/core/source"
	"fmt"
	"io/fs"
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
					nameWithoutExt := strings.TrimSuffix(d.Name(), ".gob")
					ts, err := strconv.ParseInt(nameWithoutExt, 10, 64)
					if err == nil {
						if time.Now().Unix()-ts < ttlSeconds && ts > latestTimestamp {
							latestTimestamp = ts
							validCacheFile = path
						}
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
		if err == nil {
			defer file.Close()
			caches, err := source.DeserializeList(file)
			if err == nil {
				loadedVarsities = source.LoadWithCaches(filteredDefs, caches)
				cacheUsed = true
			}
		}
	}
	if loadedVarsities == nil || len(loadedVarsities) == 0 {
		loadedVarsities = source.LoadFromDefinitions(filteredDefs)
		if len(loadedVarsities) > 0 && cacheTTL != -1 {
			_ = os.MkdirAll(cacheDir, 0755)
			newCacheFilename := filepath.Join(cacheDir, fmt.Sprintf("%d.gob", time.Now().Unix()))
			file, err := os.Create(newCacheFilename)
			if err == nil {
				defer file.Close()
				var cachesToSave []*source.VarsityDataCache
				for _, v := range loadedVarsities {
					if v.VarsityDataCache != nil {
						cachesToSave = append(cachesToSave, v.VarsityDataCache)
					}
				}
				_ = source.SerializeList(cachesToSave, file)
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
