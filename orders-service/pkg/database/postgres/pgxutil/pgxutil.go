package utils

// Scanner is a common interface implemented by both pgx.Row and pgx.Rows.
// It allows writing a single scan helper that works for both single-row
// queries (QueryRow) and multi-row queries (Query).
//
// Usage:
//
//	func scanUser(s pgxutil.Scanner) (entity.User, error) {
//	    var u entity.User
//	    if err := s.Scan(&u.ID, &u.Name); err != nil {
//	        return entity.User{}, err
//	    }
//	    return u, nil
//	}
type Scanner interface {
	Scan(dest ...any) error
}
