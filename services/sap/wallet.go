package sap

import (
	"context"
	"regexp"
	"strings"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	log "github.com/sirupsen/logrus"
)

var (
	walletURL  string
	walletUser string
	walletPswd string
)

//WalletService :
type WalletService struct {
	Code string
	DF   models.DfResp
}

//InitWalletService :
func InitWalletService(url, user, pswd string) {
	log.Infof("Initializing wallet service ...%v", url)
	walletURL = url
	walletUser = user
	walletPswd = pswd
	return
}

//GetFFResp :
func (b *WalletService) GetFFResp(ctx context.Context) (interface{}, error) {
	_, err := sapmodel.NewWalletRequest(b.Code, b.DF)
	log.Infof("error............... %v", err)
	if err != nil {

		switch err {
		case sapmodel.ErrWalletRequest:

			walletResp, _ := makeWalletResp(b.Code)
			return walletResp, sapmodel.ErrEndSessionWallet

		default:
			return "Hi Dealer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", sapmodel.ErrEndSessionWallet
		}
	}
	return "", nil
}

func makeWalletResp(code string) (string, error) {

	req := sapmodel.WalletReq{Customer: code, Bukrs: 100}
	res, err := walletDetails(&req)
	if err != nil {
		return "", nil
	}

	txt := `Hi Dealer,

Your total outstanding is ₹ ` + res.OutStanding + `

Details are as below :` + `

Range(Days)` + "\t\t\t" + `Amount(₹)` + "\n"

	if str := getNonEmptyRange(res.Bucket1, res.Range1); str != "" {
		txt = txt + str + "\n"
	}
	if str := getNonEmptyRange(res.Bucket2, res.Range2); str != "" {
		txt = txt + str + "\n"
	}
	if str := getNonEmptyRange(res.Bucket3, res.Range3); str != "" {
		txt = txt + str + "\n"
	}
	if str := getNonEmptyRange(res.Bucket4, res.Range4); str != "" {
		txt = txt + str + "\n"
	}
	if str := getNonEmptyRange(res.Bucket5, res.Range5); str != "" {
		txt = txt + str + "\n"
	}
	if str := getNonEmptyRange(res.Bucket6, res.Range6); str != "" {
		txt = txt + str + "\n"
	}
	return txt, nil
}

func getNonEmptyRange(val, timeRange string) string {
	r := regexp.MustCompile(`([^"]*) *Days`)
	if len(strings.TrimSpace(val)) > 0 && len(strings.TrimSpace(timeRange)) > 0 {
		return r.ReplaceAllString(timeRange, "${1}") + "\t\t\t\t\t " + strings.TrimSpace(val)
	}
	return ""
}

func makeRange(s string) string {
	r := regexp.MustCompile(`([^"]*) *Days`)
	s = strings.Replace(strings.TrimPrefix(s, "0"), "-0", "-", 1)
	return r.ReplaceAllString(s, "${1}")
}
