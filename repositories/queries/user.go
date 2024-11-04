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

	GetUserIdByTweetId = `
		SELECT u.id
		FROM users u
		JOIN tweets t on u.id = t.user_id
		WHERE t.id = $1 AND t.deleted_at IS NULL;
	`

	GetUserById = `
		SELECT id, name, email, password
		FROM users
		WHERE id = $1 AND deleted_at IS NULL;
	`
)
