package queries

const (
	CreateUser = `
		INSERT INTO users (name, email, password) VALUES
		($1, $2, $3);
	`
)