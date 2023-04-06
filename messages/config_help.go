package messages

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
