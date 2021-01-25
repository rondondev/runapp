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