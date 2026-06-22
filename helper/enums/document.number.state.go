package enums

// Define a custom type for the enum
type DocumentNumberState int

// Explicitly assign string values
const (
	Cancelled DocumentNumberState = 0
	Booked    DocumentNumberState = 1
	Saved     DocumentNumberState = 2
)
