package mysql

type Config struct {
	User     string `config:"user"`
	Password string `config:"password"`
	Host     string `config:"host"`
	Port     string `config:"port"`
}

func NewClient(cfg Config) *Client {
	return &Client{}
}

type Client struct {
}
