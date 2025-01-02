package instance

import (
    "fmt"
    "regexp"
    "time"
    "os"
    "os/exec"
    "path"

    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
    "github.com/Jimmy01240397/CTF-Instancer/utils/database"
)

type instance struct {
    ID string `gorm:"primaryKey"`
    User string
    Port uint16
    ExpiredAt time.Time
    SubNets subnetarr
}

func (c *instance) up() error {
    cmd := exec.Command("docker", "compose", "-p", c.ID, "up", "-d")
    cmd.Dir = config.ChalDir
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("PORT=%d", c.Port))
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("ID=%s", c.ID))
    for i, subnet := range c.SubNets {
        cmd.Env = append(cmd.Environ(), fmt.Sprintf("SUBNET%d=%s", i, subnet.String()))
    }
    err := cmd.Run()
    if err != nil {
        c.down()
        return err
    }
    return nil
}

func (c *instance) down() error {
    cmd := exec.Command("docker", "compose", "-p", c.ID, "down")
    cmd.Dir = config.ChalDir
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("PORT=%d", c.Port))
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("ID=%s", c.ID))
    for i, subnet := range c.SubNets {
        cmd.Env = append(cmd.Environ(), fmt.Sprintf("SUBNET%d=%s", i, subnet.String()))
    }
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func (c *instance) expired() {
    time.Sleep(c.ExpiredAt.Sub(time.Now()))
    lock.Lock()
    defer lock.Unlock()
    err := c.down()
    if err == nil {
        c.del()
    }
}

func newinstance(user string) (*instance, bool, error) {
    if data, exist := usermap[user]; exist {
        return data, exist, nil
    }
    reg := regexp.MustCompile(`\$\{SUBNET[0-9]+\}`)
    composefile, err := os.ReadFile(path.Join(config.ChalDir, "docker-compose.yml"))
    if err != nil {
        return nil, false, err
    }
    subnetlength := len(reg.FindAllString(string(composefile), -1))
    subnets := make(subnetarr, subnetlength)
    for i := 0; i < subnetlength; i++ {
        subnets[i] = gensubnet()
    }
    ins := instance{
        ID: genid(),
        User: user,
        Port: genport(),
        ExpiredAt: time.Now().Add(config.Validity),
        SubNets: subnets,
    }
    idmap[ins.ID] = &ins
    usermap[ins.User] = &ins
    portmap[ins.Port] = &ins
    database.GetDB().Model(&instance{}).Create(&ins)
    os.Mkdir(fmt.Sprintf("/tmp/%s", ins.ID), 0755)
    os.WriteFile(fmt.Sprintf("/tmp/%s/userid", ins.ID), []byte(ins.User), 0644)
    return &ins, false, nil
}

func (c *instance) del() {
    database.GetDB().Delete(c)
    if data, exist := usermap[c.User]; exist && data == c {
        os.RemoveAll(fmt.Sprintf("/tmp/%s", c.ID))
        delete(idmap, c.ID)
        delete(usermap, c.User)
        delete(portmap, c.Port)
        portpool.Push(c.Port)
        for _, subnet := range c.SubNets {
            subnetpool.Push(subnet)
        }
    }
}
