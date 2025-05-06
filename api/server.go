package api

import (
	"github.com/gin-gonic/gin"
	db "simple_bank.sqlc.dev/app/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// New server
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts")

	server.router = router
	return server
}
