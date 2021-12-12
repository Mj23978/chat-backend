package model

type FilterOperator string
type SortOption int32
type MessageStatus int32

const (
	Equal          FilterOperator = "$eq"
	NotEqual       FilterOperator = "$ne"
	Greater        FilterOperator = "$gt"
	GreaterOrEqual FilterOperator = "$gte"
	Less           FilterOperator = "$lt"
	LessOrEqual    FilterOperator = "$lte"
	In             FilterOperator = "$in"
	NotIn          FilterOperator = "$nin"
	Query          FilterOperator = "$q"
	AutoComplete   FilterOperator = "$autocompltete"
	Exists         FilterOperator = "$exists"
	And            FilterOperator = "$and"
	Or             FilterOperator = "$or"
	Nor            FilterOperator = "$nor"
	Contains       FilterOperator = "$contains"
)

const (
	ASC SortOption = 1
	DES SortOption = -1
)

const (
	Sending MessageStatus = iota
	Updating
	Deleting
	MSFailed
	FailedUpdate
	FailedDelete
	Sent
	NotSent
	NotViewed
	Viewed
)
