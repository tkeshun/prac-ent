package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tutorial/ent"
	"tutorial/ent/user"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// PostgreSQLの接続情報
	postgresDSN := "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"

	client := Open(postgresDSN)
	// client, err := ent.Open(dialect.Postgres, postgresDSN)

	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// u, err := CreateUser(context.TODO(), client)
	// if err != nil {
	// 	log.Fatal("Error")
	// }

	car, err := CreateCar(context.Background(), client)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(car)

	u, err := QueryUser(context.TODO(), client)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(u)
}

// Open new connection
func Open(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to Create")
	}
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) ([]*ent.User, error) {
	// u, err := client.User.Query().Where(user.Name("a8m")).Only(ctx) // レコードが１つじゃないとエラーになる

	// go run main.go
	// 2024/11/21 22:56:23 user returned:  [User(id=1, age=30, name=a8m) User(id=2, age=30, name=a8m)]
	// [User(id=1, age=30, name=a8m) User(id=2, age=30, name=a8m)]
	u, err := client.User.Query().Where(user.Name("a8m")).Limit(2).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	log.Println("user returned: ", u)
	return u, nil
}

func CreateCar(ctx context.Context, client *ent.Client) (*ent.Car, error) {
	car, err := client.Car.
		Create().
		SetModel("Tesla").SetCarNumber("111111").SetOwnerID(1).Save(ctx)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("err: %w", err)
	}
	return car, nil
}
