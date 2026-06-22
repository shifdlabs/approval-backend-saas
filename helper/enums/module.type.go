package enums

// Define a custom type for the enum
type UserRole string

// Explicitly assign string values
const (
	Document           UserRole = "Document"
	DocumentAttachment UserRole = "DocumentAttachment"
	User               UserRole = "User"
	Position           UserRole = "Position"
	NumberingGroup     UserRole = "Numbering Group"
	NumberingNumber    UserRole = "Numbering Number"
	NumberingFormat    UserRole = "Numbering Format"
	DocumentNumbers    UserRole = "Document Numbers"
	Signature          UserRole = "Signature"
	Delegator          UserRole = "Delegator"
	Authentication     UserRole = "Authentication"
	Profile            UserRole = "Profile"
	AppSettingsModule  UserRole = "App Settings"
)
