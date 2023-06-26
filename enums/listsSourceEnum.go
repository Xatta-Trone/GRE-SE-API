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
	ListFilterAll    string = "all"
	ListFilterSaved  string = "saved"
	ListFilterCrated string = "created"
)

const (
	FolderVisibilityPublic int = iota + 1
	FolderVisibilityMe
)

const (
	FolderFilterAll    string = "all"
	FolderFilterSaved  string = "saved"
	FolderFilterCrated string = "created"
)

const (
	ListMetaStatusCreated int = iota
	ListMetaStatusParsing
	ListMetaStatusComplete
	ListMetaStatusError
	ListMetaStatusURLError
)
