package main

import (
	"log"
	"os"

	"github.com/arpinfidel/tuduit/gateway/http"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "input",
			},
		},
		Action: func(ctx *cli.Context) error {
			log.Println(ctx.String("input"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	// ctx := context.Background()

	// db := initPostgres()

	// // init sqlite
	// sqlite, err := sql.Open("sqlite3", "/var/lib/sqlite3/tuduit.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer sqlite.Close()

	db, err := db.New("sqlite3", "/var/lib/sqlite3/tuduit.db", "/var/lib/sqlite3/tuduit.db")
	if err != nil {
		log.Fatal(err)
	}

	// minioClient := initMinio()

	server := http.Server{}
	server.Start(http.Dependencies{
		DB: db,
	})
}

func initPostgres() *sqlx.DB {
	db, err := sqlx.Open("postgres", "postgres://postgres:@tuduit_pg:5432/tuduit?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func initMinio() *minio.Client {
	endpoint := "tuduit-minio:9000"
	accessKeyID := "MDqjpAZx5SPxynNvANLe"
	secretAccessKey := "XqbapeuYIHBXlVaN6I9CyAZ4ZViqzTlLY3HFwHDy"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// obj, err := minioClient.GetObject(ctx, "test-bucket-1", "test-object-1", minio.GetObjectOptions{})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// b, err := io.ReadAll(obj)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println(string(b))

	// r := bytes.NewReader([]byte("hello world"))

	// info, err := minioClient.PutObject(ctx, "test-bucket-1", "test-makefolder/test-putobject-2", r, int64(len("hello world")), minio.PutObjectOptions{})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	return minioClient
}
