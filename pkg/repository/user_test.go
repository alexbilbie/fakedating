package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"fakedating/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser_GetUserByID(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mariadb",
		ExposedPorts: []string{"3306/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"MARIADB_ROOT_PASSWORD": "example",
			"MARIADB_DATABASE":      "fakedatingtest",
		},
		WaitingFor: wait.ForLog("ready for connections"),
	}
	container, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func(container testcontainers.Container, ctx context.Context) {
		_ = container.Terminate(ctx)
	}(container, ctx)

	mappedPort, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	t.Logf("Waiting for 5 seconds for container to finish booting")
	time.Sleep(time.Second * 5)

	rootDSN := fmt.Sprintf("root:example@tcp(127.0.0.1:%s)/fakedatingtest?multiStatements=true", mappedPort.Port())

	t.Logf("Database DSN: %q", rootDSN)

	db, dbOpenErr := sql.Open("mysql", rootDSN)
	if dbOpenErr != nil {
		t.Fatalf("Failed to open database: %v", dbOpenErr)
	}
	defer db.Close()

	setupDatabase(t, db)

	// Okay lets test now

	repo := NewUser(db)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("helloworld"), bcrypt.DefaultCost)
	initialUser := model.User{
		Email:        "hello@alexbilbie.com",
		PasswordHash: string(hashedPassword),
		Name:         "Alex Bilbie",
		Gender:       model.GenderMale,
		Age:          32,
		Location: model.Location{
			Latitude:  51.507831,
			Longitude: -0.076109,
		},
	}

	createdUser, createErr := repo.Create(initialUser)
	assert.NoError(t, createErr)
	assert.NotEqual(t, createdUser.ID.String(), "")

	fetchedUser, fetchErr := repo.GetByID(createdUser.ID)
	assert.NoError(t, fetchErr)
	assert.Equal(t, fetchedUser.ID.String(), createdUser.ID.String())
}

func setupDatabase(t *testing.T, db *sql.DB) {
	t.Logf("Setting up database")
	f, err := os.Open("../../setup.sql")
	if err != nil {
		t.Fatalf("Failed to open setup.sql: %v", err)
	}

	stmts, readErr := io.ReadAll(f)
	if readErr != nil {
		t.Fatalf("Failed to read setup.sql: %v", readErr)
	}

	_, execErr := db.Exec(string(stmts))
	if execErr != nil {
		t.Fatalf("Failed to setup database: %v", execErr)
	}
	t.Logf("Database setup")
}
