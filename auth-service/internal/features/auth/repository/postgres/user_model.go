package usersrepository

// UserModel represents the users table row in the database.
//
// Fields:
//   - ID: unique identifier (UUID stored as text/UUID type in PostgreSQL)
//   - FullName: user's full name (NOT NULL)
//   - PhoneNumber: user's phone number (NULLABLE)
type UserModel struct {
	ID          string  `db:"id"`
	FullName    string  `db:"full_name"`
	PhoneNumber *string `db:"phone_number"`
}