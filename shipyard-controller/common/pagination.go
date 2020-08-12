package common

import (
	"math"
	"strconv"
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
	if pageSize == nil {
		pageSize = new(int64)
		*pageSize = 20
	}
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
	} else {
		nextPageKey = 0
	}

	result.NewNextPageKey = strconv.FormatInt(nextPageKey, 10)
	result.EndIndex = upperLimit
	return result
}
