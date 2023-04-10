package enums

const (
	ListSourceWord int = iota
	ListSourceLink
)

const (
	ListVisibilityPublic int = iota + 1
	ListVisibilityMe
)

const (
	ListMetaStatusCreated int = iota
	ListMetaStatusParsing
	ListMetaStatusComplete
	ListMetaStatusError
	ListMetaStatusURLError
)