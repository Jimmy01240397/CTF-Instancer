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
    var req *http.Request
    req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/me", config.CTFDURL), nil)
    if err != nil {
        return
    }
    req.Header.Set("Authorization", fmt.Sprintf("Token %s", identity.(string)))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    var res *http.Response
    res, err = client.Do(req)
    if err != nil {
        return
    }
    defer res.Body.Close()
    result, _ := io.ReadAll(res.Body)
    if res.StatusCode != 200 {
        err = fmt.Errorf("Couldn't login as token: %s, status_code: %d, res: %s", identity.(string), res.StatusCode, string(result))
        return
    }
    var resobj map[string]any
    err = json.Unmarshal(result, &resobj)
    if err != nil {
        return
    }
    if success, ok := resobj["success"].(bool); !ok || !success {
        err = fmt.Errorf("Couldn't login as token: %s, status_code: %d, res: %s", identity.(string), res.StatusCode, string(result))
        return
    }
    name = strconv.Itoa(int(resobj["data"].(map[string]any)["id"].(float64)))
    return
}

