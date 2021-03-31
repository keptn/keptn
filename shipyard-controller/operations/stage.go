package operations

type GetStagesParams struct {

	//Pointer to the next set of items
	NextPageKey *string `form:"nextPageKey",json:"nextPageKey"`

	//The number of items to return
	PageSize *int64 `form:"pageSize",json:nextPageKey""`

	//Name of the project
	ProjectName string `form:"-"`
}
