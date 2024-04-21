package errutil

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
)

type Err struct {
    Code int
    Msg interface{}
    Data interface{}
}

func (e *Err) Error() string {
    b, _ := json.Marshal(e)
    return string(b)
}

func HtoErr(h gin.H) *Err {
    return &Err{
        Code: h["code"].(int),
        Msg: h["msg"].(string),
        Data: h["data"],
    }
}

func (e *Err) ToH() gin.H {
    return gin.H{
        "code": e.Code,
        "msg": e.Msg,
        "data": e.Data,
    }
}

func AbortAndError(c *gin.Context, err *Err) {
    c.Abort()
    c.Error(err)
}

func AbortAndStatus(c *gin.Context, code int) {
    c.Abort()
    c.Error(&Err{
        Code: code,
    })
}
