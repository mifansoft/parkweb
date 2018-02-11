package dto

type AuthDto struct {
	Account  string `json:"username"`
	Password string `json:"-"`
	Duration string `json:"-"`
	LoginIp  string `json:"login_ip"`
	Code     int    `json:"retCode"`
	Msg      string `json:"retMsg"`
	Success  bool   `json:"success"`
	Token    string `json:token`
}
