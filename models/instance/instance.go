package instance

import (
    "fmt"
    "regexp"
    "time"
    "os"
    "os/exec"
    "path"
    "slices"

    netaddr "github.com/dspinhirne/netaddr-go"

    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

type instance struct {
    ID string
    User string
    Flag string
    Ports []uint16
    ExpiredAt time.Time
    SubNets netaddr.IPv4NetList
}

func (c *instance) up() error {
    cmd := exec.Command("docker", "compose", "-p", c.ID, "up", "-d")
    cmd.Dir = config.ChalDir
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("ID=%s", c.ID))
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("FLAG=%s", c.GetFlag()))
    for i, port := range c.Ports {
        cmd.Env = append(cmd.Environ(), fmt.Sprintf("PORT%d=%d", i, port))
    }
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
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("ID=%s", c.ID))
    cmd.Env = append(cmd.Environ(), fmt.Sprintf("FLAG=%s", c.GetFlag()))
    for i, port := range c.Ports {
        cmd.Env = append(cmd.Environ(), fmt.Sprintf("PORT%d=%d", i, port))
    }
    for i, subnet := range c.SubNets {
        cmd.Env = append(cmd.Environ(), fmt.Sprintf("SUBNET%d=%s", i, subnet.String()))
    }
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func (c *instance) GetFlag() string {
    return fmt.Sprintf("%s{%s_%s}", config.FlagPrefix, config.FlagMsg, c.Flag)
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

func removeDuplicates(datas []string) (result []string) {
    for _, value := range datas {
        if !slices.Contains(result, value) {
            result = append(result, value)
        }
    }
    return
}

func newinstance(user string) (*instance, bool, error) {
    if data, exist := usermap[user]; exist {
        return data, exist, nil
    }
    composefile, err := os.ReadFile(path.Join(config.ChalDir, "docker-compose.yml"))
    if err != nil {
        return nil, false, err
    }
    reg := regexp.MustCompile(`\$\{SUBNET[0-9]+\}`)
    subnetlength := len(removeDuplicates(reg.FindAllString(string(composefile), -1)))
    subnets := make(netaddr.IPv4NetList, subnetlength)
    for i := 0; i < subnetlength; i++ {
        subnets[i], err = gensubnet()
        if err != nil {
            for j := 0; j < i; j++ {
                subnetpool.Push(subnets[j])
            }
            return nil, false, err
        }
    }

    reg = regexp.MustCompile(`\$\{PORT[0-9]+\}`)
    portlength := len(removeDuplicates(reg.FindAllString(string(composefile), -1)))
    ports := make([]uint16, portlength)
    for i := 0; i < portlength; i++ {
        ports[i], err = genport()
        if err != nil {
            for j := 0; j < i; j++ {
                portpool.Push(ports[j])
            }
            return nil, false, err
        }
    }

    ins := instance{
        ID: genid(),
        User: user,
        Flag: genid(),
        Ports: ports,
        ExpiredAt: time.Now().Add(config.Validity),
        SubNets: subnets,
    }

    idmap[ins.ID] = &ins
    usermap[ins.User] = &ins
    os.Mkdir(fmt.Sprintf("/tmp/%s", ins.ID), 0755)
    os.WriteFile(fmt.Sprintf("/tmp/%s/userid", ins.ID), []byte(ins.User), 0644)
    os.WriteFile(fmt.Sprintf("/tmp/%s/flag", ins.ID), []byte(ins.GetFlag()), 0644)
    return &ins, false, nil
}

func (c *instance) del() {
    if data, exist := usermap[c.User]; exist && data == c {
        os.RemoveAll(fmt.Sprintf("/tmp/%s", c.ID))
        delete(idmap, c.ID)
        delete(usermap, c.User)
        for _, port := range c.Ports {
            portpool.Push(port)
        }
        for _, subnet := range c.SubNets {
            subnetpool.Push(subnet)
        }
    }
}
