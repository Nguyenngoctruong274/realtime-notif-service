package handler

import (
	"github.com/gin-gonic/gin"
	"yes4all/ads-noti-api/pkg/utils/common"
)

type HealthCheck interface {
	HealthCheckInfo(c *gin.Context)
}
type healthCheck struct {
}

func NewHealthCheck() HealthCheck {
	return &healthCheck{}
}

// HealthCheckInfo health_check godoc
// @Summary HealthCheckInfo of health_check
// @Description HealthCheckInfo of health_check
// @Tags health_check
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /external/health-check/info [get]
func (m *healthCheck) HealthCheckInfo(c *gin.Context) {
	common.JSONOk(c, nil)
}
