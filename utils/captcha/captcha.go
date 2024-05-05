package captcha

import (
    "io"
    "net/http"
    "net/url"
    "encoding/json"
    "github.com/gin-gonic/gin"
    
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

func Verify(c *gin.Context) bool {
    if config.CAPTCHA_SECRET_KEY == "" {
        return true
    }
    data := map[string][]string{
        "secret": []string{config.CAPTCHA_SECRET_KEY},
        "response": []string{c.PostForm(config.CAPTCHA_RESPONSE_NAME)},
    }
    qs := url.Values(data)
    res, err := http.PostForm(config.CAPTCHA_BACKEND, qs)
    if err != nil {
        return false
    }
    defer res.Body.Close()
    result, _ := io.ReadAll(res.Body)
    var resobj map[string]any
    err = json.Unmarshal(result, &resobj)
    if err != nil {
        return false
    }
    if success, ok := resobj["success"].(bool); !ok || !success {
        return false
    }
    return true
}

