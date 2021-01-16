package sap

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"
// 	"fmt"
// 	"github.com/havells/nlp/models"
// 	sapmodel "github.com/havells/nlp/models/sapmodel"
// 	repo "github.com/havells/nlp/repo"
// 	log "github.com/sirupsen/logrus"
// )

// //SrURL ...
// var (
// 	SrURL  string
// 	PdURL  string
// 	AddURL string
// 	IsuURL string
// )

// //InitSRservice :
// func InitSRservice(url string) {
// 	log.Infof("Initializing service request ....%v", url)
// 	SrURL = url
// 	log.Infof(" service request initialized ....")
// 	return
// }

// //InitIssueservice :
// func InitIssueservice(url string) {
// 	log.Infof("Initializing issue request ....%v", url)
// 	IsuURL = url
// 	log.Infof(" issue request initialized ....")
// 	return
// }

// //InitProductdetailsservice :
// func InitProductdetailsservice(url string) {
// 	log.Infof("Initializing Product details request ....%v", url)
// 	PdURL = url
// 	log.Infof(" Product details request initialized ....")
// 	return
// }

// //InitAddservice :
// func InitAddservice(url string) {
// 	log.Infof("Initializing Address request ....%v", url)
// 	AddURL = url
// 	log.Infof(" Address request initialized ....")
// 	return
// }

// //SRservice :
// type SRservice struct {
// 	Phone     string
// 	SessionID string
// 	DF        models.DfResp
// }

// //GetFFResp :::::::::::::::::::::::::::::::::::::::::::::::::::;
// func (b *SRservice) GetFFResp(ctx context.Context) (interface{}, error) {
// 	c := http.Client{Timeout: time.Second * 60}
// 	log.Info(" MobileNumber :", b.Phone)
// 	log.Info(" SessionID : ", b.SessionID)
// 	data, err := sapmodel.NewSRrequest(b.Phone, b.SessionID, b.DF)
// 	log.Info("error : ", err)
// /////////////////
// // 	if err != nil && (err == sapmodel.ErrFallBackRegistration){

// // 		go isRegistered(b.Phone, b.SessionID)
// // 		if {

// // 		}
// // 		return
// // }

// 	//Product Details Request::::::::::::::::::::::::::::::::::::::::::::::::::;
// 	if err != nil && (err == sapmodel.ErrFallBackProduct) {
// 		log.Info("payload : ", data)
// 		req, err := NewHTTPReq("POST", PdURL, "", "", data)
// 		if err != nil {
// 			return "", err
// 		}
// 		res, err := c.Do(req)
// 		if err != nil {
// 			log.Errorf("Error fetching Product details info : %v", err)
// 			return "", err
// 		}
// 		///////////////////////////////////////////////////////////////
// 		reqAdd, err := NewHTTPReq("POST", AddURL, "", "", data)
// 		if err != nil {
// 			return "", err
// 		}
// 		resAdd, err := c.Do(reqAdd)
// 		if err != nil {
// 			log.Errorf("Error fetching Address details info : %v", err)
// 			return "", err
// 		}

//     	brAdd := []sapmodel.AddressResp{}
// 		if err := json.NewDecoder(resAdd.Body).Decode(&brAdd); err != nil {
// 			log.Errorf("Error parsing Address details response : %v", err)
// 			return "", err
// 		}
// 		go makeAddressResp(b.SessionID, brAdd)
// 		///////////////////////////////////////////////////////////////
// 		defer res.Body.Close()
// 		br := []sapmodel.ProductDetailsResp{}
// 		if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
// 			log.Errorf("Error parsing Product details response : %v", err)
// 			return "", err
// 		}
// 		return makeProductResp(b.SessionID, br), sapmodel.ErrFallBackProduct
// 	}
// 	//ProductType Details Request::::::::::::::::::::::::::::::::::::::::;;
// 	if err != nil && (err == sapmodel.ErrFallBackProductType) {
// 		log.Info("payload : ", data)
// 		req, err := NewHTTPReq("POST", PdURL, "", "", data)
// 		if err != nil {
// 			return "", err
// 		}
// 		res, err := c.Do(req)
// 		if err != nil {
// 			log.Errorf("Error fetching Product details info : %v", err)
// 			return "", err
// 		}///////////////////////////////////////////////////////////////
// 		reqAdd, err := NewHTTPReq("POST", AddURL, "", "", data)
// 		if err != nil {
// 			return "", err
// 		}
// 		resAdd, err := c.Do(reqAdd)
// 		if err != nil {
// 			log.Errorf("Error fetching Address details info : %v", err)
// 			return "", err
// 		}

//     	brAdd := []sapmodel.AddressResp{}
// 		if err := json.NewDecoder(resAdd.Body).Decode(&brAdd); err != nil {
// 			log.Errorf("Error parsing Address details response : %v", err)
// 			return "", err
// 		}
// 		go makeAddressResp(b.SessionID, brAdd)
// 		///////////////////////////////////////////////////////////////
// 		defer res.Body.Close()
// 		br := []sapmodel.ProductDetailsResp{}
// 		if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
// 			log.Errorf("Error parsing Product details response : %v", err)
// 			return "", err
// 		}
// 		return makeProductTypeResp(b.SessionID, br, b.DF), sapmodel.ErrFallBackProductType
// 	}

// 	// Issue Details request::::::::::::::::::::::::::::::::::::::;
// 	if err != nil && (err == sapmodel.ErrFallBackIssue) {
// 		req, err := NewHTTPReq("POST", IsuURL, "", "", data)
// 		if err != nil {
// 			return "", err
// 		}
// 		res, err := c.Do(req)
// 		if err != nil {
// 			log.Errorf("Error fetching Issue details info : %v", err)
// 			return "", err
// 		}
// 		defer res.Body.Close()
// 		br := []sapmodel.ProductIssueResp{}
// 		if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
// 			log.Errorf("Error parsing Issue details response : %v", err)
// 			return "", err
// 		}
// 		go makeProductIssueResp(b.SessionID, br)
// 		return sapmodel.ComplaintDescription, sapmodel.ErrFallBackIssue
// 	}

// 	// Address details request :::::::::::::::::::::::::::;
// 	if err != nil && (err == sapmodel.ErrFallBackAddress) {
// 		req, err := NewHTTPReq("POST", AddURL, "", "", data)
// 		if err != nil {
// 			return "", err
// 		}
// 		res, err := c.Do(req)
// 		if err != nil {
// 			log.Errorf("Error fetching Address details info : %v", err)
// 			return "", err
// 		}
// 		defer res.Body.Close()
// 		br := []sapmodel.AddressResp{}
// 		if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
// 			log.Errorf("Error parsing Address details response : %v", err)
// 			return "", err
// 		}
// 		return makeAddressResp(b.SessionID, br), sapmodel.ErrFallBackAddress
// 	}
// 	// Confirmation response:::::::::::::::::::::::::::;
// 	if err != nil && (err == sapmodel.ErrFallBackService || err == sapmodel.ErrEndSessionService) {
// 		return string(data), err
// 	}
// 	// Service Request:::::::::::::::::::::::::::;
// 	if err != nil {
// 		return "", err
// 	}
// 	req, err := NewHTTPReq("POST", SrURL, "", "", data)
// 	if err != nil {
// 		return "", err
// 	}
// 	res, err := c.Do(req)
// 	if err != nil {
// 		log.Errorf("Error fetching Service details info : %v", err)
// 		return "", err
// 	}
// 	defer res.Body.Close()
// 	br := &sapmodel.SRresponse{}
// 	if err := json.NewDecoder(res.Body).Decode(&br); err != nil {
// 		log.Errorf("Error parsing service details response : %v", err)
// 		return "", err
// 	}
// 	return makeSRresponse(br), nil
// }

// //makeProductresponse ::::::::::::::::::::::::::::::::::::::::::::::::::::::::
// func makeProductResp(key string, r []sapmodel.ProductDetailsResp) string {

// 	var ans []string
// 	var serials []string
// 	var productnames []string
// 	for i, v := range r {
// 		if fmt.Sprintf("%v",v.ProductCategory) == ""{
// 			return "I'm sorry but you have no product registered"
// 		}
// 		ProductCategory := fmt.Sprintf("%v",v.ProductCategory)
// 		ProductSubCategory := fmt.Sprintf("%v",v.ProductSubCategory)
// 		ProductName := fmt.Sprintf("%v",v.ProductName)
// 		resp := fmt.Sprintf("%d",i+1) + ". Product Name : " + ProductName + "\nProduct Category : " + ProductCategory +"\nProduct Sub-Category : " + ProductSubCategory + "\n\n"
// 		ans = append(ans, resp)
// 		serials = append(serials, fmt.Sprintf("%v",v.SerialNumber))
// 		productnames = append(productnames, fmt.Sprintf("%v",v.ProductName))

// 	}
// 	finalResp := "These are the list of products under your account :\n\n" + strings.Join(ans, "") + "\n\nFollowing are the list of Products, please type the Index Number from the above list of products for which you want to raise a Service Request ."

// 	go repo.UpdateProductMap(key, serials)
// 	go repo.UpdateProductNameMap(key, productnames)
// 	log.Info("getting results ::::::::::::::::::::::::::::")

// 	return finalResp
// }

// // makeProductTypeResp :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::;
// func makeProductTypeResp(key string, r []sapmodel.ProductDetailsResp, ptype models.DfResp) string {

// 	var ans []string
// 	var serials []string
// 	var productnames []string
// 	ProductType := ptype.Fields["ProductType"].GetStringValue()
// 	for c, v := range r {
// 		if strings.Contains(fmt.Sprintf("%v",v.ProductCategory), ProductType) {
// 			ProductCategory := fmt.Sprintf("%v",v.ProductCategory)
// 			ProductSubCategory := fmt.Sprintf("%v",v.ProductSubCategory)
// 			resp := fmt.Sprintf("%d",c+1) + ". Product Name : " + fmt.Sprintf("%v",v.ProductName) + "\nProduct Category : " + ProductCategory +"\nProduct Sub-Category : " + ProductSubCategory + "\n\n"
// 			ans = append(ans, resp)
// 			serials = append(serials, fmt.Sprintf("%v",v.SerialNumber))
// 			productnames = append(productnames, fmt.Sprintf("%v",v.ProductName))
// 		}
// 	}
// 	finalResp := "These are the list of products under your account :\n\n" + strings.Join(ans, "") + "\n\nFollowing are the list of Products, please type the Index Number from the above list of products for which you want to raise a Service Request ."
// 	go repo.UpdateProductMap(key, serials)
// 	go repo.UpdateProductNameMap(key, productnames)
// 	log.Infof("getting results ::::::::::::::::::::::::::::, %v", finalResp)
// 	return finalResp
// }

// //makeProductIssueResp ::::::::::::::::::::::::::::::::::::::::::::::::
// func makeProductIssueResp(key string, r []sapmodel.ProductIssueResp) {

// 	var issCache []repo.IssueCache
// 	for _, v := range r {
// 		issCache = append(issCache, repo.IssueCache{Issue: v.Name, Guid: v.GUID})
// 	}
// 	go repo.UpdateIssueMap(key, issCache)
// 	return
// }

// // makeAddressResp :::::::::::::::::::::::::::::::::::::::::::::::::::
// func makeAddressResp(key string, r []sapmodel.AddressResp) string {

// 	var ans []string
// 	var addresses []string
// 	var fulladdress []string
// 	for i, v := range r {
// 		FullAddresses := v.FullAddress
// 		resp := strconv.Itoa(i+1) + ". " + FullAddresses + "\n\n"
// 		ans = append(ans, resp)
// 		addresses = append(addresses, v.GUID)
// 		fulladdress = append(fulladdress, v.FullAddress)
// 	}
// 	finalResp := "Here is the list of registered addresses :\n\n" + strings.Join(ans, "") + "\nFollowing are the list of Addresses, please type the Index Number from the above list of addresses for which you want to raise the Service Request."
// 	go repo.UpdateAddressMap(key, addresses)
// 	go repo.UpdateFullAddressMap(key, fulladdress)
// 	return finalResp
// }

// //makeSRresponse : final response::::::::::::::::::::::::::::;
// func makeSRresponse(r *sapmodel.SRresponse) string {
// 	if r.JobID !=""{
// 	resp := "Your Job ID is : " + r.JobID + "\n\nYour Service Request has been successfully raised"
// 	return resp
// 	}
// 	return "I'm sorry, but we are not able to raise a Service Request for you ! \nPlease try again after sometime."
// 	}

// // func toCharStr(i int) string {
// // 	return string('a' + i)
// // }

// // func isRegistered(phone, key string)(*sapmpdel.RegResp, error){
// // 	r := sapmodel.RegReq{Mobile:phone, Source: 4}
// // 	res, err := registerUser(phone, key)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	if res.StatusCode == "208" {
// // 		go repo.UpdateCustomerGUID(key, res.GUID)
// // 		return res, sapmodel.ErrFallBackRegistration
// // 	}
// // 	return nil, nil
// // }
