package handler

// NewUserProvider returns a provider for User related operations.
func NewUserProvider(u userService) (up UserProvider) {
	return UserProvider{u}
}
