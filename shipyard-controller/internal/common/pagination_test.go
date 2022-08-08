package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPaginationOnePage checks whether Paginate returns a paginationInfo with a new next page
// key set to 0 for one page
func TestPaginationOnePage(t *testing.T) {
	var pageSize int64 = 20
	nextPageKey := "0"
	paginationInfo := Paginate(18, &pageSize, &nextPageKey)

	assert.Equal(t, paginationInfo.EndIndex, int64(18), "Expect end index to be set to 18")
	assert.Equal(t, paginationInfo.NewNextPageKey, "0", "Expect new next page key to be set to 0")
}

// TestPaginationGetTwoPages checks whether Paginate returns a paginationInfo with a new next page
// key set to 0 when asking for the second out of two pages.
func TestPaginationGetTwoPages(t *testing.T) {
	var pageSize int64 = 20
	nextPageKey := "20"
	paginationInfo := Paginate(21, &pageSize, &nextPageKey)

	assert.Equal(t, paginationInfo.EndIndex, int64(21), "Expect end index to be set to 18")
	assert.Equal(t, paginationInfo.NewNextPageKey, "0", "Expect new next page key to be set to 0")
}

// TestPaginationGetTwoPages checks whether Paginate returns a paginationInfo with a new next page
// key set to 40 when asking for the second out of three pages.
func TestPaginationThreePages(t *testing.T) {
	var pageSize int64 = 20
	nextPageKey := "20"
	paginationInfo := Paginate(41, &pageSize, &nextPageKey)

	assert.Equal(t, paginationInfo.EndIndex, int64(40), "Expect end index to be set to 40")
	assert.Equal(t, paginationInfo.NewNextPageKey, "40", "Expect new next page key to be set to 40")

	paginationInfo = Paginate(41, &pageSize, &paginationInfo.NewNextPageKey)

	assert.Equal(t, paginationInfo.EndIndex, int64(41), "Expect end index to be set to 41")
	assert.Equal(t, paginationInfo.NewNextPageKey, "0", "Expect new next page key to be set to 0")
}
