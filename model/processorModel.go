package model

type ScrapperResponseModel struct {
	FolderURL string                        `json:"folder_url"`
	Sets      []ScrapperSingleResponseModel `json:"sets"`
}

type ScrapperSingleResponseModel struct {
	Title   string   `json:"title"`
	GroupId int      `json:"group_id"`
	URL     string   `json:"url"`
	Words   []string `json:"words"`
}

type QuizletFolder struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

type MemriseScrapper struct {
	Title string   `json:"title"`
	Urls  []string `json:"urls"`
}

type MemriseSet struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}
