package sap

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	//"fmt"
	//"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	//repo "github.com/havells/nlp/repo"
	log "github.com/sirupsen/logrus"
)

func addressBook(r *sapmodel.AddressReq) ([]sapmodel.AddressResp, error) {
	data, err := r.AddressBookMarshal()
	if err != nil {
		log.Errorf("Error marhsaling registraion request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", AddURL, "", "", data)
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
	var br []sapmodel.AddressResp
	if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
		log.Errorf("Error parsing Address response : %v", err)
		return nil, err
	}
	return br, err
}
