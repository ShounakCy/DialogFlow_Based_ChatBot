package sapmodel

import "github.com/havells/nlp/models"

func HelloModel(phone, session string, df models.DfResp) (string, error) {

	switch df.Intent {
	case "HelloIntent":
		return "", ErrChkReg
	default:
		return "Dear Customer, I am still learning. In case I am not able to assist you, please reach us at customercare@havells.com !", ErrEndRephrase
	}

}
