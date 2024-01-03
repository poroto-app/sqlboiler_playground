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
	"time"
)

const (
	numUser      = 1000
	numPost      = 100
	numPostImage = 2
)

type UserPost struct {
	models.User      `boil:"users,bind"`
	models.Post      `boil:"posts,bind"`
	models.PostImage `boil:"post_images,bind"`
}

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

	log.Printf("start insert users")
	start := time.Now()
	var testUsers models.UserSlice
	for i := 0; i < numUser; i++ {
		testUser := models.User{ID: i + 1, Username: fmt.Sprintf("test_%d", i)}
		testUsers = append(testUsers, &testUser)
	}
	if _, err := testUsers.InsertAll(context.Background(), tx, boil.Infer()); err != nil {
		log.Fatalf("failed to insert users: %v", err)
	}
	log.Printf("insert user elapsed: %v", time.Since(start))

	log.Printf("start insert posts")
	start = time.Now()
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
	log.Printf("insert post elapsed: %v", time.Since(start))

	log.Printf("start insert postImages")
	start = time.Now()
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
	log.Printf("insert postImage elapsed: %v", time.Since(start))

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit tx: %v", err)
	}

	// =====================
	// Fetch(By Load)
	// =====================
	start = time.Now()
	_, err = models.Users(
		models.UserWhere.ID.GTE(100),
		models.UserWhere.ID.LT(200),
		qm.Load(models.UserRels.Posts),
		qm.Load(models.UserRels.Posts+"."+models.PostRels.PostImages),
	).All(context.Background(), db)
	if err != nil {
		log.Fatalf("failed to fetch userAndPost: %v", err)
	}
	log.Printf("[Load] elapsed: %v", time.Since(start))

	// =====================
	// Fetch(By JOIN)
	// =====================
	start = time.Now()
	var userPosts []UserPost
	if err := models.NewQuery(
		qm.Select("users.*", "posts.*", "post_images.*"),
		qm.From("users"),
		qm.InnerJoin("posts on posts.user_id = users.id"),
		qm.InnerJoin("post_images on post_images.post_id = posts.id"),
		models.UserWhere.ID.GTE(100),
		models.UserWhere.ID.LT(200),
	).Bind(context.Background(), db, &userPosts); err != nil {
		log.Fatalf("failed to fetch userAndPost: %v", err)
	}
	log.Printf("[JOIN] elapsed: %v", time.Since(start))
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
