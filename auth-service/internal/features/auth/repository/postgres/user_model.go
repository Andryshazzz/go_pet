package usersrepository

// UserModel represents the users table row in the database.
//
// Fields:
//   - ID: unique identifier (UUID)
//   - PhoneNumber: phone number used as login (UNIQUE, NOT NULL)
//   - PasswordHash: bcrypt hash of the password (NOT NULL)
//   - FullName: user's full name (NOT NULL)
type UserModel struct {
	ID           string `db:"id"`
	PhoneNumber  string `db:"phone_number"`
	PasswordHash string `db:"password_hash"`
	FullName     string `db:"full_name"`
}