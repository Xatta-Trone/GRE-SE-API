package model

import "time"

type ListModel struct {
	Id         uint64    `json:"id"`
	UserId     uint64    `json:"user_id" db:"user_id"`
	ListMetaId *uint64   `json:"list_meta_id" db:"list_meta_id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Visibility int       `json:"visibility"`
	Status     int       `json:"status"`
	CratedAt   time.Time `json:"crated_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ListMetaModel struct {
	Id         uint64    `json:"id"`
	UserId     uint64    `json:"user_id,omitempty" db:"user_id"`
	Name       string    `json:"name"`
	Url        *string   `json:"url"`
	Words      *string   `json:"words"`
	Visibility int       `json:"visibility"`
	Status     int       `json:"status"`
	CratedAt   time.Time `json:"crated_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ListWordModel struct {
	ListId uint64 `json:"list_id" db:"list_id"`
	WordId uint64 `json:"word_id" db:"word_id"`
}

type FolderModel struct {
	Id         uint64    `json:"id"`
	UserId     uint64    `json:"user_id" db:"user_id"`
	ListMetaId *uint64   `json:"list_meta_id" db:"list_meta_id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Visibility int       `json:"visibility"`
	Status     int       `json:"status"`
	CratedAt   time.Time `json:"crated_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type FolderListRelationModel struct {
	FolderId uint64 `json:"folder_id" db:"folder_id"`
	ListId   uint64 `json:"list_id" db:"list_id"`
}
