package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"reflect"
	"strconv"
	//	"sync"
)

type AccountContract struct {
}

//map["account"]Thinkinfo_Account
type Thinkinfo_Account struct {
	Account string `json:"account"`
	Pwd     string `json:"password"`
	Uid     string `json:"uid"`
	Gid     string `json:"gip"`
	ID      string `json:"id"` /*key to person name*/
}

/*************************************
*
*************************************/
func make_byte_buffer(args []string) []byte {

	var buffer bytes.Buffer
	buffer.WriteString("{")
	for _, arg := range args {
		buffer.WriteString("[")
		buffer.WriteString(arg)
		buffer.WriteString("]")
	}
	buffer.WriteString("}")
	return buffer.Bytes()
}

/*************************************
*
*************************************/
func get_uid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	uid := [][]byte{[]byte("Get_Uid"), []byte(nil)}
	return stub.InvokeChaincode("common", uid, "person1")
}

/*************************************
*
*************************************/
func get_gid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	gid := [][]byte{[]byte("Get_Gid"), []byte(nil)}
	return stub.InvokeChaincode("common", gid, "person1")
}

/*************************************
*
*************************************/
func query_account(stub shim.ChaincodeStubInterface, accout_data *Thinkinfo_Account, account_name string) pb.Response {
	var e_msg string

	json_bytes, err := stub.GetState(account_name)
	if err != nil {
		e_msg = fmt.Sprintf("[ %s ]", err.Error())
		return shim.Error(e_msg)
	}

	if accout_data != nil {

		if json_bytes == nil {
			e_msg = fmt.Sprintf("[ %s Invalid account ]", account_name)
			return shim.Error(e_msg)
		}

		err = json.Unmarshal(json_bytes, accout_data)
		if err != nil {
			e_msg = fmt.Sprintf("[ Failed to Unmarshal account json ]")
			return shim.Error(e_msg)
		}

	} else {

		if json_bytes != nil {
			e_msg = fmt.Sprintf("[ %s Exited. ]", account_name)
			return shim.Error(e_msg)
		}
	}
	return shim.Success(nil)

}

/*************************************
*
*************************************/
func save_account(stub shim.ChaincodeStubInterface, accout_data Thinkinfo_Account, account_name string) pb.Response {
	//	var e_msg string

	json_bytes, err := json.Marshal(accout_data)
	if err != nil {
		return shim.Error("[ Failed to Marshal accout_data json ]")
	}
	err = stub.PutState(account_name, json_bytes)
	if err != nil {
		return shim.Error("[ Failed to save accout_data ]")
	}
	return shim.Success(nil)
}

/*************************************
*
*register_account [account password id]
*************************************/
func (t *AccountContract) register_account(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var e_msg string
	var account Thinkinfo_Account

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_account(stub, nil, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	resp_uid := get_uid(stub, nil)
	if resp_uid.Status != shim.OK {
		e_msg = fmt.Sprintf("[ get_uid error. %s ]", resp_uid.Payload)
		return shim.Error(e_msg)
	}

	resp_gid := get_gid(stub, nil)
	if resp_gid.Status != shim.OK {
		e_msg = fmt.Sprintf("[ get_gid error. %s ]", resp_gid.Payload)
		return shim.Error(e_msg)
	}

	account.Account = args[0]
	account.Pwd = args[1]
	account.Uid = string(resp_uid.Payload)
	account.Gid = string(resp_gid.Payload)
	account.ID = args[2]

	return save_account(stub, account, args[0])

}

/*************************************
*
*reset_password ["account", "password", "id"]
*************************************/
func (t *AccountContract) reset_password(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var e_msg string
	var account Thinkinfo_Account

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_account(stub, &account, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	if args[2] != account.ID {
		e_msg = fmt.Sprintf("[ Invalid ID. %s]", args[2])
		return shim.Error(e_msg)
	}

	account.Pwd = args[1]

	return save_account(stub, account, args[0])
}

/*************************************
*
*query_account ["account", "account", "password"]
*************************************/
func (t *AccountContract) query_account(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	var e_msg string
	var account Thinkinfo_Account

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_account(stub, &account, args[1])
	if resp.Status != shim.OK {
		return resp
	}

	if account.Pwd != args[2] {
		return shim.Error("[ Invalid password. ]")
	}

	if account.Gid != "0" && args[1] != args[0] {
		return shim.Error("[ Restricting access. ]")
	}

	resp = query_account(stub, &account, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	json_bytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("[ Failed to Marshal account json ]")
	}

	return shim.Success(json_bytes)
	/*
		object := reflect.ValueOf(&account)
		ref := object.Elem()
		var arg []string

		for i := 0; i < ref.NumField(); i++ {
			field := ref.Field(i)
			switch field.Kind() {
			case reflect.String:
				arg = append(arg, field.Interface().(string))
			case reflect.Int:
				arg = append(arg, strconv.Itoa(field.Interface().(int)))
			}
		}

		return shim.Success(make_byte_buffer(arg))
	*/
}

/*************************************
*
*query_account_index ["account,"index","account","password"]
*************************************/
func (t *AccountContract) query_account_index(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var account Thinkinfo_Account

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_account(stub, &account, args[2])
	if resp.Status != shim.OK {
		return resp
	}

	if account.Pwd != args[3] {
		return shim.Error("[ Invalid password. ]")
	}

	if account.Gid != "0" && args[2] != args[0] {
		return shim.Error("[ Restricting access. ]")
	}

	resp = query_account(stub, &account, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	var index string
	object := reflect.ValueOf(&account)
	ref := object.Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		if ref.Type().Field(i).Name == args[1] {
			switch field.Kind() {
			case reflect.String:
				index = field.Interface().(string)
			case reflect.Int:
				index = strconv.Itoa(field.Interface().(int))
			}
		}
	}
	return shim.Success([]byte(index))
}

/*************************************
*
*************************************/
func (t *AccountContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	admin := Thinkinfo_Account{"admin", "admin", "0", "0", "0"}

	json_admin, err := json.Marshal(admin)
	if err != nil {
		return shim.Error("[ Failed to Marshal account json. ]")
	}

	err = stub.PutState("admin", json_admin)
	if err != nil {
		return shim.Error("[ Failed to register admin. ]")
	}
	return shim.Success(nil)

}

/*************************************
*
*************************************/
func (t *AccountContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "register_account":
		return t.register_account(stub, args)
	case "reset_password":
		return t.reset_password(stub, args)
	case "query_account":
		return t.query_account(stub, args)
	case "query_account_index":
		return t.query_account_index(stub, args)
	}
	return shim.Success(nil)
}

/*************************************
*
*************************************/
func main() {
	err := shim.Start(new(AccountContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
