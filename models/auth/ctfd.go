package auth

import (
    "fmt"
    "net/http"
    "regexp"
    "io"
    "encoding/json"
    "strconv"

    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

func Auth(identity any) (name string, err error) {
    reg := regexp.MustCompile(`^ctfd_[0-9a-f]{64}$`)
    if !reg.MatchString(identity.(string)) {
        err = fmt.Errorf("Token's format is invalid: %s", identity.(string))
        return
    }
    client := &http.Client{}
    var req *http.Request
    var res *http.Response
    var result []byte
    for _, ctfdurl := range config.CTFDURL {
        req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/me", ctfdurl), nil)
        if err != nil {
            return
        }
        req.Header.Set("Authorization", fmt.Sprintf("Token %s", identity.(string)))
        req.Header.Set("Content-Type", "application/json")
        res, err = client.Do(req)
        if err != nil {
            return
        }
        defer res.Body.Close()
        result, _ = io.ReadAll(res.Body)
        if res.StatusCode == 200 {
            var resobj map[string]any
            err = json.Unmarshal(result, &resobj)
            if success, ok := resobj["success"].(bool); ok && success {
                name = strconv.Itoa(int(resobj["data"].(map[string]any)["id"].(float64)))
                return
            }
        }
    }
    if res != nil {
        err = fmt.Errorf("Couldn't login as token: %s, status_code: %d, res: %s", identity.(string), res.StatusCode, string(result))
    } else {
        err = fmt.Errorf("Couldn't login as token")
    }
    return
}

