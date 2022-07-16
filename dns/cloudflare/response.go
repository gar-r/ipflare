package cloudflare

type Response[T any] struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Result   T        `json:"result"`
}

type Zone struct {
	Id string `json:"id"`
}

type Record struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}
