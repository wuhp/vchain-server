package handler

import (
    "fmt"
    "log"
    "errors"
    "os/exec"

    "gateway/model"
)

type MysqlConnection struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
}

func generateMysqlConnection(pid int64) *MysqlConnection {
    is := listActiveInstance()
    if len(is) == 0 {
        panic(errors.New("ERROR: no available mysql server"))
    }

    var host string
    var port int
    min := 10000

    for _, instance := range is {
        length := len(queryMappingByMysqlHostPort(instance.Host, instance.Port))
        if length < min {
            min = length
            host = instance.Host
            port = instance.Port
        }
    }

    mc := new(MysqlConnection)
    mc.Host = host
    mc.Port = port
    mc.User = "root"
    mc.Password = "root"
    mc.Database = fmt.Sprintf("vchain_%d", pid)

    return mc
}

func provisionMysqlInstance(mc *MysqlConnection) {
    out, err := exec.Command(
        "/vchain/server/bin/provision_mysql.sh",
        mc.Host, fmt.Sprintf("%d", mc.Port), mc.User, mc.Password, mc.Database,
    ).Output()
    if err != nil {
        log.Printf("ERROR: output = %s, err = %v\n", out, err)
        panic(errors.New("ERROR: fail to provision mysql instance"))
    }

    log.Printf("INFO: Successfully privision mysql instance %v\n", mc)
}

func migrateMysqlInstance(mc *MysqlConnection) {
    out, err := exec.Command(
        "/vchain/server/bin/migrate_mysql.sh",
        mc.Host, fmt.Sprintf("%d", mc.Port), mc.User, mc.Password, mc.Database,
    ).Output()
    if err != nil {
        log.Printf("ERROR: output = %s, err = %v\n", out, err)
        panic(errors.New("ERROR: fail to migrate mysql instance"))
    }

    log.Printf("INFO: Successfully migrate mysql instance %v\n", mc)
}

func unprovisionMysqlInstance(mc *MysqlConnection) {
    out, err := exec.Command(
        "/vchain/server/bin/unprovision_mysql.sh",
        mc.Host, fmt.Sprintf("%d", mc.Port), mc.User, mc.Password, mc.Database,
    ).Output()
    if err != nil {
        log.Printf("ERROR: output = %s, err = %v\n", out, err)
        panic(errors.New("ERROR: fail to unprovision mysql instance"))
    }

    log.Printf("INFO: Successfully unprivision mysql instance %v\n", mc)
}

func listActiveInstance() []*model.Instance {
    conditions := make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("active", "=", 1))

    return model.ListInstance(conditions, nil, nil)
}

func queryMappingByMysqlHostPort(host string, port int) []*model.Mapping {
    conditions := make([]*model.Condition, 0)
    conditions = append(conditions, model.NewCondition("host", "=", host))
    conditions = append(conditions, model.NewCondition("port", "=", port))

    return model.ListMapping(conditions, nil, nil)
}
