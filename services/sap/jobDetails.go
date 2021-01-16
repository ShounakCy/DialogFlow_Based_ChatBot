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

func jobDetails(r *sapmodel.JobStatusReq) ([]sapmodel.JobStatusResp, error) {
	data, err := r.JobDetailsMarshal()
	if err != nil {
		log.Errorf("Error marhsaling job request : %v", err)
		return nil, err
	}
	log.Info("payload : ", string(data))
	c := http.Client{Timeout: time.Second * 60}
	req, err := NewHTTPReq("POST", ServiceStatusURL, "", "", data)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		log.Errorf("Error fetching job info : %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Infof("Unknown response status code %v", res.StatusCode)
		return nil, errors.New("Unknown status code")
	}
	var br []sapmodel.JobStatusResp
	if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
		log.Errorf("Error parsing job response : %v", err)
		return nil, err
	}
	return br, err
}
