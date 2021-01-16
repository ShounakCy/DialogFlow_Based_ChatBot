package sap

import (
	"context"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	log "github.com/sirupsen/logrus"
)

var (
	creditURL  string
	creditUser string
	creditPswd string
)

//InitCreditService :
func InitCreditService(url, user, pswd string) {
	log.Infof("Initializing credit service %v", url)
	creditURL = url
	creditUser = user
	creditPswd = pswd
	log.Infof("Credit service initialized ....")
	return
}

//CreditService :
type CreditService struct {
	Code  string
	KKBER string
	DF    models.DfResp
}

//GetFFResp :
func (b *CreditService) GetFFResp(ctx context.Context) (interface{}, error) {

	_, err := sapmodel.NewCreditRequest(b.Code, b.KKBER, b.DF)
	log.Infof("error............... %v", err)
	if err != nil {

		switch err {
		case sapmodel.ErrCreditRequest:

			creditResp, _ := makeCreditResp(b.Code, b.KKBER)
			return creditResp, sapmodel.ErrEndSessionCredit

		default:
			return "Hi Dealer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", sapmodel.ErrEndSessionCredit
		}
	}

	return "", nil
}

func makeCreditResp(code, kkber string) (string, error) {

	req := sapmodel.CreditReq{Code: code, KKBER: kkber}
	res, err := creditDetails(&req)
	if err != nil {
		return "", nil
	}

	resp := `Hi Dealer,
		
Your credit summary is as follows:

Credit limit :      ₹ ` + res.EvLimit + "\n" +
		`Used :                ₹ ` + res.EvUsed + "\n" +
		`Available limit : ₹ ` + res.EvBalance

	return resp, nil

}
