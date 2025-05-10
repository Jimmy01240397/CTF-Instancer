package config

import (
    "os"
    "strconv"
    "time"
    "crypto/rand"
    "fmt"
    "strings"

    netaddr "github.com/dspinhirne/netaddr-go"
)

var Debug bool
var Port string
var ServiceMode string
var Secret []byte
var Sessionname string
var Token string
var Title string
var MaxPort uint16
var MinPort uint16
var Validity time.Duration
var BaseScheme string
var BaseHost string
var ChalDir string
var SubNetPool *netaddr.IPv4Net
var Prefix uint8
var FlagPrefix string
var FlagMsg string
var CAPTCHA_SRC string
var CAPTCHA_CLASS string
var CAPTCHA_SITE_KEY string
var CAPTCHA_SECRET_KEY string
var CAPTCHA_BACKEND string
var CAPTCHA_RESPONSE_NAME string

const (
    Forward = 0
    Proxy = 1
    Command = 2
)

func init() {
    loadenv()
    var err error
    debugstr, exists := os.LookupEnv("DEBUG")
    if !exists {
        Debug = false
    } else {
        Debug, err = strconv.ParseBool(debugstr)
        if err != nil {
            Debug = false
        }
    }
    Port = os.Getenv("PORT")
    Secret = make([]byte, 12)
    _, err = rand.Read(Secret)
    if err != nil {
        panic(err)
    }
    Sessionname = os.Getenv("SESSIONNAME")
    Title = os.Getenv("TITLE")
    Token = os.Getenv("TOKEN")
    FlagPrefix = os.Getenv("FLAGPREFIX")
    FlagMsg = os.Getenv("FLAGMSG")
    ChalDir = os.Getenv("CHALDIR")
    BaseScheme = os.Getenv("BASESCHEME")
    BaseHost = os.Getenv("BASEHOST")
    CAPTCHA_SRC = os.Getenv("CAPTCHA_SRC")
    CAPTCHA_CLASS = os.Getenv("CAPTCHA_CLASS")
    CAPTCHA_SITE_KEY = os.Getenv("CAPTCHA_SITE_KEY")
    CAPTCHA_SECRET_KEY = os.Getenv("CAPTCHA_SECRET_KEY")
    CAPTCHA_BACKEND = os.Getenv("CAPTCHA_BACKEND")
    CAPTCHA_RESPONSE_NAME = os.Getenv("CAPTCHA_RESPONSE_NAME")
    ServiceMode, exists = os.LookupEnv("SERVICEMODE")
    if !exists {
        ServiceMode = "web"
    }
    subnetpoolstr, exists := os.LookupEnv("SUBNETPOOL")
    if !exists {
        SubNetPool, _ = netaddr.ParseIPv4Net("172.16.0.0/16")
    } else {
        SubNetPool, err = netaddr.ParseIPv4Net(subnetpoolstr)
        if err != nil {
            SubNetPool, _ = netaddr.ParseIPv4Net("172.16.0.0/16")
        }
    }
    prefixstr, exists := os.LookupEnv("PREFIX")
    if !exists {
        Prefix = 24
    } else {
        tmp, err := strconv.ParseUint(prefixstr, 10, 8)
        Prefix = uint8(tmp)
        if err != nil {
            Prefix = 24
        }
    }
    maxportstr, exists := os.LookupEnv("MAXPORT")
    if !exists {
        MaxPort = 30000
    } else {
        tmp, err := strconv.ParseUint(maxportstr, 10, 16)
        MaxPort = uint16(tmp)
        if err != nil {
            MaxPort = 30000
        }
    }
    minportstr, exists := os.LookupEnv("MINPORT")
    if !exists {
        MinPort = 30000
    } else {
        tmp, err := strconv.ParseUint(minportstr, 10, 16)
        MinPort = uint16(tmp)
        if err != nil {
            MinPort = 30000
        }
    }
    validitystr, exists := os.LookupEnv("VALIDITY")
    if !exists {
        Validity = 3 * time.Minute
    } else {
        Validity, err = time.ParseDuration(validitystr)
        if err != nil {
            Validity = 3 * time.Minute
        }
    }
}

func GetMode(index int) int {
    modestr, exists := os.LookupEnv(fmt.Sprintf("MODE%d", index))
    if !exists {
        return Forward
    }
    modestr = strings.ToLower(modestr)
    switch modestr {
    case "forward":
        return Forward
    case "proxy":
        return Proxy
    case "command":
        return Command
    default:
        return Forward
    }
}

func GetCommand(index int) string {
    return os.Getenv(fmt.Sprintf("COMMAND%d", index))
}
