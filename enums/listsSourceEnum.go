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

func GetListMetaStatusText(status int) string {
	statusText := ""
	switch status {
	case ListMetaStatusCreated:
		statusText = "Created"
	case ListMetaStatusParsing:
		statusText = "Processing started."
	case ListMetaStatusComplete:
		statusText = "Processing complete."
	case ListMetaStatusError:
		statusText = "Processing error."
	case ListMetaStatusURLError:
		statusText = "Error extracting data from URL"

	}

	return statusText

}
