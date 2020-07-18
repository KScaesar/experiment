package redis

type Config struct {
	User     string `key:"user"`
	Password string `key:"password"`
	Address  string `key:"address"`
}

func NewClient(cfg Config) *Client {
	return &Client{}
}

type Client struct {
}
