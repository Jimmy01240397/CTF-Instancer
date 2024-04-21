package instance

import (
    "fmt"
    "reflect"
    "encoding/json"
    "database/sql/driver"

    netaddr "github.com/dspinhirne/netaddr-go"
)

type subnetarr netaddr.IPv4NetList

func (c *subnetarr) Scan(value interface{}) (err error) {
    var subnetstrs []string
    if val, ok := value.(string); ok {
        err = json.Unmarshal([]byte(val), &subnetstrs)
        if err != nil {
            return
        }
    } else {
        err = fmt.Errorf("sql: unsupported type %s", reflect.TypeOf(value))
        return
    }
    var tmp netaddr.IPv4NetList
    tmp, err = netaddr.NewIPv4NetList(subnetstrs)
    *c = subnetarr(tmp)
    return
}

func (c subnetarr) Value() (driver.Value, error) {
    subnetstrs := make([]string, len(c))
    for i, subnet := range c {
        subnetstrs[i] = subnet.String()
    }
    data, err := json.Marshal(subnetstrs)
    return string(data), err
}
