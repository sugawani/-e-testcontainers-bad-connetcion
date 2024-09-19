package query

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"main/models"
	"main/util"
)

func Test_Query(t *testing.T) {
	createUser := func(db *sql.DB) {
		db.Exec("INSERT INTO users (id, name) VALUES (1, 'name')")
	}
	noCreateUser := func(db *sql.DB) {}
	wantErrAssertFunc := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, sql.ErrNoRows)
	}
	cases := map[string]struct {
		createFunc func(db *sql.DB)
		want       *models.User
		assertErr  assert.ErrorAssertionFunc
	}{
		"user exists":      {createFunc: createUser, want: &models.User{ID: 1, Name: "name"}, assertErr: assert.NoError},
		"user exists2":     {createFunc: createUser, want: &models.User{ID: 1, Name: "name"}, assertErr: assert.NoError},
		"user exists3":     {createFunc: createUser, want: &models.User{ID: 1, Name: "name"}, assertErr: assert.NoError},
		"user not exists":  {createFunc: noCreateUser, want: nil, assertErr: wantErrAssertFunc},
		"user not exists2": {createFunc: noCreateUser, want: nil, assertErr: wantErrAssertFunc},
		"user not exists3": {createFunc: noCreateUser, want: nil, assertErr: wantErrAssertFunc},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			db, err := util.NewTestDB(ctx)
			if err != nil {
				t.Fatal("failed to create test db", err)
			}

			tt.createFunc(db)
			q := NewQuery(db)
			actual, err := q.Execute(1)
			tt.assertErr(t, err)
			assert.Equal(t, tt.want, actual)
		})
	}
}
