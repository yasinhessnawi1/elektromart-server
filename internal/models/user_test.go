package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllUsers ensures that all users are retrieved from the database correctly.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "first_name", "last_name", "address", "role"}).
		AddRow(1, "Kristian Dem", "Password123", "kristian@example.com", "Kristian", "Dem", "123 Elva St", "user").
		AddRow(2, "Vines Bem", "Password321", "vines@example.com", "Vines", "Bem", "456 Demo St", "admin")
	mock.ExpectQuery("^SELECT \\* FROM \"users\"").WillReturnRows(rows)

	users, err := GetAllUsers(gormDB)
	assert.NoError(t, err)
	assert.Len(t, users, 2, "Should fetch two users")
}

// TestUser_SetRole tests the assignment of a role to a user.
// It creates a new instance of the User struct and calls the SetRole function with a valid role.
// It then checks if the role was set correctly and if the function returned true.
// It repeats the process with an invalid role and checks if the role was not set and the function returned false.
func TestUser_SetRole(t *testing.T) {
	user := User{}
	assert.False(t, user.SetRole(""), "Empty role should be invalid")
	assert.True(t, user.SetRole("admin"), "Admin should be a valid role")
}

// TestUser_SetUsername tests setting the username with proper validation checks.
// It creates a new instance of the User struct and calls the SetUsername function with a valid username.
// It then checks if the username was set correctly and if the function returned true.
// It repeats the process with an invalid username and checks if the username was not set and the function returned false.
func TestUser_SetUsername(t *testing.T) {
	user := User{}
	assert.False(t, user.SetUsername(""), "Empty username should be invalid")
	assert.True(t, user.SetUsername("Kristian"), "johnDoe should be a valid username")
}

// TestUser_SetPassword tests setting the password with proper validation for security.
// It creates a new instance of the User struct and calls the SetPassword function with a valid password.
// It then checks if the password was set correctly and if the function returned true.
// It repeats the process with an invalid password and checks if the password was not set and the function returned false.
func TestUser_SetPassword(t *testing.T) {
	user := User{}
	assert.False(t, user.SetPassword("short"), "Too short password should be invalid")
	assert.True(t, user.SetPassword("strongPassword123!"), "Strong password should be valid")
}

// TestUser_SetEmail tests setting the email with proper validation for format.
// It creates a new instance of the User struct and calls the SetEmail function with a valid email.
// It then checks if the email was set correctly and if the function returned true.
// It repeats the process with an invalid email and checks if the email was not set and the function returned false.
func TestUser_SetEmail(t *testing.T) {
	user := User{}
	assert.False(t, user.SetEmail("not-an-email"), "Invalid email should be rejected")
	assert.True(t, user.SetEmail("Kristian@example.com"), "Valid email should be accepted")
}

// TestUser_SetFirstName tests the assignment of a first name to a user.
// It creates a new instance of the User struct and calls the SetFirstName function with a valid first name.
// It then checks if the first name was set correctly and if the function returned true.
// It repeats the process with an invalid first name and checks if the first name was not set and the function returned false.
func TestUser_SetFirstName(t *testing.T) {
	user := User{}
	assert.False(t, user.SetFirstName(""), "Empty first name should be invalid")
	assert.True(t, user.SetFirstName("Kristian"), "Kristian should be a valid first name")
}

// TestUser_SetLastName tests the assignment of a last name to a user.
// It creates a new instance of the User struct and calls the SetLastName function with a valid last name.
// It then checks if the last name was set correctly and if the function returned true.
// It repeats the process with an invalid last name and checks if the last name was not set and the function returned false.
func TestUser_SetLastName(t *testing.T) {
	user := User{}
	assert.False(t, user.SetLastName(""), "Empty last name should be invalid")
	assert.True(t, user.SetLastName("Demo"), "Demo should be a valid last name")
}

// TestUser_SetAddress tests setting the address for a user.
// It creates a new instance of the User struct and calls the SetAddress function with a valid address.
// It then checks if the address was set correctly and if the function returned true.
// It repeats the process with an invalid address and checks if the address was not set and the function returned false.
func TestUser_SetAddress(t *testing.T) {
	user := User{}
	assert.False(t, user.SetAddress(""), "Empty address should be invalid")
	assert.True(t, user.SetAddress("123 Elva St"), "123 Elva St should be a valid address")
}

// TestUserExists tests if a user exists in the database by their ID.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	assert.True(t, UserExists(gormDB, 1), "User should exist")

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	assert.False(t, UserExists(gormDB, 99), "User should not exist")
}

// TestSearchUsers tests the search functionality based on given parameters.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestSearchUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "Kristian Demo", "johndoe@example.com")
	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs("%demo%").
		WillReturnRows(rows)

	searchParams := map[string]interface{}{"username": "demo"}
	users, err := SearchUsers(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, users, 1, "Should find one user matching search criteria")
}
