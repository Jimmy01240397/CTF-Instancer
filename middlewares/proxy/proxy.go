package proxy

import (
    "fmt"
    "regexp"
    "net/http/httputil"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/Jimmy01240397/CTF-Instancer/models/instance"
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

func Proxy(c *gin.Context) {
    re := regexp.MustCompile(`:[0-9]+$`)
    for id, ins := range instance.GetIDMap() {
        if c.Request.Host == fmt.Sprintf("%s.%s%s", id, config.BaseHost, re.FindString(c.Request.Host)) {
            proxy := &httputil.ReverseProxy{}
            proxy.Director = func(req *http.Request) {
                req.Header = c.Request.Header
                req.Host = c.Request.Host
                req.URL.Scheme = "http"
                req.URL.Host = fmt.Sprintf("localhost:%d", ins.Port)
                req.URL.Path = c.Request.URL.Path
            }
            proxy.ModifyResponse = func(resp *http.Response) error {
                return nil
            }
            proxy.ServeHTTP(c.Writer, c.Request)
            c.Abort()
        }
    }
}
