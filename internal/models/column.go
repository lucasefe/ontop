package models

// Column constants for workflow stages
const (
	ColumnInbox      = "inbox"
	ColumnInProgress = "in_progress"
	ColumnDone       = "done"
)

// ValidColumns returns all valid column values
func ValidColumns() []string {
	return []string{ColumnInbox, ColumnInProgress, ColumnDone}
}

// IsValidColumn checks if a column value is valid
func IsValidColumn(column string) bool {
	for _, valid := range ValidColumns() {
		if column == valid {
			return true
		}
	}
	return false
}
