package enum

type UserRole string

const (
	// Platform-level roles
	PlatformAdmin UserRole = "platform_admin" // Can manage all gyms and platform settings

	// Gym-level roles
	GymAdmin UserRole = "gym_admin" // Can manage their gym, users, billing, and settings
	Trainer  UserRole = "trainer"   // Can work with members, create workouts, view member data
	Member   UserRole = "member"    // Verified, paying gym member with full access
	Guest    UserRole = "guest"     // Temporary/trial access, limited permissions
)

func (r UserRole) IsValid() bool {
	switch r {
	case PlatformAdmin, GymAdmin, Trainer, Member, Guest:
		return true
	}
	return false
}

// IsPlatformLevel returns true if the role has platform-wide permissions
func (r UserRole) IsPlatformLevel() bool {
	return r == PlatformAdmin
}

// IsGymLevel returns true if the role is specific to a gym
func (r UserRole) IsGymLevel() bool {
	return r == GymAdmin || r == Trainer || r == Member || r == Guest
}

// CanManageUsers returns true if the role can manage other users
func (r UserRole) CanManageUsers() bool {
	return r == PlatformAdmin || r == GymAdmin
}

// CanManageGym returns true if the role can manage gym settings
func (r UserRole) CanManageGym() bool {
	return r == PlatformAdmin || r == GymAdmin
}

// CanManageWorkouts returns true if the role can create/manage workouts
func (r UserRole) CanManageWorkouts() bool {
	return r == PlatformAdmin || r == GymAdmin || r == Trainer
}

// CanAccessMemberData returns true if the role can view member workout data
func (r UserRole) CanAccessMemberData() bool {
	return r == PlatformAdmin || r == GymAdmin || r == Trainer
}

// IsVerifiedMember returns true if the role represents a verified gym member
func (r UserRole) IsVerifiedMember() bool {
	return r == Member
}

// HasLimitedAccess returns true if the role has restricted permissions
func (r UserRole) HasLimitedAccess() bool {
	return r == Guest
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
