package data

import "time"

// NOTE: any non-exported (not starting with a capital letter) fields
// won't be included when encoding the struct to JSON.
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // hide in JSON
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`    // when released
	Minutes   int32     `json:"minutes,omitzero"` // duration
	Genres    []string  `json:"genres,omitempty"` // romance, comedy, etc.
	Version   int32     `json:"version"`          // starts at 1 and will be incremented each time movie info is updated
}
