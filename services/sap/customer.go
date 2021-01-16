package sap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"

	log "github.com/sirupsen/logrus"
)

//RegCustomerURL ...
var (
	RegCustomerURL string
)

//RegCustomer :
type RegCustomer struct {
	Phone string
	DF    models.DfResp
}

//InitRegCustomer :
func InitRegCustomer(url string) {
	log.Infof("Initializing customer registration ....%v", url)
	RegCustomerURL = url
	log.Infof("Customer registration initialized ....")
	return
}

//GetFFResp :::::::::::::::::::::::::::::::::::::::::::::::::::;
func (b *RegCustomer) GetFFResp(ctx context.Context) (interface{}, error) {
	log.Infof("sap/customer.go, fetching Brochure info %s,%s,%s", RegCustomerURL)
	regReq, err := sapmodel.NewRegRequest(b.Phone, b.DF)
	if err != nil {
		return fmt.Sprintf("%v", err), err
	}
	br, err := registerUser(regReq)
	if err != nil {
		return nil, err
	}
	return br.Make()
}

func registerUser(r *sapmodel.RegReq) (*sapmodel.RegResp, error) {
	data, err := r.Marshal()
	if err != nil {
		log.Errorf("Error marhsaling registraion request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", RegCustomerURL, "", "", data)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("Error fetching customerReg info : %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Infof("Unknown response status code %v", res.StatusCode)
		return nil, errors.New("Unknown status code")
	}
	br := &sapmodel.RegResp{}
	if err := json.NewDecoder(res.Body).Decode(br); err != nil {
		log.Errorf("Error parsing leave response : %v", err)
		return nil, err
	}
	return br, err
}
