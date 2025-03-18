package enum

type SystemRole string

const (
	AdminRole     SystemRole = "admin"
	DeveloperRole SystemRole = "developer"
	ViewerRole    SystemRole = "viewer"
)
