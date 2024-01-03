package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"poroto.app/sns/models"
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
	testUser := models.User{Username: "test"}
	if err := testUser.Insert(context.Background(), db, boil.Infer()); err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	fmt.Println(testUser)

	// =====================
	// Fetch
	// =====================
	user, err := models.Users(models.UserWhere.Username.EQ("test")).One(context.Background(), db)
	if err != nil {
		panic(err)
	}
	if diff := cmp.Diff(user.Username, "test"); diff != "" {
		log.Fatalf("user.Username mismatch (-want +got):\n%s", diff)
	}
}

func cleanup(ctx context.Context, db *sql.DB) {
	models.Users().DeleteAll(context.Background(), db)
}
