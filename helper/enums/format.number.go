package enums

// Define a custom type for the enum
type FormatCode string

// Explicitly assign string values
const (
	Number      FormatCode = "{N}"
	MonthNumber FormatCode = "{MN}"
	MonthRoman  FormatCode = "{MR}"
	Year        FormatCode = "{Y}"
)
