package fixtures

import (
	"context"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattes/migrate/source/file"
	"github.com/omer1998/booking_api/config"
	"github.com/omer1998/booking_api/database"
	"github.com/omer1998/booking_api/services"
)

type TestEnv struct {
	Config   config.Config
	Db       database.Database
	T        *testing.T
	ConnPool *pgxpool.Pool
}

func NewTestEnv(t *testing.T, path string) *TestEnv {

	conf, err := config.NewWithEnvPath(path)
	if err != nil {
		t.Error(err)
	}
	pool := database.ConnectPool(context.Background(), conf)
	db := database.NewPostgresDbPool(pool, context.Background())

	return &TestEnv{Db: db, Config: *conf, T: t, ConnPool: pool}
}

func (t *TestEnv) SetUpDb() {
	m, err := migrate.New("file://../../migrations", t.Config.GetConnectionDbUrl())

	if err != nil {
		t.T.Errorf("Error creating migrating instance %s", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.T.Errorf("Error running migration %s", err.Error())
	}
	t.T.Log("Migration up")

}

func (t *TestEnv) TearDownDb() {
	// we will only drop the data from all tables (first for doctor table)
	_, err := t.ConnPool.Exec(context.Background(), "TRUNCATE TABLE doctors CASCADE")
	if err != nil {
		t.T.Errorf("Error truncating table %s", err.Error())
	}
}
func (t *TestEnv) TestAddDoctor(envPath string) {

	testEnv := NewTestEnv(t.T, envPath)
	// perform migration of the init schema (first migration to setup the db)
	// if the db is not already setup
	testEnv.SetUpDb()

	authServ := services.NewAuthenticationService(testEnv.Db)
	if err := authServ.RegisterDoctor("alifaris11@gmail.com", "1234567890"); err != nil {
		t.T.Errorf("Error registering doctor: %v", err)
	} else {
		t.T.Log("Doctor registered")
	}

	// if _, err := authServ.LoginDoctor("omerfaris11@gmail.com", "1234567890"); err != nil {
	// 	t.Errorf("Error logging in doctor: %v", err)
	// } else {
	// 	t.Log("Doctor logged in")
	// }

	t.T.Cleanup(testEnv.TearDownDb)

}
