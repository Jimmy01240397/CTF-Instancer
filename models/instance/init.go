package instance

import (
    "sync"
    "time"
    
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
    "github.com/Jimmy01240397/CTF-Instancer/utils/database"
)

func init() {
    lock = new(sync.RWMutex)
    idmap = make(map[string]*instance)
    usermap = make(map[string]*instance)
    portmap = make(map[uint16]*instance)
    portpool = portqueue{rangedata{
        min: uint64(config.MinPort),
        max: uint64(config.MaxPort),
    }}

    subnetpool = subnetqueue{rangedata{
        min: 0,
        max: uint64(config.SubNetPool.SubnetCount(uint(config.Prefix)) - 1),
    }}

    database.GetDB().AutoMigrate(&instance{})
    var data []*instance
    database.GetDB().Model(&instance{}).Find(&data)
    for _, a := range data {
        if time.Now().After(a.ExpiredAt) {
            err := a.down()
            if err == nil {
                database.GetDB().Delete(a)
            }
        } else {
            portpool.Remove(a.Port)
            for _, subnet := range a.SubNets {
                subnetpool.Remove(subnet)
            }
            idmap[a.ID] = a
            usermap[a.User] = a
            portmap[a.Port] = a
            go a.expired()
        }
    }
}
