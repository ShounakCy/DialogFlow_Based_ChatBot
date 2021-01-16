package sap

import (
	sapmodel "github.com/havells/nlp/models/sapmodel"
	"github.com/micro/go-micro/config"
	log "github.com/sirupsen/logrus"
)

func categoryDetails() ([]sapmodel.PrdCategories, error) {

	var br []sapmodel.PrdCategories
	if err := config.Get("product_categories").Scan(&br); err != nil {
		log.Errorf("Error reading product categories from config file : %v", err)
		return nil, err
	}

	return br, nil
}
