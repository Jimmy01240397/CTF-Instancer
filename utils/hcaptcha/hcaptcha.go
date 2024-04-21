package hcaptcha

import (
    "io"
    "net/http"
    "net/url"
    "encoding/json"
    
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

func Verify(responseToken string) bool {
    data := map[string][]string{
        "secret": []string{config.HCAPTCHA_SECRET_KEY},
        "response": []string{responseToken},
    }
    qs := url.Values(data)
    res, err := http.PostForm("https://api.hcaptcha.com/siteverify", qs)
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

