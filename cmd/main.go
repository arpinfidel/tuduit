package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/gateway/wabot"
	"github.com/arpinfidel/tuduit/pkg/cron"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/log"
	checkinrepo "github.com/arpinfidel/tuduit/repo/checkin"
	taskrepo "github.com/arpinfidel/tuduit/repo/task"
	userrepo "github.com/arpinfidel/tuduit/repo/user"
	checkinuc "github.com/arpinfidel/tuduit/usecase/checkin"
	taskuc "github.com/arpinfidel/tuduit/usecase/task"
	useruc "github.com/arpinfidel/tuduit/usecase/user"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"
)

type Start struct {
	ctx    context.Context
	cancel context.CancelFunc

	l  *log.Logger
	db *sqlx.DB
}

func main() {
	main := &Start{}

	main.ctx, main.cancel = context.WithCancel(context.Background())
	time.Local = time.UTC

	l, err := zap.NewDevelopment()
	if err != nil {
		l.Sugar().Fatalf("zap.NewDevelopment: %v", err)
	}
	main.l = log.New(l)

	// app := &cli.App{
	// 	Flags: []cli.Flag{
	// 		&cli.StringFlag{
	// 			Name:    "input",
	// 			Aliases: []string{"i"},
	// 			Usage:   "input",
	// 		},
	// 	},
	// 	Action: func(ctx *cli.Context) error {
	// 		main.l.Infof(ctx.String("input"))
	// 		return nil
	// 	},
	// }

	// if err := app.Run(os.Args); err != nil {
	// 	main.l.Fatalln(err)
	// }
	// ctx := context.Background()

	// db := initPostgres()

	// // init sqlite
	// sqlite, err := sql.Open("sqlite3", "/var/lib/sqlite3/tuduit.db")
	// if err != nil {
	// 	main.l.Fatalln(err)
	// }
	// defer sqlite.Close()

	if _, err := os.Stat("/var/lib/sqlite3/tuduit.db"); os.IsNotExist(err) {
		// Create the file if it does not exist
		_, err = os.Create("/var/lib/sqlite3/tuduit.db")
		if err != nil {
			main.l.Fatalf("failed to create file: %v", err)
		}
	}

	db, err := db.New("sqlite3", "file:/var/lib/sqlite3/tuduit.db?_foreign_keys=on", "file:/var/lib/sqlite3/tuduit.db?_foreign_keys=on")
	if err != nil {
		main.l.Fatalln(err)
	}

	waBot, err := main.initWaBot(db)
	if err != nil {
		main.l.Fatalln(err)
	}
	defer waBot.Disconnect()

	// x := 0
	// err = db.GetMaster().GetContext(main.ctx, &x, "SELECT 1+1")
	// if err != nil {
	// 	main.l.Fatalln(err)
	// }
	// err = db.GetSlave().GetContext(main.ctx, &x, "SELECT COUNT(1) FROM task")
	// if err != nil {
	// 	main.l.Fatalln(err)
	// }
	// fmt.Printf(" >> debug >> x: %#v\n", x)

	taskRepo := taskrepo.New(taskrepo.Dependencies{
		DB: db,
	})
	userRepo := userrepo.New(userrepo.Dependencies{
		DB: db,
	})
	checkinRepo := checkinrepo.New(checkinrepo.Dependencies{
		DB: db,
	})

	taskUC := taskuc.New(taskuc.Dependencies{
		Repo: taskRepo,
	})
	userUC := useruc.New(useruc.Dependencies{
		Repo: userRepo,
	})
	checkinUC := checkinuc.New(checkinuc.Dependencies{
		Repo: checkinRepo,
	})

	cron := cron.New(main.ctx, main.l)

	a := app.New(main.l, app.Dependencies{
		TaskUC:    taskUC,
		UserUC:    userUC,
		CheckinUC: checkinUC,

		Cron: cron,
	}, app.Config{})

	// minioClient := initMinio()

	server := wabot.New(main.ctx, main.l, wabot.Dependencies{
		WaClient: waBot,

		App: a,
	})

	go server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	main.cancel()
}

func (main *Start) initPostgres() *sqlx.DB {
	db, err := sqlx.Open("postgres", "postgres://postgres:@tuduit_pg:5432/tuduit?sslmode=disable")
	if err != nil {
		main.l.Fatalln(err)
	}
	return db
}

func (main *Start) initMinio() *minio.Client {
	endpoint := "tuduit-minio:9000"
	accessKeyID := "MDqjpAZx5SPxynNvANLe"
	secretAccessKey := "XqbapeuYIHBXlVaN6I9CyAZ4ZViqzTlLY3HFwHDy"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		main.l.Fatalln(err)
	}

	// obj, err := minioClient.GetObject(ctx, "test-bucket-1", "test-object-1", minio.GetObjectOptions{})
	// if err != nil {
	// 	main.l.Fatallnln(err)
	// }
	// b, err := io.ReadAll(obj)
	// if err != nil {
	// 	main.l.Fatallnln(err)
	// }

	// main.l.Println(string(b))

	// r := bytes.NewReader([]byte("hello world"))

	// info, err := minioClient.PutObject(ctx, "test-bucket-1", "test-makefolder/test-putobject-2", r, int64(len("hello world")), minio.PutObjectOptions{})
	// if err != nil {
	// 	main.l.Fatallnln(err)
	// }

	return minioClient
}

func (main *Start) initWaBot(db *db.DB) (client *whatsmeow.Client, err error) {
	// https://godocs.io/go.mau.fi/whatsmeow#example-package
	var dbLog waLog.Logger = nil
	// dbLog = waLog.Stdout("Database", "DEBUG", true)
	container := sqlstore.NewWithDB(db.GetMaster().DB, "sqlite3", dbLog)
	err = container.Upgrade()
	if err != nil {
		return nil, err
	}

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
