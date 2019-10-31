package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"

	"github.com/iris-contrib/parrot/parrot-api/api"
	"github.com/iris-contrib/parrot/parrot-api/auth"
	"github.com/iris-contrib/parrot/parrot-api/datastore"
)

func init() {
	// Config log
	golog.SetOutput(os.Stdout)
	golog.SetLevel("info")
}

// TODO: refactor this into cli to start server
func main() {
	// init environment variables
	err := godotenv.Load()
	if err != nil {
		golog.Info(err)
	}

	// init and ping datastore
	dbName := os.Getenv("PARROT_API_DB")
	dbURL := os.Getenv("PARROT_API_DB_URL")
	if dbName == "" || dbURL == "" {
		golog.Fatal("no db set in env")
	}

	ds, err := datastore.NewDatastore(dbName, dbURL)
	if err != nil {
		golog.Fatal(err)
	}
	defer ds.Close()

	// Ping DB until service is up, block meanwhile
	blockAndRetry(5*time.Second, func() bool {
		if err = ds.Ping(); err != nil {
			golog.Error(fmt.Sprintf("failed to ping datastore.\nerr: %s", err))
			return false
		}
		return true
	})

	migrate(dbName, ds)

	app := iris.New()

	app.Use(
		recover.New(),
		func(ctx iris.Context) {
			ctx.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			ctx.Next()
		},
		logger.New(),
	)

	signingKey := os.Getenv("PARROT_AUTH_SIGNING_KEY")
	if signingKey == "" {
		golog.Fatal("no auth signing key set")
	}
	issuerName := os.Getenv("PARROT_AUTH_ISSUER_NAME")
	if signingKey == "" {
		golog.Warn("no auth issuer name set, resorting to default")
		issuerName = "parrot-default"
	}

	tp := auth.TokenProvider{Name: issuerName, SigningKey: []byte(signingKey)}
	app.Configure(auth.NewRouter(ds, tp))
	app.Configure(api.NewRouter(ds, tp))

	// config and init server
	addr := ":8080"
	if os.Getenv("PARROT_API_HOST_PORT") != "" {
		addr = os.Getenv("PARROT_API_HOST_PORT")
	}

	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	app.Run(iris.Server(srv))
}

func migrate(dbName string, ds datastore.Store) {
	migrationStrategy := os.Getenv("PARROT_DB_MIGRATION_STRATEGY")
	if migrationStrategy != "" {
		golog.Infof("migration strategy is set to '%s'", migrationStrategy)
	}

	dirPath := os.Getenv("PARROT_DB_MIGRATIONS_DIR")
	if dirPath == "" {
		dirPath = fmt.Sprintf("./datastore/%s/migrations", dbName)
		golog.Infof("migrations directory not set, using default one: '%s'", dirPath)
	}

	var fn func(string) error

	switch migrationStrategy {
	// Case when we want to start clean each time
	case "down,up":
		fn = func(path string) error {
			err := ds.MigrateDown(path)
			if err != nil {
				return err
			}
			err = ds.MigrateUp(path)
			if err != nil {
				return err
			}
			return nil
		}
	// Case when we want to apply migrations if needed
	case "up":
		fn = ds.MigrateUp
	// Case when we want to simply drop everything
	case "down":
		fn = ds.MigrateDown
	default:
		golog.Fatalf("could not recognize migration strategy '%s'", migrationStrategy)
	}

	golog.Info("migrating...")
	err := fn(dirPath)
	if err != nil {
		golog.Fatal(err)
	}
	golog.Info("migration completed successfully")
}

func blockAndRetry(d time.Duration, fn func() bool) {
	for !fn() {
		golog.Infof("retrying in %s...\n", d.String())
		time.Sleep(d)
	}
}
