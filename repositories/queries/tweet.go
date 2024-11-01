package queries

const (
	CreateTweet = `
		INSERT INTO users (tweet, user_id) VALUES
		($1, $2);
	`
)
