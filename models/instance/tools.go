package instance

import (
    "regexp"
    "fmt"

    "github.com/google/uuid"
    netaddr "github.com/dspinhirne/netaddr-go"

    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

type rangedata struct {
    min uint64
    max uint64
}

type portqueue []rangedata
type subnetqueue []rangedata

var idmap map[string]*instance
var usermap map[string]*instance

var portpool portqueue
var subnetpool subnetqueue

func push(src []rangedata, num uint64) (dst []rangedata) {
    if len(src) > 0 && src[len(src) - 1].min <= src[len(src) - 1].max && src[len(src) - 1].max + 1 == num {
        src[len(src) - 1].max++
        dst = src
    } else if len(src) > 0 && src[len(src) - 1].min >= src[len(src) - 1].max && src[len(src) - 1].max - 1 == num {
        src[len(src) - 1].max--
        dst = src
    } else {
        dst = append(src, rangedata{
            min: num,
            max: num,
        })
    }
    return
}

func pop(src []rangedata) (dst []rangedata, num uint64, err error) {
    if len(src) < 1 {
        return src, 0, fmt.Errorf("Pop OOB")
    }
    num = src[0].min
    if src[0].min == src[0].max {
        dst = (src)[1:]
    } else if src[0].min < src[0].max {
        src[0].min++
        dst = src
    } else if src[0].min > src[0].max {
        src[0].min--
        dst = src
    }
    return
}

func remove(src []rangedata, num uint64) (dst []rangedata) {
    for i, _ := range src {
        if src[i].min == src[i].max && src[i].min == num {
            dst = append(src[:i], src[i+1:]...)
            return
        } else if src[i].min < src[i].max && src[i].min <= num && num <= src[i].max {
            tmp := src[i]
            dst = src[:i]
            if tmp.min < num {
                dst = append(dst, rangedata{
                    min: tmp.min,
                    max: num - 1,
                })
            }
            if num < tmp.max {
                dst = append(dst, rangedata{
                    min: num + 1,
                    max: tmp.max,
                })
            }
            dst = append(dst, src[i+1:]...)
            return
        } else if src[i].min > src[i].max && src[i].min >= num && num >= src[i].max {
            tmp := src[i]
            dst = src[:i]
            if tmp.min > num {
                dst = append(dst, rangedata{
                    min: tmp.min,
                    max: num + 1,
                })
            }
            if num > tmp.max {
                dst = append(dst, rangedata{
                    min: num - 1,
                    max: tmp.max,
                })
            }
            dst = append(dst, src[i+1:]...)
            return
        }
    }
    dst = src
    return
}


func (c *portqueue) Push(num uint16) {
    *c = portqueue(push([]rangedata(*c), uint64(num)))
}

func (c *portqueue) Pop() (num uint16, err error) {
    tmpqueue, tmpnum, errtmp := pop(*c)
    err = errtmp
    num = uint16(tmpnum)
    *c = portqueue(tmpqueue)
    return
}

func (c *portqueue) Remove(num uint16) {
    *c = portqueue(remove([]rangedata(*c), uint64(num)))
}

func (c *subnetqueue) Push(net *netaddr.IPv4Net) {
    length := config.SubNetPool.Resize(uint(config.Prefix)).Len()
    index := uint64((net.Network().Addr() - config.SubNetPool.Network().Addr()) / length)
    *c = subnetqueue(push([]rangedata(*c), index))
}

func (c *subnetqueue) Pop() (net *netaddr.IPv4Net, err error) {
    tmpqueue, tmpnum, errtmp := pop(*c)
    err = errtmp
    num := uint32(tmpnum)
    *c = subnetqueue(tmpqueue)
    net = config.SubNetPool.NthSubnet(uint(config.Prefix), num)
    return
}

func (c *subnetqueue) Remove(net *netaddr.IPv4Net) {
    length := config.SubNetPool.Resize(uint(config.Prefix)).Len()
    index := uint64((net.Network().Addr() - config.SubNetPool.Network().Addr()) / length)
    *c = subnetqueue(remove([]rangedata(*c), index))
}


func genid() (id string) {
    reg := regexp.MustCompile(`[^a-zA-Z0-9]`)
    id = reg.ReplaceAllString(uuid.NewString(), "")
    for _, exist := idmap[id]; exist; _, exist = idmap[id] {
        id = reg.ReplaceAllString(uuid.NewString(), "")
    }
    return
}

func genport() (uint16, error) {
    return portpool.Pop()
}

func gensubnet() (*netaddr.IPv4Net, error) {
    return subnetpool.Pop()
}

