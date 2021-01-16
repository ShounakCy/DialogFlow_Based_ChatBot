package sap

import (
	"context"
	"fmt"
	"strings"

	"github.com/havells/nlp/models"
	sapmodel "github.com/havells/nlp/models/sapmodel"
	repo "github.com/havells/nlp/repo"
	log "github.com/sirupsen/logrus"
)

//SrURL ...
var (
	SrURL  string
	PdURL  string
	AddURL string
	IsuURL string
)

//InitSRservice :
func InitSRservice(url string) {
	log.Infof("Initializing service request ....%v", url)
	SrURL = url
	log.Infof(" service request initialized ....")
	return
}

//InitIssueservice :
func InitIssueservice(url string) {
	log.Infof("Initializing issue request ....%v", url)
	IsuURL = url
	log.Infof(" issue request initialized ....")
	return
}

//InitProductdetailsservice :
func InitProductdetailsservice(url string) {
	log.Infof("Initializing Product details request ....%v", url)
	PdURL = url
	log.Infof(" Product details request initialized ....")
	return
}

//InitAddservice :
func InitAddservice(url string) {
	log.Infof("Initializing Address request ....%v", url)
	AddURL = url
	log.Infof(" Address request initialized ....")
	return
}

//SRservice :
type SRservice struct {
	Phone     string
	SessionID string
	DF        models.DfResp
}

//GetFFResp :
func (b *SRservice) GetFFResp(ctx context.Context) (interface{}, error) {
	data, err := sapmodel.NewServiceRequest(b.Phone, b.SessionID, b.DF)
	log.Infof("sap/servicerequest.go, data : %v", fmt.Sprintf("%v", data))
	if err != nil {
		switch err {
		case sapmodel.ErrChkReg:
			log.Infof("sap/servicerequest.go, Checking registration : %v", b.SessionID)
			log.Infof("sap/servicerequest.go, fetching Brochure info %s,%s,%s", RegCustomerURL)
			ok, err := isRegistered(b.Phone, b.SessionID)
			if err != nil {
				return "", nil
			}
			if ok {
				log.Infof("sap/servicerequest.go, Checking address book", b.SessionID)
				respAddresses, err := isAddressBook(b.SessionID)
				if err == sapmodel.ErrFallBackAddress {
					return respAddresses, sapmodel.ErrFallBackProduct
				}
				return `Hi, I am Iyris, your friendly service bot from Havells. !
 Kindly add your ADRRESS with VALID PINCODE, we will be happy to assist you further.`, nil
			}
			return "Hi there, !\nKindly register yourself with NAME and EMAIL-ID first before raising a request, we will be happy to assist you further.", nil

		case sapmodel.ErrFallBackSRnewAddress:
			a := b.DF.Fields["address"].GetStringValue()
			z := b.DF.Fields["zip-code"].GetStringValue()

			if a != "" && z != "" {
				ok, err := isValidPin(z, b.SessionID)
				if err != nil {
					return "Please try again later", nil
				}
				if ok {
					return "Would you like to set it as Primary Address. Reply with a Yes or No.", sapmodel.ErrPin
				}
				return "Please enter 6 digit VALID PINCODE", sapmodel.ErrAddressLines
			}
			return "Please enter ADDRESS and 6 digit VALID PINCODE", sapmodel.ErrAddressLines

		case sapmodel.ErrFallBackProduct:
			log.Infof("sap/servicerequest.go, Checking Products: %v", b.SessionID)
			custproductsGUID := repo.GetCustumerGUID(b.SessionID)
			respProducts, err := isProducts(custproductsGUID, b.SessionID)

			if err == sapmodel.ErrFallBackProduct {
				log.Infof("sap/servicerequest.go, Checking the products details : %v", b.SessionID)
				return respProducts, sapmodel.ErrFallBackIssue
			}
			return "I'm sorry, but we are not able to raise a Service Request for this selected Product !\nPlease try again after sometime.", nil

		case sapmodel.ErrFallBackIssue:
			log.Infof("sap/servicerequest.go, Checking Issues : %v", b.SessionID)
			ProductIndex := b.DF.Fields["number1"].GetNumberValue()
			var ProductName string
			var SerialCode string

			if ProductIndex != 0 {

				ProductName = repo.GetProductName(b.SessionID, int(ProductIndex-1))
				SerialCode = repo.GetSerialCode(b.SessionID, int(ProductIndex-1))
				CategoryGuid := repo.GetCategoryGuidCode(b.SessionID, int(ProductIndex-1))
				CategoryName := repo.GetCategoryName(b.SessionID, int(ProductIndex-1))

				log.Infof("sap/servicerequest.go, Serial Code :- %v", SerialCode)
				log.Infof("sap/servicerequest.go, Product Name :- %v", ProductName)
				log.Infof("sap/servicerequest.go, Category Code :- %v", CategoryGuid)
				log.Infof("sap/servicerequest.go, Category Name :- %v", CategoryName)

				return sapmodel.ComplaintDescription, sapmodel.ErrFallBackComplaint
			}
			respCategories, _ := isCategories(b.SessionID)
			return respCategories, sapmodel.ErrFallBackComplaint

		case sapmodel.ErrFallBackComplaint:
			log.Infof("sap/servicerequest.go, Checking Complaint: %v", b.SessionID)
			ProductIndex := b.DF.Fields["number1"].GetNumberValue()
			ProductCode := b.DF.Fields["ProductCode"].GetStringValue()
			SerialCode := repo.GetSerialCode(b.SessionID, int(ProductIndex-1))
			log.Infof("sap/servicerequest.go, Serial Code : %v", SerialCode)
			log.Infof("sap/servicerequest.go, Product Code : %v", ProductCode)
			if ProductCode == "" {
				ProductCode = SerialCode
			}
			log.Infof("sap/servicerequest.go, Serial Code : %v", SerialCode)
			CategoryGuid := repo.GetCategoryGuidCode(b.SessionID, int(ProductIndex-1))
			CategoryName := repo.GetCategoryName(b.SessionID, int(ProductIndex-1))
			AddressIndex := b.DF.Fields["number"].GetNumberValue()

			Address := repo.GetAddress(b.SessionID, int(AddressIndex-1))
			log.Infof("sap/servicerequest.go, Address : %v", Address)

			IssueType := b.DF.Fields["IssueType"].GetStringValue()
			log.Infof("sap/servicerequest.go, IssueType : %v", IssueType)
			if IssueType == "" {
				IssueType = "Breakdown"
			}
			ComplaintDesc := b.DF.Query
			repo.UpdateComplaintMap(b.SessionID, ComplaintDesc)
			if CategoryGuid != "" {
				return `Please find the details of your raised service request :-

		Product Category :` + CategoryName + `
		
		Issue Type : ` + IssueType + `
			
		Complaint Description : ` + ComplaintDesc + "\n\nPlease check your details to proceed !", sapmodel.ErrFallBackService

			}
			return `Please find the details of your raised service request :-

		Product Code :` + ProductCode + `
					
		Issue Type : ` + IssueType + `
			
		Complaint Description : ` + ComplaintDesc + "\n\nPlease check your details to proceed !", sapmodel.ErrFallBackService

		default:
			return "Dear Customer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", err
		}
	}
	custGUID := repo.GetCustumerGUID(b.SessionID)
	AddressIndex := b.DF.Fields["number"].GetNumberValue()
	AddGUID := repo.GetAddressGUID(b.SessionID, int(AddressIndex-1))
	ComplaintDesc := repo.GetComplaint(b.SessionID)
	IssueType := b.DF.Fields["IssueType"].GetStringValue()
	log.Infof("sap/servicerequest.go, IssueType : %v", IssueType)
	if IssueType == "" {
		IssueType = "Breakdown"
	}

	ProductIndex := b.DF.Fields["number1"].GetNumberValue()
	SerialCode := repo.GetSerialCode(b.SessionID, int(ProductIndex-1))
	log.Infof("sap/servicerequest.go, Serial Code : %v", SerialCode)

	CategoryGuid := repo.GetCategoryGuidCode(b.SessionID, int(ProductIndex-1))

	if CategoryGuid != "" {
		SerialCode = ""
		srResp, _ := RegService(AddGUID, ComplaintDesc, IssueType, SerialCode, CategoryGuid, custGUID)
		return srResp, nil
	}

	return "Dear Customer, Iyris is not able to recognize any product", nil
}

func RegService(addGuid, Complaint, issueType, Code, Name, custGUID string) (string, error) {
	if issueType == "" {
		return "I'm sorry, but we are not able to raise a Service Request for this selected Product ! \nPlease register a product with valid Serial Number or choose from the Other Category option. To proceed you may raise the service request again. Thanks !", sapmodel.ErrEndSessionService
	}
	req := sapmodel.SRrequest{AddressGUID: addGuid, ChiefComplaint: Complaint,
		CustomerGuid: custGUID, NOCName: issueType, SerialNumber: Code, ProductCategoryGuid: Name,
		SourceOfJob: 12}
	res, err := srDetails(&req)
	log.Infof("sap/servicerequest.go, Testing panic : %v, %v", res, err)
	if err != nil {
		return "", err
	}
	if res.Status != "200" {
		return "I'm sorry, we are not able to raise a Service Request for this selected product ! \nPlease try again after sometime.", sapmodel.ErrEndSessionService
	}
	return "Your Service Request ID is : " + res.ID + "\n\nDear Customer, We have registered a service call request for your issue. Our representative will contact you shortly. You can also check the status of your Service Request.", sapmodel.ErrEndSessionService
}

//makeProductIssueResp ::::::::::::::::::::::::::::::::::::::::::::::::
func makeProductIssueResp(serial, key string) {
	req := sapmodel.ProductIssueReq{SrNo: serial}
	res, _ := issueDetails(&req)
	log.Infof("sap/servicerequest.go, Testing panic : %v, %v", res)
	var issCache []repo.IssueCache
	for _, v := range res {
		issCache = append(issCache, repo.IssueCache{Issue: v.Name, Guid: v.GUID})
	}
	repo.UpdateIssueMap(key, issCache)
	return
}

func isProducts(guid, key string) (string, error) {
	req := sapmodel.ProductDetailsReq{CustGUID: guid}
	res, err := productDetails(&req)
	log.Infof("sap/servicerequest.go, Testing panic : %v, %v", res, err)
	if err != nil {
		return "", err
	}
	var ans []string
	var serials []string
	var productnames []string

	var Product string
	var resp string

	for i, v := range res {

		if fmt.Sprintf("%v", v.StatusCode) != "200" {
			respCategories, err := isCategories(key)
			return respCategories, err
		}

		ProductSerialNumber := fmt.Sprintf("%v", v.SerialNumber)
		Product = fmt.Sprintf("%v", v.ProductCategory)
		if Product == "<nil>" {
			Product = "***"
		}
		resp = fmt.Sprintf("%d", i+1) + ". Product : " + Product + "\nSerial Number : " + ProductSerialNumber + "\n\n"
		ans = append(ans, resp)
		serials = append(serials, fmt.Sprintf("%v", v.SerialNumber))
		productnames = append(productnames, fmt.Sprintf("%v", v.ProductCategory))
	}
	Product = "* OTHER PRODUCT CATEGORIES *"
	resp = fmt.Sprintf("%d", 0) + ". Product : " + Product
	ans = append(ans, resp)
	productnames = append(productnames, fmt.Sprintf("%v", Product))

	finalResp := "Here are the list of products :\n\n" + strings.Join(ans, "") + "\n\nDear Customer,\n You may reply with respective OPTION NUMBER to raise the service request or enter  0  to select from other category of products."
	go repo.UpdateProductMap(key, serials)
	go repo.UpdateProductNameMap(key, productnames)
	return finalResp, sapmodel.ErrFallBackProduct
}
func isCategories(key string) (string, error) {

	res, err := categoryDetails()
	if err != nil {
		return "", err
	}

	var ans []string
	var categorynames []string
	var categoryguids []string
	for i, v := range res {
		resp := fmt.Sprintf("%d", i+1) + " . " + v.WaName + "\n\n"
		ans = append(ans, resp)
		categorynames = append(categorynames, fmt.Sprintf("%v", v.WaName))
		categoryguids = append(categoryguids, fmt.Sprintf("%v", v.ID))
	}
	finalResp := "Here is a list of Havells Products for which you would like to raise a Service Request :-\n\n" + strings.Join(ans, "") + "\n\nPlease reply with respective OPTION NUMBER to raise a Service Request"
	go repo.UpdateCategoryMap(key, categorynames)
	go repo.UpdateCategoryGuidMap(key, categoryguids)
	return finalResp, sapmodel.ErrFallBackProduct
}

func isAddressBook(key string) (string, error) {

	custGUID := repo.GetCustumerGUID(key)
	log.Infof("sap/servicerequest.go, Forming address request : %v", custGUID)
	reg := sapmodel.AddressReq{CustomerGUID: custGUID}
	res, err := addressBook(&reg)

	log.Infof("Address Response", fmt.Sprintf("%v", res))
	if err != nil {
		return "", err
	}

	var ans []string
	var addresses []string
	var fulladdress []string
	for i, v := range res {
		if v.StatusCode == "200" {
			FullAddresses := fmt.Sprintf("%v", v.FullAddress)
			resp := fmt.Sprintf("%d", i+1) + ". " + FullAddresses + "\n\n"
			ans = append(ans, resp)
			addresses = append(addresses, fmt.Sprintf("%v", v.AddressGUID))
			fulladdress = append(fulladdress, fmt.Sprintf("%v", v.FullAddress))
		}
		if v.StatusCode == "204" {
			return "Dear Customer, We dont have your Address. So please add your ADDRESS with a VALID PINCODE first and then raise Service Request.", nil
		}
	}
	FullAddresses := "*** ADD NEW ADDRESS ***"
	resp := fmt.Sprintf("%d", (len(ans)+1)) + ". " + FullAddresses
	ans = append(ans, resp)
	fulladdress = append(fulladdress, fmt.Sprintf("%v", FullAddresses))
	finalResp := "Here are the list of registered addresses :\n\n" + strings.Join(ans, "") + "\n\nDear Customer, Please reply with respective OPTION NUMBER where you want our service executive to visit"
	go repo.UpdateAddressMap(key, addresses)
	go repo.UpdateFullAddressMap(key, fulladdress)
	return finalResp, sapmodel.ErrFallBackAddress
}
