package router
import (
    "fmt"
    "bytes"
    "time"
    "regexp"
    "strconv"
    "text/template"
    htemplate "html/template"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    "github.com/Jimmy01240397/CTF-Instancer/middlewares/token"
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
    "github.com/Jimmy01240397/CTF-Instancer/utils/captcha"
    "github.com/Jimmy01240397/CTF-Instancer/utils/errutil"
    "github.com/Jimmy01240397/CTF-Instancer/models/auth"
    "github.com/Jimmy01240397/CTF-Instancer/models/instance"
)

type statusdata struct {
    AccessPoint htemplate.HTML `json:"accesspoint"`
    ExpiredAt time.Time `json:"expiredat"`
}

type userdata struct {
    ID int `json:"userid"`
}

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    switch config.ServiceMode {
    case "web":
        router.GET("/", index)
        router.POST("/create", create)
        router.POST("/destroy", destroy)
    case "api":
        router.GET("/", token.CheckAuth, index)
        router.GET("/flag", token.CheckAuth, flag)
        router.POST("/create", token.CheckAuth, create)
        router.POST("/destroy", token.CheckAuth, destroy)
    default:
        panic("bad mode")
    }
}

func index(c *gin.Context) {
    var name string
    var session sessions.Session
    data := gin.H{
        "Title": config.Title,
        "CAPTCHA_SRC": config.CAPTCHA_SRC,
        "CAPTCHA_CLASS": config.CAPTCHA_CLASS,
        "CAPTCHA_SITE_KEY": config.CAPTCHA_SITE_KEY,
        "Now": time.Now(),
    }
    switch config.ServiceMode {
    case "web":
        session = sessions.Default(c)
        name, _ = session.Get("name").(string)
        ins := instance.GetInstance(name)
        if ins == nil {
            data["InstanceId"] = ""
            c.HTML(200, "index.tmpl", data)
            return
        }
    case "api":
        name = c.Query("userid")
        if name == "" {
            errutil.AbortAndStatus(c, 404)
            return
        }
        ins := instance.GetInstance(name)
        if ins == nil {
            errutil.AbortAndStatus(c, 404)
            return
        }
    default:
        errutil.AbortAndStatus(c, 500)
        return
    }
    ins := instance.GetInstance(name)
    status := statusdata{
        ExpiredAt: ins.ExpiredAt,
        AccessPoint: "",
    }
    re := regexp.MustCompile(`:[0-9]+$`)
    for i, port := range ins.Ports {
        switch config.GetMode(i) {
        case config.Forward:
            tmp := fmt.Sprintf("%s://%s:%d", config.BaseScheme, config.BaseHost, port)
            status.AccessPoint += htemplate.HTML(fmt.Sprintf("<a href=\"%s\">%s</a><br/>", tmp, tmp))
        case config.Proxy:
            tmp := fmt.Sprintf("%s://%s%d.%s%s", config.BaseScheme, ins.ID, i, config.BaseHost, re.FindString(c.Request.Host))
            status.AccessPoint += htemplate.HTML(fmt.Sprintf("<a href=\"%s\">%s</a><br/>", tmp, tmp))
        case config.Command:
            tmp, err := template.New("command").Parse(config.GetCommand(i))
            if err != nil {
                errutil.AbortAndStatus(c, 500)
                return
            }
            var buf bytes.Buffer
            err = tmp.Execute(&buf, struct {
                BaseHost string
                Port uint16
            } {
                BaseHost: config.BaseHost,
                Port: port,
            })
            if err != nil {
                errutil.AbortAndStatus(c, 500)
                return
            }
            status.AccessPoint += htemplate.HTML(fmt.Sprintf("<code>%s</code><br/>", buf.String()))
        }
    }
    switch config.ServiceMode {
    case "web":
        data["InstanceId"] = ins.ID
        data["Status"] = status
        c.HTML(200, "index.tmpl", data)
        return
    case "api":
        c.JSON(200, status)
        return
    default:
        errutil.AbortAndStatus(c, 500)
        return
    }
}

func flag(c *gin.Context) {
    name := c.Query("userid")
    if name == "" {
        errutil.AbortAndStatus(c, 404)
        return
    }
    ins := instance.GetInstance(name)
    if ins == nil {
        errutil.AbortAndStatus(c, 404)
        return
    }
    c.String(200, ins.GetFlag())
}

func create(c *gin.Context) {
    var name string
    var err error
    var session sessions.Session
    switch config.ServiceMode {
    case "web":
        token := c.PostForm("token")
        if !captcha.Verify(c) {
            c.String(400, "Captcha verification failed.")
            return
        }
        name, err = auth.Auth(token)
        if err != nil {
            c.String(400, "Invalid Token.")
            return
        }
    case "api":
        var user userdata
        if err := c.ShouldBindJSON(&user); err != nil {
            errutil.AbortAndStatus(c, 400)
            return
        }
        name = strconv.Itoa(user.ID)
    default:
        errutil.AbortAndStatus(c, 500)
        return
    }
    _, err = instance.Up(name)
    if err != nil {
        panic(err)
    }
    switch config.ServiceMode {
    case "web":
        session = sessions.Default(c)
        session.Set("name", name)
        session.Save()
        c.Redirect(301, "/")
        return
    case "api":
        c.JSON(200, true)
        return
    default:
        errutil.AbortAndStatus(c, 500)
        return
    }
}

func destroy(c *gin.Context) {
    var name string
    var err error
    var session sessions.Session
    switch config.ServiceMode {
    case "web":
        session = sessions.Default(c)
        name, _ = session.Get("name").(string)
    case "api":
        var user userdata
        if err := c.ShouldBindJSON(&user); err != nil {
            errutil.AbortAndStatus(c, 400)
            return
        }
        name = strconv.Itoa(user.ID)
    default:
        errutil.AbortAndStatus(c, 500)
        return
    }
    err = instance.Down(name)
    if err != nil {
        panic(err)
    }
    switch config.ServiceMode {
    case "web":
        session.Clear()
        session.Save()
        c.Redirect(301, "/")
        return
    case "api":
        c.JSON(200, true)
        return
    default:
        errutil.AbortAndStatus(c, 500)
        return
    }
}
