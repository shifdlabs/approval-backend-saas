package enums

// Define a custom type for the enum
type ActionLog string

// Explicitly assign string values
const (
	Create  ActionLog = "Create"
	Update  ActionLog = "Update"
	Delete  ActionLog = "Delete"
	Approve ActionLog = "Approve"
	Reject  ActionLog = "Reject"
	Cancel  ActionLog = "Cancel"
	Recall  ActionLog = "Recall"
	Login   ActionLog = "Login"
	Logout  ActionLog = "Logout"
)
