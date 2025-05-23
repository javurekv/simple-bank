package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	db "simple_bank.sqlc.dev/app/db/sqlc"
	"simple_bank.sqlc.dev/app/util"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
