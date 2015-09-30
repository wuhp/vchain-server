package discover

var gateway string

func InitGateway(g string) error {
    gateway = g
    return nil
}

func DiscoverGateway() (string, error) {
    return gateway, nil
}
