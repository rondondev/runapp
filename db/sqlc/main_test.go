package db

import (
	"database/sql"
	"github.com/jaswdr/faker"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:6432/runapp?sslmode=disable&timezone=UTC"
)

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

type DbTestSuite struct {
	suite.Suite
	q	 *Queries
	f faker.Faker
}

func (s *DbTestSuite) SetupSuite() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// todo: run the tests in a fresh containerized DB
	// clean DB
	_, err = conn.Exec(`DELETE FROM training_feedback`)
	s.Require().NoError(err)
	_, err = conn.Exec(`DELETE FROM training`)
	s.Require().NoError(err)
	_, err = conn.Exec(`DELETE FROM users`)
	s.Require().NoError(err)

	s.q = New(conn)
	s.f = faker.New()
}
