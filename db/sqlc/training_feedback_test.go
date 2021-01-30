package db

import (
	"context"
	"github.com/emvi/null"
	"time"
)

func (s *DbTestSuite) createTrainingFeedback(trainingID int64, borgScale int32) TrainingFeedback {

	arg := CreateTrainingFeedbackParams{
		TrainingID: trainingID,
		BorgScale: borgScale,
	}

	feedback, err := s.q.CreateTrainingFeedback(context.Background(), arg)
	s.Require().NoError(err)
	s.NotEmpty(feedback)

	s.Equal(arg.TrainingID, feedback.TrainingID)
	s.Equal(arg.BorgScale, feedback.BorgScale)

	return feedback
}

func (s *DbTestSuite) TestCreateTrainingFeedback() {
	s.createUser(UserTypeAdmin, false)
}

func (s *DbTestSuite) TestDeleteTrainingFeedback() {
	u := s.createUser(UserTypeAthlete, true)
	t := s.createTraining(u.ID)
	s.createTrainingFeedback(t.ID, 10)

	err := s.q.DeleteTrainingFeedback(context.Background(), t.ID)
	s.Require().NoError(err)

	feedback, err := s.q.GetTrainingFeedback(context.Background(), t.ID)
	s.Require().Error(err)
	s.Empty(feedback)
}

func (s *DbTestSuite) TestGetTrainingFeedback() {
	u := s.createUser(UserTypeAthlete, true)
	t := s.createTraining(u.ID)
	f := s.createTrainingFeedback(t.ID, 10)

	feedback, err := s.q.GetTrainingFeedback(context.Background(), t.ID)
	s.Require().NoError(err)
	s.NotEmpty(feedback)

	s.Equal(f.ID, feedback.ID)
	s.Equal(f.TrainingID, feedback.TrainingID)
	s.Equal(f.BorgScale, feedback.BorgScale)
}

func (s *DbTestSuite) TestListTrainingFeedbacksByUser() {
	u1 := s.createUser(UserTypeAthlete, true)
	u2 := s.createUser(UserTypeAthlete, true)

	// 3 trainings, 3 feedbacks for user 1
	for i := 0; i< 3; i++ {
		t := s.createTraining(u1.ID)
		s.createTrainingFeedback(t.ID, 10)
	}

	// 3 trainings, 1 feedback for user 2
	var t Training
	for i := 0; i< 3; i++ {
		t = s.createTraining(u2.ID)
	}
	s.createTrainingFeedback(t.ID, 10)


	trainings, err := s.q.ListTrainingFeedbacksByUser(context.Background(), u1.ID)
	s.Require().NoError(err)
	s.NotEmpty(trainings)
	s.Len(trainings, 3)

	trainings, err = s.q.ListTrainingFeedbacksByUser(context.Background(), u2.ID)
	s.Require().NoError(err)
	s.NotEmpty(trainings)
	s.Len(trainings, 1)
}

func (s *DbTestSuite) TestListTrainingFeedbackByUserInPeriod() {
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
		training, err := s.q.CreateTraining(context.Background(), t)
		s.Require().NoError(err)
		f := CreateTrainingFeedbackParams{
			TrainingID: training.ID,
			BorgScale: 10,
		}
		_, err = s.q.CreateTrainingFeedback(context.Background(), f)
		s.Require().NoError(err)
	}

	startDate := date.AddDate(0, 0, 3)
	endDate := date.AddDate(0, 0, 7)

	arg := ListTrainingFeedbacksByUserInPeriodParams{
		u.ID,
		startDate,
		endDate,
	}

	trainings, err := s.q.ListTrainingFeedbacksByUserInPeriod(context.Background(), arg)
	s.Require().NoError(err)
	s.NotEmpty(trainings)
	s.Len(trainings, 5)
}


func (s *DbTestSuite) TestUpdateTrainingFeedback() {
	u := s.createUser(UserTypeAthlete, true)
	t := s.createTraining(u.ID)
	f := s.createTrainingFeedback(t.ID, 10)


	arg := UpdateTrainingFeedbackParams{
		TrainingID: t.ID,
		BorgScale: 15,
	}

	_, err := s.q.UpdateTrainingFeedback(context.Background(), arg)
	s.Require().NoError(err)

	feedback, err := s.q.GetTrainingFeedback(context.Background(), t.ID)
	s.Require().NoError(err)
	s.NotEmpty(feedback)

	s.Equal(f.ID, feedback.ID)
	s.Equal(arg.BorgScale, feedback.BorgScale)
}