package discover

var gateway string
var vchaind string
var consumer string

func InitGateway(g string) error {
    gateway = g
    return nil
}

func DiscoverGateway() (string, error) {
    return gateway, nil
}

func InitVchaind(v string) error {
    vchaind = v
    return nil
}

func DiscoverVchaind() (string, error) {
    return vchaind, nil
}

func InitConsumer(c string) error {
    consumer = c
    return nil
}

func DiscoverConsumer() (string, error) {
    return consumer, nil
}
