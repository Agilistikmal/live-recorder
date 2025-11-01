package configs

type RequestConfig struct {
	UserAgent string
	Referer   string
	Cookie    string
}

func NewRequestConfig(userAgent string, referer string, cookie string) RequestConfig {
	return RequestConfig{
		UserAgent: userAgent,
		Referer:   referer,
		Cookie:    cookie,
	}
}
