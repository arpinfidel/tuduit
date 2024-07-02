package entity

type Task struct {
	StdFields

	Name        string `json:"name"         db:"name"`
	Description string `json:"description"  db:"description"`
	Status      string `json:"status"       db:"status"`
	Started     bool   `json:"started"      db:"started"`
	StartedAt   string `json:"started_at"   db:"started_at"`
	Completed   bool   `json:"completed"    db:"completed"`
	CompletedAt string `json:"completed_at" db:"completed_at"`
	Archived    bool   `json:"archived"     db:"archived"`
	ArchivedAt  string `json:"archived_at"  db:"archived_at"`
}
