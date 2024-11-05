DROP DATABASE IF EXISTS moana_healthcare_db; 
CREATE DATABASE moana_healthcare_db;

create table users (
	id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    photo_profile VARCHAR,
    fcm_token VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

create table tweets (
	id BIGSERIAL PRIMARY KEY,
	tweet VARCHAR not null,
	user_id BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

create table comments (
	id BIGSERIAL PRIMARY KEY,
	comment VARCHAR not null,
	user_id BIGINT NOT NULL,
	tweet_id BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (tweet_id) REFERENCES tweets(id)
);
