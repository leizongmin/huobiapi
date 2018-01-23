package client

type Client struct {
	Sign    *Sign
	options ClientOptions
}

type ClientOptions struct {
	AccessKeyId      string
	AccessKeySecret  string
	SignatureMethod  string
	SignatureVersion string
	Host             string
}

func NewClient(options ClientOptions) *Client {
	return &Client{
		Sign:    NewSign(options.AccessKeyId, options.AccessKeySecret),
		options: options,
	}
}

func (c *Client) Request(method, path string, data ParamData) ([]byte, error) {
	return SendRequest(c.Sign, method, c.options.Host, path, data)
}
