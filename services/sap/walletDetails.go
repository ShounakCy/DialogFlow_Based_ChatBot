package sap

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	sapmodel "github.com/havells/nlp/models/sapmodel"
	log "github.com/sirupsen/logrus"
)

func walletDetails(r *sapmodel.WalletReq) (*sapmodel.WalletResp, error) {
	data, err := r.WalletDetailsMarshal()
	if err != nil {
		log.Errorf("Error marhsaling wallet request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", walletURL, walletUser, walletPswd, data)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("Error fetching Wallet info : %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Infof("Unknown response status code %v", res.StatusCode)
		return nil, errors.New("Unknown status code")
	}
	br := &sapmodel.WalletResp{}
	if err := json.NewDecoder(res.Body).Decode(br); err != nil {
		log.Errorf("Error parsing Wallet response : %v", err)
		return nil, err
	}
	return br, err
}
