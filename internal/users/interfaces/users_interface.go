package interfaces

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersRepository interface {
	// CreateUser creates a new user in the repository.
	CreateUser(user User) error
	// GetUserByID retrieves a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername retrieves a user by username.
	GetUserByUsername(username string) (User, error)
	// UpdateUser updates an existing user.
	UpdateUser(user User) error
	// DeleteUser removes a user by ID.
	DeleteUser(id string) error
}

type UsersService interface {
	// RegisterUser registers a new user.
	RegisterUser(user User) error
	// GetUserByID retrieves a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername retrieves a user by username.
	GetUserByUsername(username string) (User, error)
	// UpdateUser updates an existing user.
	UpdateUser(user User) error
	// DeleteUser removes a user by ID.
	DeleteUser(id string) error
}

type UsersController interface {
	// RegisterUser handles user registration.
	RegisterUser(user User) error
	// GetUserByID handles retrieving a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername handles retrieving a user by username.
	GetUserByUsername(username string) (User, error)
	// UpdateUser handles updating an existing user.
	UpdateUser(user User) error
	// DeleteUser handles removing a user by ID.
	DeleteUser(id string) error
}