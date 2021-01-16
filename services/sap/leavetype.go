package sap

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"github.com/havells/nlp/models"
// 	sapmodel "github.com/havells/nlp/models/sapmodel"
// 	log "github.com/sirupsen/logrus"
// )

// // ApplyLeaveService :
// type ApplyLeaveService struct {
// 	EmpCode string
// 	DF      models.DfResp
// }

// var (
// 	// ApplyLeaveURL :
// 	ApplyLeaveURL string
// 	// ApplyLeaveUser :
// 	ApplyLeaveUser string
// 	// ApplyLeavePswd :
// 	ApplyLeavePswd string
// )

// //InitApplyLeaveService :
// func InitApplyLeaveService(url, user, pswd string) {
// 	log.Infof("Initializing Leavtype Service %v", url)
// 	ApplyLeaveURL = url
// 	ApplyLeaveUser = user
// 	ApplyLeavePswd = pswd
// 	log.Infof("ApplyLeave service initialized ....")
// 	return
// }

// //GetFFResp :
// func (b *ApplyLeaveService) GetFFResp(ctx context.Context) (interface{}, error) {
// 	log.Infof("fetching Brochure info %s,%s,%s", ApplyLeaveURL, ApplyLeaveUser, ApplyLeavePswd)
// 	c := http.Client{Timeout: time.Second * 60}
// 	data, err := sapmodel.NewApplyLeaveReq(b.EmpCode, b.DF)
// 	if err != nil && (err == sapmodel.ErrFallBack || err == sapmodel.ErrEndSession || err == sapmodel.ErrWrongDate) {
// 		return string(data), err
// 	}
// 	req, err := NewHTTPReq("POST", ApplyLeaveURL, ApplyLeaveUser, ApplyLeavePswd, data)
// 	if err != nil {
// 		return "", err
// 	}
// 	res, err := c.Do(req)
// 	if err != nil {
// 		log.Errorf("Error fetching leave info : %v", err)
// 		return "", err
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != http.StatusOK {
// 		log.Infof("Unknow response status code %v", res.StatusCode)
// 		return "I am not able to apply leave right now, Please try again after some time.", nil
// 	}
// 	br := &sapmodel.ApplyLeaveResp{}
// 	if err := json.NewDecoder(res.Body).Decode(br); err != nil {
// 		log.Errorf("Error parsing leave response : %v", err)
// 		return nil, err
// 	}
// 	return makeApplyLeaveResp(br), nil
// }
// func makeApplyLeaveResp(r *sapmodel.ApplyLeaveResp) string {
// 	if r.TYPEFIELD == "E" {
// 		return "I'm sorry but you have already applied leave for this date !!!"
// 	}
// 	log.Infof("apply leave service response %v", r)
// 	resp := "Your leave request has been sent for approval !!!"
// 	return resp
// }
