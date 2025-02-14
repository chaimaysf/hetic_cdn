package models

import "time"

// Structure représentant un dossier
type Folder struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ParentID  *int      `json:"parent_id,omitempty"`
}
