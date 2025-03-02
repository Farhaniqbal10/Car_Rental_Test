package driver

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(r *gin.Engine, db *sqlx.DB) error {
	initCar(r, db)

	return nil
}
