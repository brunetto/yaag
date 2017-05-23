package gin

import (

	"github.com/brunetto/yaag/middleware"
	"github.com/brunetto/yaag/yaag"
	"github.com/brunetto/yaag/yaag/models"
	"gopkg.in/gin-gonic/gin.v1"
	"strings"
	"log"
	"net/http"
	"bytes"
)

// Ref: https://stackoverflow.com/a/38548555/1283701
type responseCloneWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseCloneWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Document() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !yaag.IsOn() {
			return
		}
		apiCall := models.ApiCall{}
		middleware.Before(&apiCall, c.Request)

		blw := &responseCloneWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode != 404 {
			apiCall.MethodType = c.Request.Method
			apiCall.CurrentPath = strings.Split(c.Request.RequestURI, "?")[0]
			apiCall.ResponseBody = "omitted, non JSON"
			apiCall.ResponseCode = c.Writer.Status()
			apiCall.ResponseCodeString = http.StatusText(apiCall.ResponseCode)
			headers := map[string]string{}
			for k, v := range c.Writer.Header() {
				if yaag.Debug {
					log.Println(k, v)
				}
				headers[k] = strings.Join(v, " ")

				// Write response body if JSON
				if k == "Content-Type" {
					for _, h := range v {
						if strings.Contains(h, "application/json") {
							apiCall.ResponseBody = blw.body.String()
						}
					}
				}
			}
			apiCall.ResponseHeader = headers
			go yaag.GenerateHtml(&apiCall)
		}
	}
}
