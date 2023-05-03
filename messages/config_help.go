package messages

import "strings"

type Config struct {
	Name  string
	Users []DeliverConfig
}

type DeliverConfig struct {
	Host string
	Key  string
}

func MakeTargets(conf Config) *[]GroupDetails {
	var recips []GroupDetails
	for _, user_conf := range conf.Users {
		recips = append(recips, GroupDetails{
			DestinationHostPort: user_conf.Host,
			PublicKey:           ParsePublicKey(user_conf.Key),
		})
	}
	return &recips
}

// returns only the IP address of an IP:PORT string
func StripPort(host_port string) string {
	host := strings.Split(host_port, ":")[0]
	return host
}

// turns the list of Users in the Config object into a map where
// the IP address is the key and the value is the DeliverConfig struct
func MakeConfigMap(config Config) map[string]DeliverConfig {
	confMap := make(map[string]DeliverConfig)
	for _, user := range config.Users {
		confMap[StripPort(user.Host)] = user
	}
	return confMap
}
