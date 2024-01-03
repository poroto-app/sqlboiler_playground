package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"poroto.app/sns/models"
)

type UserAndPost struct {
	models.User `boil:"users,bind"`
	models.Post `boil:"posts,bind"`
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
	boil.DebugMode = true

	defer db.Close()

	cleanup(context.Background(), db)

	// =====================
	// Insert
	// =====================
	testUser := models.User{Username: "test"}
	if err := testUser.Insert(context.Background(), db, boil.Infer()); err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	testPost := models.Post{UserID: testUser.ID, Content: "My first post"}
	if err := testPost.Insert(context.Background(), db, boil.Infer()); err != nil {
		log.Fatalf("failed to insert post: %v", err)
	}

	// =====================
	// Fetch
	// =====================
	var userAndPost UserAndPost
	if err := models.NewQuery(
		qm.Select("users.*", "posts.*"),
		qm.From("users"),
		qm.InnerJoin("posts on posts.user_id = users.id"),
		qm.Where("users.id = ?", testUser.ID),
	).Bind(context.Background(), db, &userAndPost); err != nil {
		log.Fatalf("failed to fetch userAndPost: %v", err)
	}

	if diff := cmp.Diff(testPost, userAndPost.Post); diff != "" {
		log.Fatalf("post mismatch  (-want +got):\n%s", diff)
	}
}

func cleanup(ctx context.Context, db *sql.DB) error {
	if _, err := models.Posts().DeleteAll(ctx, db); err != nil {
		return err
	}
	if _, err := models.Users().DeleteAll(ctx, db); err != nil {
		return err
	}
	return nil
}
