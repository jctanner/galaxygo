package utils

import (
    "github.com/gin-gonic/gin"
)

func GetRequestHostFromContext(c *gin.Context) (string, error) {

    // get the request host ...
    scheme := "http"
    if c.Request.TLS != nil {
        scheme = "https"
    }
    rhost := scheme + "://" + c.Request.Host

    return rhost, nil
}
