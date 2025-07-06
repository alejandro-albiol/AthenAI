package userrole_enum

type UserRole string

const (
	Admin  UserRole = "admin"
	User   UserRole = "user"
	Guest  UserRole = "guest"
)
