package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories/queries"
)

type UserRepoOpt struct {
	Db              *sql.DB
	FirestoreClient *firestore.Client
}

type UserRepository interface {
	RegisterUser(ctx context.Context, req entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserIdByTweetId(ctx context.Context, tweetId int64) (*entities.User, error)
	GetUserById(ctx context.Context, userId int64) (*entities.User, error)
	UpdateFcmToken(ctx context.Context, fcmToken string, userId int64) error
	AddUserToFirestore(ctx context.Context, user entities.User) error
	AddFollowerToFirestore(ctx context.Context, userId string, followerId string) error
}

type UserRepositoryDb struct {
	db              *sql.DB
	FirestoreClient *firestore.Client
}

func NewUserRepositoryDb(urOpt *UserRepoOpt) UserRepository {
	return &UserRepositoryDb{
		db:              urOpt.Db,
		FirestoreClient: urOpt.FirestoreClient,
	}
}

func (r *UserRepositoryDb) RegisterUser(ctx context.Context, req entities.User) (*entities.User, error) {
	var err error
	u := entities.User{}

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, queries.CreateUser, req.Name, req.Email, req.Password).Scan(&u.Id, &u.Name, &u.Email)
	} else {
		err = r.db.QueryRowContext(ctx, queries.CreateUser, req.Name, req.Email, req.Password).Scan(&u.Id, &u.Name, &u.Email)
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDb) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	u := entities.User{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, queries.GetUserByEmail, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.FcmToken)
	} else {
		err = r.db.QueryRowContext(ctx, queries.GetUserByEmail, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.FcmToken)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound()
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDb) GetUserIdByTweetId(ctx context.Context, tweetId int64) (*entities.User, error) {
	u := entities.User{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, queries.GetUserIdByTweetId, tweetId).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.FcmToken)
	} else {
		err = r.db.QueryRowContext(ctx, queries.GetUserIdByTweetId, tweetId).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.FcmToken)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound()
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDb) GetUserById(ctx context.Context, userId int64) (*entities.User, error) {
	u := entities.User{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, queries.GetUserById, userId).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.FcmToken)
	} else {
		err = r.db.QueryRowContext(ctx, queries.GetUserById, userId).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.FcmToken)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound()
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDb) UpdateFcmToken(ctx context.Context, fcmToken string, userId int64) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, queries.GetUserById, fcmToken, userId)
	} else {
		_, err = r.db.ExecContext(ctx, queries.GetUserById, fcmToken, userId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryDb) AddUserToFirestore(ctx context.Context, user entities.User) error {
	userIdStr := strconv.Itoa(int(user.Id))
	_, err := r.FirestoreClient.Collection("users").Doc(userIdStr).Set(ctx, entities.UserFirestore{
		Name: user.Name,
		Email: user.Email,
		FollowerCount: 0,
		FollowingCount: 0,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryDb) AddFollowerToFirestore(ctx context.Context, userId string, followerId string) error {
	userRef := r.FirestoreClient.Collection("users").Doc(userId)
	followerRef := r.FirestoreClient.Collection("users").Doc(followerId)
	usersFollowersRef := r.FirestoreClient.Collection("users").Doc(userId).Collection("followers").Doc(followerId)
	followersFollowingRef := r.FirestoreClient.Collection("users").Doc(followerId).Collection("following").Doc(userId)
	
	doc, err := usersFollowersRef.Get(ctx)
	if err != nil {
		return err
	}
	if doc.Exists() {
		return custom_errors.AlreadyFollowed()
	}

	err = r.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		err := t.Set(usersFollowersRef, entities.FollowerFirestore{
			FollowerId: followerId,
			FollowedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		err = t.Set(followersFollowingRef, entities.FollowingFirestore{
			FollowingId: userId,
			FollowedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		err = t.Update(userRef, []firestore.Update{
			{Path: "FollowerCount", Value: firestore.Increment(1)},
		})
		if err != nil {
			return err
		}

		err = t.Update(followerRef, []firestore.Update{
			{Path: "FollowingCount", Value: firestore.Increment(1)},
		})
		if err != nil {
			return err
		}

		return nil
	})
	
	if err != nil {
		return err
	}

	return nil
}
