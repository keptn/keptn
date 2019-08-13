package common

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/keptn/keptn/configuration-service/models"
)

// PaginationResult contains pagination info
type PaginationResult struct {
	// Pointer to next page, base64 encoded
	NewNextPageKey string
	NextPageKey    int64
	// Size of returned page
	PageSize float64

	// End Index
	EndIndex int64
}

// Paginate paginates an array
func Paginate(totalCount int, pageSize *int64, nextPageKeyString *string) *PaginationResult {
	var result = &PaginationResult{}
	var newNextPageKey int64
	var nextPageKey int64

	if nextPageKeyString != nil {
		tmpNextPageKey, _ := strconv.Atoi(*nextPageKeyString)
		nextPageKey = int64(tmpNextPageKey)
		newNextPageKey = nextPageKey + *pageSize
	} else {
		newNextPageKey = *pageSize
	}
	pagesize := *pageSize

	upperLimit := int64(math.Floor(math.Min(float64(totalCount), float64(nextPageKey+pagesize))))
	result.NextPageKey = nextPageKey
	if newNextPageKey < int64(totalCount) {
		nextPageKey = newNextPageKey
	}

	result.NewNextPageKey = strconv.FormatInt(nextPageKey, 10)
	result.EndIndex = upperLimit
	return result
}

// GetPaginatedResources returns a paginates resources set
func GetPaginatedResources(dir string, pageSize *int64, nextPageKey *string) *models.Resources {
	var result = &models.Resources{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Resources:   []*models.Resource{},
	}
	var files = []string{}
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(path, ".git") {
				return nil
			}
			// don't expose the internal directory structure of the container
			cutPrefix := strings.TrimPrefix(strings.TrimPrefix(dir, "./"), "/")
			path = strings.Replace(path, cutPrefix, "", 1)
			if !info.IsDir() {
				files = append(files, strings.TrimPrefix(path, "/"))
			}
			return nil
		})
	if err != nil {
		return result
	}

	paginationInfo := Paginate(len(files), pageSize, nextPageKey)

	totalCount := len(files)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, resourceURI := range files[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			var tmp = resourceURI
			var resource = &models.Resource{ResourceURI: &tmp}
			result.Resources = append(result.Resources, resource)
		}
	}

	result.TotalCount = float64(totalCount)
	result.NextPageKey = paginationInfo.NewNextPageKey

	return result
}
