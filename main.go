package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"bytes"
)

type CouchDBChaincode struct {


}

func(t *CouchDBChaincode) Init(stub shim.ChaincodeStubInterface)peer.Response  {
	return shim.Success(nil)
}
func (t *CouchDBChaincode)Invoke(stub shim.ChaincodeStubInterface)peer.Response  {
	fun, args := stub.GetFunctionAndParameters()
	if fun == "billInit" {
		return billInit(stub,args)
	}else if fun == "queryBills"{
		return queryBills(stub,args)
	}else if fun == "queryWaitBills" {
		return queryWaitBills(stub,args)
	}
	return shim.Error("非法操作,指定的fun无效")
}

//初始化票据数据
func billInit(stub shim.ChaincodeStubInterface,args []string)peer.Response  {
	bill := BillStruct{
		ObjectType:"billObj",
		BillInfoID:"POC101",
		BillInfoAmt:"1000",
		BillInfoType:"111",
		BillInfoBgDate:"2018-11-11",
		BillInfoEdDate:"2018-12-12",

		HoldrAcct:"AAA",
		HoldrCmID:"AAAID",
		WaitEndorseAcct:"",
		WaitEndorseCmID:"",
	}
	billByte, _ := json.Marshal(bill)
	err := stub.PutState(bill.BillInfoID,billByte)
	if err != nil {
		return shim.Error("初始化票据1失败"+err.Error())
	}

	bill2 := BillStruct{
		ObjectType:"billObj",
		BillInfoID:"POC102",
		BillInfoAmt:"2000",
		BillInfoType:"111",
		BillInfoBgDate:"2018-09-11",
		BillInfoEdDate:"2018-12-00",

		HoldrAcct:"AAA",
		HoldrCmID:"AAAID",
		WaitEndorseAcct:"BBB",
		WaitEndorseCmID:"BBBID",
	}
	billByte2, _ := json.Marshal(bill2)
	err = stub.PutState(bill2.BillInfoID,billByte2)
	if err != nil {
		return shim.Error("初始化票据2失败"+err.Error())
	}

	bill3 := BillStruct{
		ObjectType:"billObj",
		BillInfoID:"POC103",
		BillInfoAmt:"3000",
		BillInfoType:"111",
		BillInfoBgDate:"2018-01-11",
		BillInfoEdDate:"2018-12-10",

		HoldrAcct:"BBB",
		HoldrCmID:"BBBID",
		WaitEndorseAcct:"CCC",
		WaitEndorseCmID:"CCCID",
	}
	billByte3, _ := json.Marshal(bill3)
	err = stub.PutState(bill3.BillInfoID,billByte3)
	if err != nil {
		return shim.Error("初始化票据2失败"+err.Error())
	}

	bill4 := BillStruct{
		ObjectType:"billObj",
		BillInfoID:"POC104",
		BillInfoAmt:"4000",
		BillInfoType:"111",
		BillInfoBgDate:"2018-10-11",
		BillInfoEdDate:"2018-11-12",

		HoldrAcct:"CCC",
		HoldrCmID:"CCCID",
		WaitEndorseAcct:"BBB",
		WaitEndorseCmID:"BBBID",
	}
	billByte4, _ := json.Marshal(bill4)
	err =stub.PutState(bill4.BillInfoID,billByte4)
	if err != nil {
		return shim.Error("初始化票据2失败"+err.Error())
	}

	return shim.Success([]byte("票据初始化成功"))
}
//批量查询持票人票据
func queryBills(stub shim.ChaincodeStubInterface,args []string)peer.Response  {
	if len(args) != 1 {
		return shim.Error("必须指定持票人证件号码")
	}
	HoldrCmID := args[0]
	//拼接couchDB所需的查询字符串=》json串
	//”{\"key\":{\"key\":\"value\",\"k\":\"v\"}}“
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"billObj\", \"HoldrCmID\":\"%s\"}}",HoldrCmID)

	result, err := getBillsByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据持票人的证件号码批量查询持票人票据:"+err.Error())
	}
	return shim.Success(result)
}
//根据待背书人的证件号码批量查询待背书列表
func queryWaitBills(stub shim.ChaincodeStubInterface,args []string)peer.Response  {
	if len(args) != 1 {
		return shim.Error("必须指定持票人证件号码")
	}
	WaitEndorseCmID := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"billObj\", \"WaitEndorseCmID\":\"%s\"}}",WaitEndorseCmID)
	result, err := getBillsByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据待背书人的证件号码批量查询待背书列表:"+err.Error())
	}
	return shim.Success(result)
}
//查询数据
func getBillsByQueryString(stub shim.ChaincodeStubInterface,queryString string) ([]byte, error) {

	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil,err
	}
	defer iterator.Close()

	var buffer bytes.Buffer
	var isSpilt bool

	for iterator.HasNext() {
		result, err := iterator.Next()
		if err != nil {
			return nil,err
		}
		if isSpilt {
			buffer.WriteString("; ")
		}

		buffer.WriteString("key: ")
		buffer.WriteString(result.Key)
		buffer.WriteString(", value: ")
		buffer.WriteString(string(result.Value))

		isSpilt = true

	}
	return buffer.Bytes(),nil
}

func main() {
	err := shim.Start(new(CouchDBChaincode))
	if err != nil {
		fmt.Errorf("启动链码失败：%v",err)
	}
}