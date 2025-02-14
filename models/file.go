package models

import "time"

type File struct {
	ID        int       `json:"id"`
	UserID    *int      `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	IsFolder  bool      `json:"is_folder"`
	ParentID  *int      `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
	FolderID  *int      `json:"folder_id,omitempty"`
}