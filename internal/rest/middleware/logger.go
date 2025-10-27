package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = io.ReadAll(c.Request.Body)
			if err == nil {
				log.Infof("Request Start => Method: %s| Path:  %s| Client IP:  %s| Request Body: %s",
					c.Request.Method, c.Request.URL.Path, c.ClientIP(), string(requestBody))
			} else {
				log.Errorf("Failed to read request body: %v", err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		log.Infof("Request End => Client IP:  %s|Code: %d| Latency: %s| Response Body: %s",
			c.ClientIP(), c.Writer.Status(), latency, blw.body.String())

	}
}
