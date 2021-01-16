package sap

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	sapmodel "github.com/havells/nlp/models/sapmodel"

	log "github.com/sirupsen/logrus"
)

func validatePin(r *sapmodel.PinCodeReq) (*sapmodel.AddrPinCode, error) {
	data, err := r.ValidateMarshal()
	if err != nil {
		log.Errorf("Error marhsaling validate PinCode request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", VPinURL, "", "", data)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("Error fetching PinCode Request info : %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Infof("Unknown response status code %v", res.StatusCode)
		return nil, errors.New("Unknown status code")
	}
	br := &sapmodel.AddrPinCode{}
	if err := json.NewDecoder(res.Body).Decode(br); err != nil {
		log.Errorf("Error parsing PinCode response : %v", err)
		return nil, err
	}
	return br, err
}
