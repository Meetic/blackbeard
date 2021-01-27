package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func jsonLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			fields := struct {
				Doctype      string  `json:"document_type"`
				Time         string  `json:"time"`
				Verb         string  `json:"verb"`
				Request      string  `json:"request"`
				User         string  `json:"user"`
				Httpversion  string  `json:"http_version"`
				Useragent    string  `json:"user_agent"`
				Remoteaddr   string  `json:"remoteaddr"`
				Status       int     `json:"status"`
				Responsetime float64 `json:"response_time"`
				Error        string  `json:"error"`
			}{
				Doctype:      "accesslog-blackbeard",
				Time:         param.TimeStamp.Format(time.RFC3339),
				Verb:         param.Method,
				Request:      param.Path,
				User:         param.Request.Header.Get("Remote-User"),
				Httpversion:  param.Request.Proto,
				Useragent:    param.Request.UserAgent(),
				Remoteaddr:   param.ClientIP,
				Status:       param.StatusCode,
				Responsetime: param.Latency.Seconds(),
				Error:        param.ErrorMessage,
			}

			jsonLog, _ := json.Marshal(fields)

			return fmt.Sprintf(string(jsonLog) + "\n")
		},
	})
}
