package main

type BillStruct struct {

	ObjectType string `json:"docType"`
	BillInfoID string `json:"BillInfoID"`
	BillInfoAmt string`json:"BillInfoAmt"`
	BillInfoType string `json:"BillInfoType"`

	BillInfoBgDate string `json:"BillInfoBgDate"` //出票日期
	BillInfoEdDate string `json:"BillInfoEdDate"` //到期时间

	HoldrAcct string `json:"HoldrAcct"` //当前持有人名称
	HoldrCmID string `json:"HoldrCmID"` //当前持有人证件ｉｄ

	WaitEndorseAcct string `json:"WaitEndorseAcct"` //带背书人名称
	WaitEndorseCmID string `json:"WaitEndorseCmID"` //带背书人证件ｉｄ


}