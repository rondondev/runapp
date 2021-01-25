package db

import (
	"context"
	"math"
	"time"

	"github.com/emvi/null"
)

func (s *DbTestSuite) createUser(userType UserType, active bool) User {
	arg := CreateUserParams{
		Type: userType,
		Name: s.f.Person().Name(),
		Email: s.f.Internet().Email(),
		PasswordHash: "hash",
		Phone: null.NewString("12345678", true),
		Birth: null.NewTime(time.Now().UTC(), true),
		Active: active,
	}

	user, err := s.q.CreateUser(context.Background(), arg)
	s.Require().NoError(err)
	s.NotEmpty(user)

	s.Equal(arg.Type, user.Type)
	s.Equal(arg.Name, user.Name)
	s.Equal(arg.Email, user.Email)
	s.Equal(arg.PasswordHash, user.PasswordHash)
	s.Equal(arg.Phone, user.Phone)
	s.Equal(arg.Birth.Time.Format("2006-01-02"), user.Birth.Time.Format("2006-01-02"))
	s.Equal(arg.Active, user.Active)

	s.NotEmpty(user.CreatedAt)
	s.Empty(user.DeletedAt)

	return user
}

func (s *DbTestSuite) TestCreateUser() {
	s.createUser(UserTypeAdmin, false)
}

func (s *DbTestSuite) TestDeleteUser() {
	u := s.createUser(UserTypeAdmin, false)

	err := s.q.DeleteUser(context.Background(),u.ID)
	s.Require().NoError(err)

	_, err = s.q.GetUser(context.Background(), u.ID)
	s.Require().Error(err)
}

func (s *DbTestSuite) TestGetUser() {
	u := s.createUser(UserTypeAdmin, false)

	user, err := s.q.GetUser(context.Background(), u.ID)
	s.Require().NoError(err)
	s.NotEmpty(user)

	s.Equal(u.ID, user.ID)
	s.Equal(u.Type, user.Type)
	s.Equal(u.Name, user.Name)
	s.Equal(u.Email, user.Email)
	s.Equal(u.PasswordHash, user.PasswordHash)
	s.Equal(u.Phone, user.Phone)
	s.Equal(u.Birth.Time, user.Birth.Time)
	s.Equal(u.Active, user.Active)
	s.Equal(u.CreatedAt, user.CreatedAt)
}

func (s *DbTestSuite) TestListActiveUsers() {
	arg := ListActiveUsersParams{
		Limit: math.MaxInt32,
		Offset: 0,
	}

	// Get the current number of active users
	users, err := s.q.ListActiveUsers(context.Background(), arg)
	s.Require().NoError(err)
	initialUsers := len(users)

	// Create a few inactive users
	for i := 0; i<5; i++ {
		s.createUser(UserTypeAthlete, false)
	}

	// Get the number of active users - shouldn't have changed
	users, err = s.q.ListActiveUsers(context.Background(), arg)
	s.Require().NoError(err)
	s.Require().Len(users, initialUsers)

	// Create five active users
	for i := 0; i<5; i++ {
		s.createUser(UserTypeAthlete, true)
	}

	// Get the number of active users - shouldn't have increased by 5
	users, err = s.q.ListActiveUsers(context.Background(), arg)
	s.Require().NoError(err)
	s.Require().Len(users, initialUsers+5)
}

func (s *DbTestSuite) TestListAllUsers() {
	arg := ListAllUsersParams{
		Limit: math.MaxInt32,
		Offset: 0,
	}

	// Get the list of users
	users, err := s.q.ListAllUsers(context.Background(), arg)
	s.Require().NoError(err)
	initialUsers := len(users)

	// Create a few inactive users
	for i := 0; i<5; i++ {
		s.createUser(UserTypeAthlete, false)
	}
	// Create a few active users
	for i := 0; i<5; i++ {
		s.createUser(UserTypeAthlete, false)
	}
	// Create a few deleted users
	for i := 0; i<5; i++ {
		u := s.createUser(UserTypeAthlete, false)
		err := s.q.DeleteUser(context.Background(), u.ID)
		s.NoError(err)
	}

	// // Get the list of users - shouldn't have increased by 15
	users, err = s.q.ListAllUsers(context.Background(), arg)
	s.Require().NoError(err)
	s.Require().Len(users, initialUsers+15)
}

func (s *DbTestSuite) TestListUsers() {
	arg := ListUsersParams{
		Limit: math.MaxInt32,
		Offset: 0,
	}

	// Get the list of users
	users, err := s.q.ListUsers(context.Background(), arg)
	s.Require().NoError(err)
	initialUsers := len(users)

	// Create a few inactive users
	for i := 0; i<5; i++ {
		s.createUser(UserTypeAthlete, false)
	}
	// Create a few active users
	for i := 0; i<5; i++ {
		s.createUser(UserTypeAthlete, false)
	}
	// Create a few deleted users
	for i := 0; i<5; i++ {
		u := s.createUser(UserTypeAthlete, false)
		err := s.q.DeleteUser(context.Background(), u.ID)
		s.NoError(err)
	}

	// // Get the list of users - shouldn't have increased by 10
	users, err = s.q.ListUsers(context.Background(), arg)
	s.Require().NoError(err)
	s.Require().Len(users, initialUsers+10)
}

func (s *DbTestSuite) TestUpdateUser() {
	u := s.createUser(UserTypeAdmin, false)

	arg := UpdateUserParams{
		ID: u.ID,
		Type: u.Type,
		Name: s.f.Person().Name(),
		Email: s.f.Internet().Email(),
		PasswordHash: u.PasswordHash,
		Phone: u.Phone,
		Birth: u.Birth,
		Active: true,
	}

	user, err := s.q.UpdateUser(context.Background(), arg)
	s.Require().NoError(err)
	s.NotEmpty(user)

	s.Equal(u.ID, user.ID)
	s.Equal(u.Type, user.Type)
	s.Equal(arg.Name, user.Name)
	s.Equal(arg.Email, user.Email)
	s.Equal(u.PasswordHash, user.PasswordHash)
	s.Equal(u.Phone, user.Phone)
	s.Equal(u.Birth.Time, user.Birth.Time)
	s.Equal(true, user.Active)
	s.Equal(u.CreatedAt, user.CreatedAt)
}
