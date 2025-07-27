package registry

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/trueegorletov/analabit/core/source"
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
		currentTime := time.Now().Unix()

		if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
			log.Printf("Scanning cache directory: %s (TTL: %d minutes)", cacheDir, cacheTTL)

			var candidateFiles []struct {
				path      string
				timestamp int64
				age       int64
			}

			err := filepath.WalkDir(cacheDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					log.Printf("Error walking path %s: %v", path, err)
					return err
				}

				if !d.IsDir() && strings.HasSuffix(d.Name(), ".gob") {
					name := strings.TrimSuffix(d.Name(), ".gob")
					ts, parseErr := strconv.ParseInt(name, 10, 64)

					if parseErr != nil {
						log.Printf("Failed to parse timestamp from filename %s: %v", d.Name(), parseErr)
						return nil // Continue walking, don't fail the entire operation
					}

					age := currentTime - ts
					candidateFiles = append(candidateFiles, struct {
						path      string
						timestamp int64
						age       int64
					}{path, ts, age})

					log.Printf("Found cache file: %s (timestamp: %d, age: %d seconds, TTL: %d seconds)",
						d.Name(), ts, age, ttlSeconds)

					// Check if file is within TTL and newer than current latest
					if age < ttlSeconds && ts > latestTimestamp {
						latestTimestamp = ts
						validCacheFile = path
						log.Printf("Selected as latest valid cache: %s", d.Name())
					} else if age >= ttlSeconds {
						log.Printf("Cache file %s is expired (age: %d >= TTL: %d)", d.Name(), age, ttlSeconds)
					} else {
						log.Printf("Cache file %s is older than current latest (%d <= %d)", d.Name(), ts, latestTimestamp)
					}
				}
				return nil
			})

			if err != nil {
				return nil, fmt.Errorf("error walking cache dir: %w", err)
			}

			log.Printf("Cache scan complete. Found %d .gob files, selected: %s",
				len(candidateFiles), validCacheFile)
		} else {
			log.Printf("Cache directory %s does not exist", cacheDir)
		}
	}

	var loadedVarsities []*source.Varsity
	if validCacheFile != "" {
		file, err := os.Open(validCacheFile)
		if err != nil {
			log.Println("Error happened")
			log.Printf("Failed to open cache file: %v", err)
		} else {
			defer file.Close()
			caches, err := source.DeserializeList(file)
			if err != nil {
				log.Printf("Failed to deserialize cache file: %v", err)
			} else {
				var cacheComplete bool
				loadedVarsities, cacheComplete = source.LoadWithCaches(filteredDefs, caches)
				cacheUsed = true

				// If cache was incomplete, save the updated cache
				if !cacheComplete && len(loadedVarsities) > 0 && cacheTTL != -1 {
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
