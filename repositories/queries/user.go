package queries

const (
	CreateUser = `
		INSERT INTO users (name, email, password) VALUES
		($1, $2, $3);
	`

	GetUserByEmail = `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1 AND deleted_at IS NULL;
	`
)
