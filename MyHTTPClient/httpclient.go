package myhttpclient

func Get(rawURL string) (Response, error) {
	req := Request{
		Method: "GET",
		URL:    rawURL,
		Headers: map[string]string{
			"User-Agent": "MyHTTPClient/1.0",
			"Accept":     "*/*",
		},
	}
	return Fetch(req, 0)
}
