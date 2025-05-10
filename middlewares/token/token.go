package token

import (
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/Jimmy01240397/CTF-Instancer/utils/errutil"
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

func CheckAuth(c *gin.Context) {
    if isAuth, exist := c.Get("isAuth"); !exist || !isAuth.(bool) {
        errutil.AbortAndStatus(c, 401)
    }
}

func AddMeta(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if fmt.Sprintf("Token %s", config.Token) == token {
        c.Set("isAuth", true)
    }
}
