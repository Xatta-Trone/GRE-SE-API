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
	// relations
	User      *UserModel `json:"user,omitempty"` // for one2one relations
	WordCount *int       `json:"word_count,omitempty" db:"word_count"`
}

type ListMetaModel struct {
	Id         uint64    `json:"id"`
	UserId     uint64    `json:"user_id,omitempty" db:"user_id"`
	FolderId   *uint64   `json:"folder_id,omitempty" db:"folder_id"`
	Name       string    `json:"name"`
	Url        *string   `json:"url"`
	Words      *string   `json:"words"`
	Visibility int       `json:"visibility"`
	Status     int       `json:"status"`
	CratedAt   time.Time `json:"crated_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ListWordModel struct {
	ListId    uint64 `json:"list_id" db:"list_id"`
	WordId    uint64 `json:"word_id" db:"word_id"`
	WordCount *int   `json:"word_count,omitempty" db:"word_count"`
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
	// relations
	User       *UserModel `json:"user,omitempty"` // for one2one relations
	ListsCount *int       `json:"lists_count,omitempty" db:"lists_count"`
}

type FolderListRelationModel struct {
	FolderId uint64 `json:"folder_id" db:"folder_id"`
	ListId   uint64 `json:"list_id" db:"list_id"`
	ListCount *int   `json:"list_count,omitempty" db:"list_count"`

}
