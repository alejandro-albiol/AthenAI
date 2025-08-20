package enum

type UserRole string

const (
	Admin UserRole = "admin"
	User  UserRole = "user"
	Guest UserRole = "guest"
)

func (r UserRole) IsValid() bool {
	switch r {
	case Admin, User, Guest:
		return true
	}
	return false
}

type UserVerificationStatus string

const (
	Verified   UserVerificationStatus = "verified"
	Unverified UserVerificationStatus = "unverified"
	Demo       UserVerificationStatus = "demo"
)

func (v UserVerificationStatus) IsValid() bool {
	switch v {
	case Verified, Unverified, Demo:
		return true
	}
	return false
}
