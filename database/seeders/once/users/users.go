package usersonce

import "github.com/dracory/userstore"

// defaultUser describes one account to seed.
type defaultUser struct {
	email     string
	firstName string
	lastName  string
	password  string
	role      string
}

// defaultUsers returns the accounts created by UsersSeed on first boot.
// After the seed has run the database owns the rows — edits here do not
// propagate.
func defaultUsers() []defaultUser {
	return []defaultUser{
		{
			email:     "admin@example.com",
			firstName: "Admin",
			lastName:  "User",
			password:  "password",
			role:      userstore.USER_ROLE_ADMINISTRATOR,
		},
		{
			email:     "user@example.com",
			firstName: "Test",
			lastName:  "User",
			password:  "password",
			role:      userstore.USER_ROLE_USER,
		},
	}
}
