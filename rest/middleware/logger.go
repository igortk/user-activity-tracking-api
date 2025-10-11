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

		log.Infof("--- Request Start ---")
		log.Infof("Method: %s", c.Request.Method)
		log.Infof("Path:  %s", c.Request.URL.Path)
		log.Infof("Client IP:  %s", c.ClientIP())

		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = io.ReadAll(c.Request.Body)
			if err == nil {
				log.Infof("Request Body: %s", string(requestBody))
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

		log.Infof("Status Code: %d", c.Writer.Status())
		log.Infof("Response Latency: %s", latency)
		log.Infof("Response Body: %s", blw.body.String())
		log.Infof("--- Request End ---")
	}
}
