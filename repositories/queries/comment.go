package queries

const (
	CreateComment = `
		INSERT INTO comments (comment, user_id, tweet_id) VALUES
		($1, $2, $3);
	`
)
