package util

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbName     = "mysql"
	dbPortNat  = nat.Port("3306/tcp")
	mysqlImage = "mysql:8.0"
)

func NewTestDB(ctx context.Context) (*sql.DB, error) {
	mysqlC, err := createMySQLContainer(ctx)
	if err != nil {
		return nil, err
	}

	db, err := createDBConnection(ctx, mysqlC)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createMySQLContainer(ctx context.Context) (testcontainers.Container, error) {
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: mysqlImage,
			Env: map[string]string{
				"MYSQL_DATABASE":             dbName,
				"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
			},
			ExposedPorts: []string{"3306/tcp"},
			Tmpfs:        map[string]string{"/var/lib/mysql": "rw"},
			WaitingFor:   wait.ForLog("port: 3306  MySQL Community Server"),
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      "../migrations/create_user.sql",
					ContainerFilePath: "/docker-entrypoint-initdb.d/create_user.sql",
					FileMode:          644,
				},
			},
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}
	return mysqlC, nil
}

func createDBConnection(ctx context.Context, mysqlC testcontainers.Container) (*sql.DB, error) {
	host, err := mysqlC.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := mysqlC.MappedPort(ctx, dbPortNat)
	if err != nil {
		return nil, err
	}
	cfg := mysql.Config{
		DBName:    dbName,
		User:      "root",
		Addr:      fmt.Sprintf("%s:%d", host, port.Int()),
		Net:       "tcp",
		ParseTime: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}
