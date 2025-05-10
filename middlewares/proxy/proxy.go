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
        for idx, port := range ins.Ports {
            if config.GetMode(idx) == config.Proxy && c.Request.Host == fmt.Sprintf("%s%d.%s%s", id, idx, config.BaseHost, re.FindString(c.Request.Host)) {
                proxy := &httputil.ReverseProxy{}
                proxy.Director = func(req *http.Request) {
                    req.Header = c.Request.Header
                    req.Host = c.Request.Host
                    req.URL.Scheme = "http"
                    req.URL.Host = fmt.Sprintf("localhost:%d", port)
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
}
