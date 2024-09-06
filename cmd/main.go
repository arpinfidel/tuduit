package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/gateway/http"
	"github.com/arpinfidel/tuduit/gateway/wabot"
	"github.com/arpinfidel/tuduit/pkg/cron"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/jwt"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/pkg/messenger/whatsapp"
	checkinrepo "github.com/arpinfidel/tuduit/repo/checkin"
	otprepo "github.com/arpinfidel/tuduit/repo/otp"
	schedulerepo "github.com/arpinfidel/tuduit/repo/schedule"
	taskrepo "github.com/arpinfidel/tuduit/repo/task"
	userrepo "github.com/arpinfidel/tuduit/repo/user"
	checkinuc "github.com/arpinfidel/tuduit/usecase/checkin"
	otpuc "github.com/arpinfidel/tuduit/usecase/otp"
	scheduleuc "github.com/arpinfidel/tuduit/usecase/schedule"
	taskuc "github.com/arpinfidel/tuduit/usecase/task"
	useruc "github.com/arpinfidel/tuduit/usecase/user"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Main struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg Config

	l *log.Logger
}

type Config struct {
	Postgres struct {
		Master struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Database string `yaml:"database"`
			SSLMode  string `yaml:"ssl_mode"`
		}
		Slave struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Database string `yaml:"database"`
			SSLMode  string `yaml:"ssl_mode"`
		}
	} `yaml:"postgres"`
	JWT struct {
		SigningMethod string `yaml:"signing_method"`
		PrivateKey    string `yaml:"private_key"`
		PublicKey     string `yaml:"public_key"`
	} `yaml:"jwt"`
}

func main() {
	main := &Main{}

	main.ctx, main.cancel = context.WithCancel(context.Background())
	time.Local = time.UTC

	// read from ./files/etc/secret.yaml
	b, err := os.ReadFile("./files/etc/secret.yaml")
	if err != nil {
		main.l.Fatalln(err)
	}
	err = yaml.Unmarshal(b, &main.cfg)
	if err != nil {
		main.l.Fatalln(err)
	}

	// read from ./files/etc/config.yaml

	zapCfg := zap.NewDevelopmentConfig()
	// zapCfg.Level.SetLevel(zap.WarnLevel)
	l, err := zapCfg.Build()
	if err != nil {
		l.Sugar().Fatalf("zap.NewDevelopment: %v", err)
	}
	main.l = log.New(l)

	main.l.Infof("initializing database")
	db := main.initPostgres()
	main.l.Infof("initialized database")

	main.l.Infof("initializing jwt")
	jwt, err := jwt.New(main.cfg.JWT.SigningMethod, []byte(main.cfg.JWT.PrivateKey), []byte(main.cfg.JWT.PublicKey))
	if err != nil {
		main.l.Fatalln(err)
	}
	main.l.Infof("initialized jwt")

	main.l.Infof("initializing wa bot")
	waBot, err := main.initWaBot(db)
	if err != nil {
		main.l.Fatalln(err)
	}
	defer waBot.Disconnect()
	waClientWrapper := whatsapp.New(whatsapp.Dependencies{
		WaClient: waBot,
	})

	main.l.Infof("initialized wa bot")

	main.l.Infof("initializing repos")
	taskRepo := taskrepo.New(taskrepo.Dependencies{
		DB:     db,
		Logger: main.l,
	})
	scheduleRepo := schedulerepo.New(schedulerepo.Dependencies{
		DB:     db,
		Logger: main.l,
	})
	userRepo := userrepo.New(userrepo.Dependencies{
		DB:     db,
		Logger: main.l,
	})
	checkinRepo := checkinrepo.New(checkinrepo.Dependencies{
		DB:     db,
		Logger: main.l,
	})
	otpRepo := otprepo.New(otprepo.Dependencies{
		DB:     db,
		Logger: main.l,
	})
	main.l.Infof("initialized repos")

	main.l.Infof("initializing usecases")
	taskUC := taskuc.New(taskuc.Dependencies{
		Repo: taskRepo,
	})
	scheduleUC := scheduleuc.New(scheduleuc.Dependencies{
		Repo: scheduleRepo,
	})
	userUC := useruc.New(useruc.Dependencies{
		Repo: userRepo,
	})
	checkinUC := checkinuc.New(checkinuc.Dependencies{
		Repo: checkinRepo,
	})
	otpUC := otpuc.New(otpuc.Dependencies{
		Repo: otpRepo,
	})
	main.l.Infof("initialized usecases")

	main.l.Infof("initializing cron")
	cron := cron.New(main.ctx, main.l)
	main.l.Infof("initialized cron")

	main.l.Infof("creating app")
	a := app.New(main.l, app.Dependencies{
		TaskUC:     taskUC,
		ScheduleUC: scheduleUC,
		UserUC:     userUC,
		CheckinUC:  checkinUC,
		OTPUC:      otpUC,

		Cron: cron,
		JWT:  jwt,

		WaClient: waClientWrapper,

		DB: db,
	}, app.Config{})
	main.l.Infof("created app")

	// minioClient := initMinio()

	main.l.Infof("initializing wabot")
	wabot := wabot.New(main.ctx, main.l, wabot.Dependencies{
		Messenger: waClientWrapper,

		App: a,
	})
	go wabot.Start()
	main.l.Infof("initialized wabot")

	main.l.Infof("initializing http server")
	httpServer := http.New(main.l, http.Dependencies{
		App: a,
	})
	go httpServer.Start()
	main.l.Infof("initialized http server")

	main.l.Infoln("app started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	main.cancel()
}

func (main *Main) initPostgres() *db.DB {
	mcfg := main.cfg.Postgres.Master
	master := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", mcfg.User, mcfg.Password, mcfg.Host, mcfg.Port, mcfg.Database, mcfg.SSLMode)
	scfg := main.cfg.Postgres.Slave
	slave := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", scfg.User, scfg.Password, scfg.Host, scfg.Port, scfg.Database, scfg.SSLMode)
	db, err := db.New("postgres", master, slave)
	if err != nil {
		main.l.Fatalln(err)
	}
	return db
}

func (main *Main) initMinio() *minio.Client {
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

func (main *Main) initWaBot(db *db.DB) (client *whatsmeow.Client, err error) {
	// https://godocs.io/go.mau.fi/whatsmeow#example-package
	var dbLog waLog.Logger = nil
	// dbLog = waLog.Stdout("Database", "DEBUG", true)
	container := sqlstore.NewWithDB(db.GetMaster().DB, "postgres", dbLog)
	err = container.Upgrade()
	if err != nil {
		return nil, err
	}

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	// clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, nil)

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
