package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jitsucom/jitsu/configurator/authorization"
	"github.com/jitsucom/jitsu/configurator/storages"
	enmiddleware "github.com/jitsucom/jitsu/server/middleware"
	"net/http"
)

//Configuration is used for system endpoint
//current authorization configuration and amount of registered users
type Configuration struct {
	Authorization          string `json:"authorization"`
	Users                  bool   `json:"users"`
	SMTP                   bool   `json:"smtp"`
	SelfHosted             bool   `json:"selfhosted"`
	SupportWidget          bool   `json:"support_widget"`
	DefaultS3Bucket        bool   `json:"default_s3_bucket"`
	SupportTrackingDomains bool   `json:"support_tracking_domains"`
	TelemetryUsageDisabled bool   `json:"telemetry_usage_disabled"`
}

type SystemHandler struct {
	authService          *authorization.Service
	configurationService *storages.ConfigurationsService
	smtp                 bool
	selfHosted           bool
}

type ConfigurationResponse struct {
	ProjectID string `json:"projectID"`
}

func NewSystemHandler(authService *authorization.Service, configurationService *storages.ConfigurationsService, smtp, selfHosted bool) *SystemHandler {
	return &SystemHandler{authService: authService, configurationService: configurationService, smtp: smtp, selfHosted: selfHosted}
}

//GetHandler returns JSON with current authorization type and users existence
func (sh *SystemHandler) GetHandler(c *gin.Context) {
	exist, err := sh.authService.UsersExist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, enmiddleware.ErrorResponse{
			Message: "Error checking users existence",
			Error:   err.Error(),
		})
		return
	}

	telemetryConfig, err := sh.configurationService.GetParsedTelemetry()
	if err != nil && err != storages.ErrConfigurationNotFound {
		errorMessage := fmt.Sprintf("Error getting telemetry configuration: %v", err)
		c.JSON(http.StatusInternalServerError, enmiddleware.ErrorResponse{Error: errorMessage, Message: "Telemetry error"})
		return
	}

	var telemetryUsageDisabled bool
	if telemetryConfig != nil && telemetryConfig.Disabled != nil {
		usageDisabled, ok := telemetryConfig.Disabled[telemetryUsageKey]
		if ok {
			telemetryUsageDisabled = usageDisabled
		}
	}

	currentConfiguration := Configuration{
		Authorization:          sh.authService.GetAuthorizationType(),
		Users:                  exist,
		SMTP:                   sh.smtp,
		SelfHosted:             sh.selfHosted,
		SupportWidget:          !sh.selfHosted,
		DefaultS3Bucket:        !sh.selfHosted,
		SupportTrackingDomains: !sh.selfHosted,
		TelemetryUsageDisabled: telemetryUsageDisabled,
	}

	c.JSON(http.StatusOK, currentConfiguration)
}
