package queries

const (
	CreateTweet = `
		INSERT INTO tweets (tweet, user_id) VALUES
		($1, $2);
	`
)
