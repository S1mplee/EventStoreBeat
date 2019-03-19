// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Brokers  []string `config:"brokers"`
	Stream   string   `config:"stream"`
	Password string   `config:"password"`
	Username string   `config:"username"`
}

var DefaultConfig = Config{
	Brokers:  []string{"http://localhost:2113"},
	Stream:   "$ce-holdings.allModels",
	Password: "changeit",
	Username: "admin",
}
