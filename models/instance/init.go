package instance

import (
    "sync"
    
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

func init() {
    lock = new(sync.RWMutex)
    idmap = make(map[string]*instance)
    usermap = make(map[string]*instance)
    portpool = portqueue{rangedata{
        min: uint64(config.MinPort),
        max: uint64(config.MaxPort),
    }}

    subnetpool = subnetqueue{rangedata{
        min: 0,
        max: uint64(config.SubNetPool.SubnetCount(uint(config.Prefix)) - 1),
    }}
}
