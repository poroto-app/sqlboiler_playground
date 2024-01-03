package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"poroto.app/sns/models"
)

const (
	numUser      = 10
	numPost      = 2
	numPostImage = 1
)

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&loc=%s&tls=%v&interpolateParams=%v",
		"root",
		"password",
		"localhost:3306",
		"sns",
		"Asia%2FTokyo",
		false,
		true,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	boil.SetDB(db)
	boil.DebugMode = true

	defer db.Close()

	cleanup(context.Background(), db)

	// =====================
	// Insert
	// =====================
	var testUsers []models.User
	for i := 0; i < numUser; i++ {
		testUser := models.User{Username: fmt.Sprintf("test_%d", i)}
		if err := testUser.Insert(context.Background(), db, boil.Infer()); err != nil {
			log.Fatalf("failed to insert user: %v", err)
		}
		testUsers = append(testUsers, testUser)
	}

	var testPosts []models.Post
	for _, testUser := range testUsers {
		for i := 0; i < numPost; i++ {
			testPost := models.Post{UserID: testUser.ID, Content: fmt.Sprintf("My %d post", i)}
			if err := testPost.Insert(context.Background(), db, boil.Infer()); err != nil {
				log.Fatalf("failed to insert post: %v", err)
			}
			testPosts = append(testPosts, testPost)
		}
	}

	for _, testPost := range testPosts {
		for i := 0; i < numPostImage; i++ {
			testPostImage := models.PostImage{PostID: testPost.ID, ImageURL: fmt.Sprintf("https://example.com/%d.jpg", i)}
			if err := testPostImage.Insert(context.Background(), db, boil.Infer()); err != nil {
				log.Fatalf("failed to insert postImage: %v", err)
			}
		}
	}

	// =====================
	// Fetch
	// =====================
	//var userAndPost UserPost
	//if err := models.NewQuery(
	//	models.UserWhere.Username.LIKE("test%"),
	//	qm.Load(models.UserRels.Posts),
	//	qm.Load(models.UserRels.Posts+"."+models.PostRels.PostImages),
	//).Bind(context.Background(), db, &userAndPost); err != nil {
	//	log.Fatalf("failed to fetch userAndPost: %v", err)
	//}
	users, err := models.Users(
		models.UserWhere.Username.LIKE("test%"),
		qm.Load(models.UserRels.Posts),
		qm.Load(models.UserRels.Posts+"."+models.PostRels.PostImages),
	).All(context.Background(), db)
	if err != nil {
		log.Fatalf("failed to fetch userAndPost: %v", err)
	}

	if len(users) != numUser {
		log.Fatalf("wrong user size: expected %d actual %d", len(users), numUser)
	}

	for _, user := range users {
		if len(user.R.Posts) != numPost {
			log.Fatalf("wrong user post size: expected %d actual %d", len(user.R.Posts), numPost)
		}
		for _, post := range user.R.Posts {
			if len(post.R.PostImages) != numPostImage {
				log.Fatalf("wrong post image size: expected %d actual %d", len(post.R.PostImages), numPostImage)
			}

			for _, postImage := range post.R.PostImages {
				log.Println(postImage)
			}
		}
	}
}

func cleanup(ctx context.Context, db *sql.DB) error {
	if _, err := models.PostImages().DeleteAll(ctx, db); err != nil {
		return err
	}
	if _, err := models.Posts().DeleteAll(ctx, db); err != nil {
		return err
	}
	if _, err := models.Users().DeleteAll(ctx, db); err != nil {
		return err
	}
	return nil
}
