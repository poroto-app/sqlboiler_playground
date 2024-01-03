package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"poroto.app/sns/models"
)

const (
	numUser      = 10000
	numPost      = 1000
	numPostImage = 2
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
	//boil.DebugMode = true

	defer db.Close()

	cleanup(context.Background(), db)

	// =====================
	// Insert
	// =====================
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to begin tx: %v", err)
	}

	var testUsers models.UserSlice
	for i := 0; i < numUser; i++ {
		testUser := models.User{ID: i + 1, Username: fmt.Sprintf("test_%d", i)}
		testUsers = append(testUsers, &testUser)
	}
	if _, err := testUsers.InsertAll(context.Background(), tx, boil.Infer()); err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	var testPosts models.PostSlice
	for iUser, testUser := range testUsers {
		for iPost := 0; iPost < numPost; iPost++ {
			testPost := models.Post{ID: iUser*numPost + iPost + 1, UserID: testUser.ID, Content: fmt.Sprintf("My %d post", iPost)}
			testPosts = append(testPosts, &testPost)
		}
	}
	if _, err := testPosts.InsertAllByPage(context.Background(), tx, boil.Infer()); err != nil {
		log.Fatalf("failed to insert post: %v", err)
	}

	var testPostImages models.PostImageSlice
	for iPost, testPost := range testPosts {
		for iPostImage := 0; iPostImage < numPostImage; iPostImage++ {
			testPostImage := models.PostImage{ID: iPost*numPostImage + iPostImage + 1, PostID: testPost.ID, ImageURL: fmt.Sprintf("https://example.com/%d.jpg", iPostImage)}
			testPostImages = append(testPostImages, &testPostImage)
		}
	}
	if _, err := testPostImages.InsertAllByPage(context.Background(), tx, boil.Infer()); err != nil {
		log.Fatalf("failed to insert postImage: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit tx: %v", err)
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
