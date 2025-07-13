package userrole_enum

type UserRole string

const (
	Admin UserRole = "admin"
	User  UserRole = "user"
	Guest UserRole = "guest"
)

type UserVerificationStatus string

const (
	Verified   UserVerificationStatus = "verified"
	Unverified UserVerificationStatus = "unverified"
	Demo       UserVerificationStatus = "demo"
)
