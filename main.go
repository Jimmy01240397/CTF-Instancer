package main
import (
//    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/go-errors/errors"

    "github.com/Jimmy01240397/CTF-Instancer/router"
    "github.com/Jimmy01240397/CTF-Instancer/middlewares/proxy"
    "github.com/Jimmy01240397/CTF-Instancer/middlewares/token"
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
    "github.com/Jimmy01240397/CTF-Instancer/utils/errutil"
)

func main() {
    if !config.Debug {
        gin.SetMode(gin.ReleaseMode)
    }
    backend := gin.Default()
    backend.Use(proxy.Proxy)
    backend.Use(errorHandler)
    backend.Use(gin.CustomRecovery(panicHandler))
    
    switch config.ServiceMode {
    case "web":
        store := cookie.NewStore(config.Secret)
        store.Options(sessions.Options{
            Path:     "/",
            MaxAge: 2592000,
            HttpOnly: true,
            Secure:   false,
        })
        backend.Use(sessions.Sessions(config.Sessionname, store))
        backend.LoadHTMLGlob("template/*")
    case "api":
        backend.Use(token.AddMeta)
    default:
        panic("bad mode")
    }
    
    router.Init(&backend.RouterGroup)
    backend.Run(":"+string(config.Port))
}

func panicHandler(c *gin.Context, err any) {
    goErr := errors.Wrap(err, 2)
    errmsg := ""
    if config.Debug {
        errmsg = goErr.Error()
    }
    errutil.AbortAndError(c, &errutil.Err{
        Code: 500,
        Msg: "Internal server error",
        Data: errmsg,
    })
}

func errorHandler(c *gin.Context) {
    c.Next()

    for _, e := range c.Errors {
        err := e.Err
        if myErr, ok := err.(*errutil.Err); ok {
            if myErr.Msg != nil {
                c.JSON(myErr.Code, myErr.ToH())
            } else {
                c.Status(myErr.Code)
            }
        } else {
            errmsg := ""
            if config.Debug {
                errmsg = err.Error()
            }
            c.JSON(500, gin.H{
                "code": 500,
                "msg": "Internal server error",
                "data": errmsg,
            })
        }
        return
    }
}
