package discover

var gatewayHost string

func Init(gateway string) {
    gatewayHost = gateway
}

func Gateway() string {
    return gatewayHost
}
