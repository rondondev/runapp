package db

import (
	"context"
	"github.com/emvi/null"
	"time"
)

func (s *DbTestSuite) createTraining(userID int64) Training {

	arg := CreateTrainingParams{
		UserID: userID,
		Date: time.Now().UTC(),
		Sport: TrainingSportRunning,
		Type: null.String{},
		Intensity: null.String{},
		Details: "random training",
		Status: TrainingStatusNew,
	}

	training, err := s.q.CreateTraining(context.Background(), arg)
	s.Require().NoError(err)
	s.NotEmpty(training)

	s.Equal(arg.Type, training.Type)
	s.Equal(arg.UserID, training.UserID)
	s.Equal(arg.Date.Format("2006-01-02"), training.Date.Format("2006-01-02"))
	s.Equal(arg.Sport, training.Sport)
	s.Equal(arg.Type, training.Type)
	s.Equal(arg.Intensity, training.Intensity)
	s.Equal(arg.Details, training.Details)
	s.Equal(arg.Status, training.Status)

	s.NotEmpty(training.CreatedAt)
	s.Empty(training.DeletedAt)

	return training
}

func (s *DbTestSuite) TestCreateTraining() {
	s.createUser(UserTypeAdmin, false)
}

func (s *DbTestSuite) TestDeleteTraining() {
	u := s.createUser(UserTypeAthlete, true)
	t := s.createTraining(u.ID)

	err := s.q.DeleteTraining(context.Background(), t.ID)
	s.Require().NoError(err)

	training, err := s.q.GetTraining(context.Background(), t.ID)
	s.Require().Error(err)
	s.Empty(training)
}

func (s *DbTestSuite) TestGetTraining() {
	u := s.createUser(UserTypeAthlete, true)
	t := s.createTraining(u.ID)

	training, err := s.q.GetTraining(context.Background(), t.ID)
	s.Require().NoError(err)
	s.NotEmpty(training)

	s.Equal(t.ID, training.ID)
	s.Equal(t.Type, training.Type)
	s.Equal(t.UserID, training.UserID)
	s.Equal(t.Date, training.Date)
	s.Equal(t.Sport, training.Sport)
	s.Equal(t.Type, training.Type)
	s.Equal(t.Intensity, training.Intensity)
	s.Equal(t.Details, training.Details)
	s.Equal(t.Status, training.Status)
	s.Equal(t.CreatedAt, training.CreatedAt)
}

func (s *DbTestSuite) TestListTrainingByUser() {
	u1 := s.createUser(UserTypeAthlete, true)
	u2 := s.createUser(UserTypeAthlete, true)

	for i := 0; i< 3; i++ {
		s.createTraining(u1.ID)
	}
	for i := 0; i< 4; i++ {
		s.createTraining(u2.ID)
	}

	trainings, err := s.q.ListTrainingsByUser(context.Background(), u1.ID)
	s.Require().NoError(err)
	s.NotEmpty(trainings)
	s.Len(trainings, 3)
	for _, t := range trainings {
		s.Equal(t.UserID, u1.ID)
	}

	trainings, err = s.q.ListTrainingsByUser(context.Background(), u2.ID)
	s.Require().NoError(err)
	s.NotEmpty(trainings)
	s.Len(trainings, 4)
	for _, t := range trainings {
		s.Equal(t.UserID, u2.ID)
	}
}

func (s *DbTestSuite) TestListTrainingsByUserInPeriod() {
	u := s.createUser(UserTypeAthlete, true)

	t := CreateTrainingParams{
		UserID: u.ID,
		Sport: TrainingSportRunning,
		Type: null.String{},
		Intensity: null.String{},
		Details: "random training",
		Status: TrainingStatusNew,
	}

	date := time.Now()
	for i := 0; i<10; i++ {
		t.Date = date.AddDate(0, 0, i)
		_, err := s.q.CreateTraining(context.Background(), t)
		s.Require().NoError(err)
	}

	startDate := date.AddDate(0, 0, 3)
	endDate := date.AddDate(0, 0, 7)

	arg := ListTrainingsByUserInPeriodParams{
		u.ID,
		startDate,
		endDate,
	}

	trainings, err := s.q.ListTrainingsByUserInPeriod(context.Background(), arg)
	s.Require().NoError(err)
	s.NotEmpty(trainings)
	s.Len(trainings, 5)
	for _, t := range trainings {
		s.Equal(t.UserID, u.ID)
	}
}


func (s *DbTestSuite) TestUpdateTraining() {
	u := s.createUser(UserTypeAthlete, true)
	t := s.createTraining(u.ID)

	d, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	arg := UpdateTrainingParams{
		ID: t.ID,
		Date: d,
		Sport: TrainingSportCycling,
		Type: null.NewString("interval", true),
		Intensity: null.NewString("high", true),
		Details: "details 2",
		Status: TrainingStatusNotified,
	}

	_, err := s.q.UpdateTraining(context.Background(), arg)
	s.Require().NoError(err)

	training, err := s.q.GetTraining(context.Background(), arg.ID)
	s.Require().NoError(err)
	s.NotEmpty(training)

	s.Equal(t.ID, training.ID)
	s.Equal(t.UserID, training.UserID)
	s.Equal(arg.Date.Format("2006-01-02"), training.Date.Format("2006-01-02"))
	s.Equal(arg.Sport, training.Sport)
	s.Equal(arg.Type, training.Type)
	s.Equal(arg.Intensity, training.Intensity)
	s.Equal(arg.Details, training.Details)
	s.Equal(arg.Status, training.Status)

	s.Equal(t.CreatedAt, training.CreatedAt)
	s.Equal(t.DeletedAt, training.DeletedAt)
}