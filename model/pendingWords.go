package model

type PendingWordModel struct {
	Word     string `json:"word"`
	ListId   uint64 `json:"list_id" db:"list_id"`
	Approved int    `json:"approved"`
	Tried    int    `json:"tried"`
}
