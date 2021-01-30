package db

import (
	"database/sql"
	"log"
	"testing"

	"github.com/rondondev/runapp/util"

	"github.com/jaswdr/faker"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

type DbTestSuite struct {
	suite.Suite
	q *Queries
	f faker.Faker
}

func (s *DbTestSuite) SetupSuite() {
	config, err := util.LoadConfig("../..", "test")
	if err != nil {
		log.Fatal("cannot read config file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// todo: run the tests in a fresh containerized DB
	// clean up test DB
	_, err = conn.Exec(`DELETE FROM training_feedback`)
	s.Require().NoError(err)
	_, err = conn.Exec(`DELETE FROM training`)
	s.Require().NoError(err)
	_, err = conn.Exec(`DELETE FROM users`)
	s.Require().NoError(err)

	s.q = New(conn)
	s.f = faker.New()
}
