package adapters

import (
	"errors"
	"fmt"
	"github.com/jitsucom/jitsu/server/logging"
	"github.com/jitsucom/jitsu/server/typing"
	"net/http"
	"net/url"
	"time"
)

const defaultEventType = "event"

var (
	//GA doesn't use types
	SchemaToGoogleAnalytics = map[typing.DataType]string{
		typing.STRING:    "string",
		typing.INT64:     "string",
		typing.FLOAT64:   "string",
		typing.TIMESTAMP: "string",
		typing.BOOL:      "string",
		typing.UNKNOWN:   "string",
	}

	gaEventTypeMapping = map[string]string{
		"pageview":    "pageview",
		"screenview":  "screenview",
		"event":       "event",
		"conversion":  "transaction",
		"transaction": "transaction",
		"item":        "item",
		"social":      "social",
		"exception":   "exception",
		"timing":      "timing",
	}
)

type GoogleAnalyticsConfig struct {
	TrackingID string `mapstructure:"tracking_id" json:"tracking_id,omitempty" yaml:"tracking_id,omitempty"`
}

func (gac *GoogleAnalyticsConfig) Validate() error {
	if gac == nil {
		return errors.New("google_analytics config is required")
	}
	if gac.TrackingID == "" {
		return errors.New("tracking_id is required parameter")
	}

	return nil
}

type GoogleAnalytics struct {
	config      *GoogleAnalyticsConfig
	client      *http.Client
	debugLogger *logging.QueryLogger
}

func NewGoogleAnalytics(config *GoogleAnalyticsConfig, requestDebugLogger *logging.QueryLogger) *GoogleAnalytics {
	return &GoogleAnalytics{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: 1000,
			},
		},
		debugLogger: requestDebugLogger,
	}
}

//Send HTTP GET request to GoogleAnalytics with query parameters
//remove system fields and map event type
func (ga *GoogleAnalytics) Send(object map[string]interface{}) error {
	uv := make(url.Values)
	uv.Add("tid", ga.config.TrackingID)
	uv.Add("v", "1")

	for k, v := range object {
		strValue, ok := v.(string)
		if !ok {
			strValue = fmt.Sprint(v)
		}

		//override event type
		if k == "t" {
			mapped, ok := gaEventTypeMapping[strValue]
			if !ok {
				mapped = defaultEventType
			}

			strValue = mapped
		}

		uv.Add(k, strValue)
	}

	reqURL := "https://www.google-analytics.com/collect?" + uv.Encode()
	ga.debugLogger.LogQuery(reqURL)

	r, err := ga.client.Get(reqURL)
	if r != nil && r.Body != nil {
		r.Body.Close()
	}

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("Google Analytics response code: %d", r.StatusCode)
	}

	return nil
}

//GetTableSchema always return empty schema
func (ga *GoogleAnalytics) GetTableSchema(tableName string) (*Table, error) {
	return &Table{
		Name:           tableName,
		Columns:        Columns{},
		PKFields:       map[string]bool{},
		DeletePkFields: false,
		Version:        0,
	}, nil
}

//CreateTable GA doesn't use tables
func (ga *GoogleAnalytics) CreateTable(schemaToCreate *Table) error {
	return nil
}

//PatchTableSchema GA doesn't use tables
func (ga *GoogleAnalytics) PatchTableSchema(schemaToAdd *Table) error {
	return nil
}

func (ga *GoogleAnalytics) Close() error {
	ga.client.CloseIdleConnections()

	return nil
}
