package models

//AppReq :
type AppReq struct {
	Msg          string `json:"msg"`
	SessionID    string `json:"session_id"`
	LangCode     string `json:"lang_code"`
	MobileNumber string `json:"mobile_num"`
	UserCode     string `json:"user_code"`
	}

//AppResp :
type AppResp struct {
	Data interface{} `json:"msg"`
	SessionID	string	`json:"session_id"`
}
