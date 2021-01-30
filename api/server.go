package api

import (
	db "github.com/rondondev/runapp/db/sqlc"

	"github.com/gin-gonic/gin"

)

// Server serves HTTP requests
type Server struct {
	store db.Querier
	router *gin.Engine
}

type idRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Querier) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Users
	router.GET("/users", server.listUsers)
	router.GET("/users/active", server.listActiveUsers)
	router.GET("/users/all", server.listAllUsers)
	router.GET("/user/:id", server.getUser)
	router.POST("/user", server.createUser)
	router.PUT("/user/:id", server.updateUser)
	router.DELETE("/user/:id", server.deleteUser)

	// User trainings
	router.GET("/trainings/feedback/user/:id", server.listTrainingFeedbacksByUser)

	// Training
	router.GET("/training/:id/feedback", server.getTrainingFeedback)
	router.POST("/training/:id/feedback", server.createTrainingFeedback)
	router.PUT("/training/:id/feedback", server.updateTrainingFeedback)
	router.DELETE("/training/:id/feedback", server.deleteTrainingFeedback)

	// Training feedbacks
	router.GET("/trainings/user/:id", server.listTrainingsByUser)

	// Training
	router.GET("/training/:id", server.getTraining)
	router.POST("/training", server.createTraining)
	router.PUT("/training/:id", server.updateTraining)
	router.DELETE("/training/:id", server.deleteTraining)


	server.router = router

	return server
}

// Start runs the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}