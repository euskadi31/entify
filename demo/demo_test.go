package demo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/euskadi31/entify/entify/entity"
	"github.com/stretchr/testify/assert"
)

func TestUserClientUpdateOne(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE `users` SET `email` = ? WHERE `id` = ?").
		WithArgs("user@email.tld", "fdgfgh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	u, err := c.User.UpdateOneID("fdgfgh").
		SetEmail("user@email.tld").
		Save(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user@email.tld", u.GetEmail())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserClientDeleteOne(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM `users` WHERE `id` = ?").
		WithArgs("ertyht").
		WillReturnResult(sqlmock.NewResult(0, 1))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	err = c.User.DeleteOneID("ertyht").Exec(ctx)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserClientQueryFindOne(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT * FROM users WHERE id = ?").
		WithArgs("fdgdgdgh").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).AddRow("fdgdgdgh", "user@email.tld", "passw0rd"))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	user, err := c.User.Query().FindOne(ctx, "SELECT * FROM users WHERE id = ?", "fdgdgdgh")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, "fdgdgdgh", user.GetID())
	assert.Equal(t, "user@email.tld", user.GetEmail())
	assert.Equal(t, "passw0rd", user.GetPassword())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserClientQueryFindAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT * FROM users WHERE id in(?, ?)").
		WithArgs("fdgdgdgh", "yrtyrtgr").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow("fdgdgdgh", "user@email.tld", "passw0rd").
				AddRow("yrtyrtgr", "user2@email.tld", "p1ssw0rd"),
		)

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	users, err := c.User.Query().FindAll(ctx, "SELECT * FROM users WHERE id in(?, ?)", "fdgdgdgh", "yrtyrtgr")
	assert.NoError(t, err)
	assert.NotNil(t, users)

	assert.Equal(t, 2, len(users))

	assert.Equal(t, "fdgdgdgh", users[0].GetID())
	assert.Equal(t, "user@email.tld", users[0].GetEmail())
	assert.Equal(t, "passw0rd", users[0].GetPassword())

	assert.Equal(t, "yrtyrtgr", users[1].GetID())
	assert.Equal(t, "user2@email.tld", users[1].GetEmail())
	assert.Equal(t, "p1ssw0rd", users[1].GetPassword())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO `users` (`id`, `email`, `password`) VALUES (?, ?, ?)").WithArgs("fdgfgh", "user@email.tld", "fdghfghgfh").WillReturnResult(sqlmock.NewResult(0, 1))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	u, err := c.User.Create().
		SetID("fdgfgh").
		SetEmail("user@email.tld").
		SetPassword("fdghfghgfh").
		Save(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user@email.tld", u.GetEmail())
	assert.Equal(t, "fdghfghgfh", u.GetPassword())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserUpdateOne(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO `users` (`id`, `email`, `password`) VALUES (?, ?, ?)").
		WithArgs("fdgfgh", "user@email.tld", "fdghfghgfh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("UPDATE `users` SET `email` = ? WHERE `id` = ?").
		WithArgs("user+test2@email.tld", "fdgfgh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	u, err := c.User.Create().
		SetID("fdgfgh").
		SetEmail("user@email.tld").
		SetPassword("fdghfghgfh").
		Save(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user@email.tld", u.GetEmail())
	assert.Equal(t, "fdghfghgfh", u.GetPassword())

	u, err = u.Update().SetEmail("user+test2@email.tld").Save(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user+test2@email.tld", u.GetEmail())
	assert.Equal(t, "fdghfghgfh", u.GetPassword())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserUpdateOneWithNullableEmail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO `users` (`id`, `email`, `firstname`, `password`) VALUES (?, ?, ?, ?)").
		WithArgs("fdgfgh", "user@email.tld", "foo", "fdghfghgfh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("UPDATE `users` SET `firstname` = ? WHERE `id` = ?").
		WithArgs(nil, "fdgfgh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	u, err := c.User.Create().
		SetID("fdgfgh").
		SetEmail("user@email.tld").
		SetFirstname("foo").
		SetPassword("fdghfghgfh").
		Save(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user@email.tld", u.GetEmail())
	assert.Equal(t, "foo", u.GetFirstname())
	assert.Equal(t, "fdghfghgfh", u.GetPassword())

	u, err = u.Update().ClearFirstname().Save(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, u)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user@email.tld", u.GetEmail())
	assert.Equal(t, "", u.GetFirstname())
	assert.Equal(t, "fdghfghgfh", u.GetPassword())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserDeleteOne(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO `users` (`id`, `email`, `password`) VALUES (?, ?, ?)").
		WithArgs("fdgfgh", "user@email.tld", "fdghfghgfh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM `users` WHERE `id` = ?").
		WithArgs("fdgfgh").
		WillReturnResult(sqlmock.NewResult(0, 1))

	c := entity.NewClient("mysql", db)

	ctx := context.Background()

	u, err := c.User.Create().
		SetID("fdgfgh").
		SetEmail("user@email.tld").
		SetPassword("fdghfghgfh").
		Save(ctx)

	assert.NoError(t, err)

	assert.Equal(t, "fdgfgh", u.GetID())
	assert.Equal(t, "user@email.tld", u.GetEmail())
	assert.Equal(t, "fdghfghgfh", u.GetPassword())

	err = u.Delete().Exec(ctx)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
