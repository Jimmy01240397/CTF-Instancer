package router
import (
    "fmt"
    "time"
    "regexp"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
    "github.com/Jimmy01240397/CTF-Instancer/utils/hcaptcha"
    "github.com/Jimmy01240397/CTF-Instancer/models/auth"
    "github.com/Jimmy01240397/CTF-Instancer/models/instance"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.GET("/", index)
    router.POST("/create", create)
    router.POST("/stop", stop)
    //user.Init(router.Group("/user"))
}

func index(c *gin.Context) {
    session := sessions.Default(c)
    name := ""
    name, _ = session.Get("name").(string)
    ins := instance.GetInstance(name)
    data := gin.H{
        "Title": config.Title,
        "HCAPTCHA_SITE_KEY": config.HCAPTCHA_SITE_KEY,
        "Now": time.Now(),
    }
    if ins == nil {
        data["InstanceId"] = ""
    } else {
        data["InstanceId"] = ins.ID
        data["ExpiredAt"] = ins.ExpiredAt
        if config.ProxyMode {
            re := regexp.MustCompile(`:[0-9]+$`)
            data["URL"] = fmt.Sprintf("%s://%s.%s%s", config.BaseScheme, ins.ID, config.BaseHost, re.FindString(c.Request.Host))
        } else {
            data["URL"] = fmt.Sprintf("%s://%s:%d", config.BaseScheme, config.BaseHost, ins.Port)
        }
    }
    c.HTML(200, "index.tmpl", data)
}

func create(c *gin.Context) {
    token := c.PostForm("token")
    hcaptchares := c.PostForm("h-captcha-response")
    if !hcaptcha.Verify(hcaptchares) {
        c.String(400, "Captcha verification failed.")
        return
    }
    name, err := auth.Auth(token)
    if err != nil {
        c.String(400, "Invalid Token.")
        return
    }
    _, err = instance.Up(name)
    if err != nil {
        panic(err)
    }
    session := sessions.Default(c)
    session.Set("name", name)
    session.Save()
    c.Redirect(301, "/")
}

func stop(c *gin.Context) {
    session := sessions.Default(c)
    name := ""
    name, _ = session.Get("name").(string)
    err := instance.Down(name)
    if err != nil {
        panic(err)
    }
    session.Clear()
    session.Save()
    c.Redirect(301, "/")
}
