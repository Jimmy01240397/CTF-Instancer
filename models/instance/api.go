package instance

import (
    "sync"
)

var lock *sync.RWMutex

func Up(user string) (ins *instance, err error) {
    lock.Lock()
    defer lock.Unlock()
    var exist bool
    ins, exist, err = newinstance(user)
    if err != nil {
        return
    }
    if !exist {
        err = ins.up()
        if err != nil {
            ins.del()
            return
        }
        go ins.expired()
    }
    return
}

func Down(user string) error {
    lock.Lock()
    defer lock.Unlock()
    if data, exist := usermap[user]; exist {
        err := data.down()
        if err != nil {
            return err
        }
        data.del()
    }
    return nil
}

func GetInstance(user string) (ins *instance) {
    var exist bool
    if ins, exist = usermap[user]; !exist {
        ins = nil
    }
    return
}

func GetIDMap() map[string]*instance {
    return idmap
}
