package interfaces

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserCreationDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersRepository interface {
	// CreateUser creates a new user in the repository.
	CreateUser(user UserCreationDTO) error
	// GetUserByID retrieves a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername retrieves a user by username.
	GetUserByUsername(username string) (User, error)
	// GetUserByEmail retrieves a user by email.
	GetUserByEmail(email string) (User, error)
	// UpdateUser updates an existing user.
	UpdateUser(user User) error
	// DeleteUser removes a user by ID.
	DeleteUser(id string) error
}

type UsersService interface {
	// RegisterUser registers a new user.
	RegisterUser(user UserCreationDTO) error
	// GetUserByID retrieves a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername retrieves a user by username.
	GetUserByUsername(username string) (User, error)
	// GetUserByEmail retrieves a user by email.
	GetUserByEmail(email string) (User, error)
	// UpdateUser updates an existing user.
	UpdateUser(user User) error
	// DeleteUser removes a user by ID.
	DeleteUser(id string) error
}

type UsersHandler interface {
	// RegisterUser handles user registration.
	RegisterUser(dto UserCreationDTO) (User, error)
	// GetUserByID handles retrieving a user by ID.
	GetUserByID(id string) (User, error)
	// GetUserByUsername handles retrieving a user by username.
	GetUserByUsername(username string) (User, error)
	// GetUserByEmail handles retrieving a user by email.
	GetUserByEmail(email string) (User, error)
	// UpdateUser handles updating an existing user.
	UpdateUser(id string, dto UserCreationDTO) (User, error)
	// DeleteUser handles removing a user by ID.
	DeleteUser(id string) error
}
