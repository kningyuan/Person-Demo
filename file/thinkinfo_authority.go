package main

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type AuthorityContract struct {
}

type Thinkinfo_Authority struct {
	Tables_name       string `json:"tables_name"`
	Owner_uid         string `json:"owner_id"`
	R_group           string `json:"r_group"`
	R_user_auth       string `json:"r_user_auth"`
	R_admin_auth      string `json:"r_admin_auth"`
	Cnt               string `json:"cnt`
	Admin_collocation string `json:"admin_collocation"`
}

var cnt_mutex sync.Mutex

/*************************************
* query_account_index ["index", "account", "password"]
*************************************/
func query_account_index(stub shim.ChaincodeStubInterface, index string, account string, pwd string) pb.Response {
	cmd := [][]byte{[]byte("query_account_index"), []byte(account), []byte(index), []byte(account), []byte(pwd)}
	return stub.InvokeChaincode("account", cmd, "person1")
}

/*************************************
*
*************************************/
func query_authority(stub shim.ChaincodeStubInterface, authority_data *Thinkinfo_Authority, tables_name string) pb.Response {
	var e_msg string

	json_bytes, err := stub.GetState(tables_name + "_" + "authority")
	if err != nil {
		e_msg = fmt.Sprintf("[ %s ]", err.Error())
		return shim.Error(e_msg)
	}

	if authority_data != nil {

		if json_bytes == nil {
			e_msg = fmt.Sprintf("[query authority %s Invalid account ]", tables_name)
			return shim.Error(e_msg)
		}

		err = json.Unmarshal(json_bytes, authority_data)
		if err != nil {
			e_msg = fmt.Sprintf("[ Failed to Unmarshal account json ]")
			return shim.Error(e_msg)
		}

	} else {

		if json_bytes != nil {
			e_msg = fmt.Sprintf("[ %s Exited. ]", tables_name)
			return shim.Error(e_msg)
		}
	}
	return shim.Success(nil)

}

/*************************************
*
*************************************/
func save_authority(stub shim.ChaincodeStubInterface, authority_data Thinkinfo_Authority, tables_name string) pb.Response {
	//	var e_msg string

	json_bytes, err := json.Marshal(authority_data)
	if err != nil {
		return shim.Error("[ Failed to Marshal authority_data json ]")
	}
	err = stub.PutState(tables_name+"_"+"authority", json_bytes)
	if err != nil {
		return shim.Error("[ Failed to save authority_data ]")
	}
	return shim.Success(nil)
}

/*************************************
*
*************************************/
func query_all_authority(stub shim.ChaincodeStubInterface) pb.Response {
	var e_msg string

	json_bytes, err := stub.GetState("Authority")
	if err != nil {
		e_msg = fmt.Sprintf("[ %s ]", err.Error())
		return shim.Error(e_msg)
	}

	return shim.Success(json_bytes)
}

/*************************************
*
*************************************/
func add_all_authority(stub shim.ChaincodeStubInterface, tables_name string) pb.Response {
	resp := query_all_authority(stub)
	if resp.Status != shim.OK {
		return resp
	}

	err := stub.PutState("Authority", []byte(string(resp.Payload)+"["+tables_name+"]"))

	if err != nil {
		return shim.Error("[ Failed to save authority list ]")
	}

	return shim.Success(nil)
}

/*************************************
*
*add_authority_list ["tables","OK",account",password]
*************************************/
func (t *AuthorityContract) add_authority_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var authority Thinkinfo_Authority

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_authority(stub, nil, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	Uid_resp := query_account_index(stub, "Uid", args[2], args[3])
	if Uid_resp.Status != shim.OK {
		return shim.Error("[ invalid account or password ]")
	}

	Gid_resp := query_account_index(stub, "Gid", args[2], args[3])
	if Gid_resp.Status != shim.OK {
		return shim.Error("[ invalid account or password ]")
	}

	authority.Tables_name = args[0]
	authority.Admin_collocation = args[1]
	authority.Owner_uid = string(Uid_resp.Payload)
	authority.R_group = string(Gid_resp.Payload)
	authority.R_admin_auth = ""
	authority.Cnt = "1"

	if "OK" == authority.Admin_collocation {
		authority.R_user_auth = "[" + args[2] + "]" + "[" + "admin" + "]"
	} else {
		authority.R_user_auth = "[" + args[2] + "]"
	}

	resp = add_all_authority(stub, authority.Tables_name)
	if resp.Status != shim.OK {
		return resp
	}

	return save_authority(stub, authority, authority.Tables_name)
}

/*************************************
*
*query_authority_index ["tables,"account","password"]
*************************************/
func (t *AuthorityContract) update_records_cnt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var authority Thinkinfo_Authority

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	Uid_resp := query_account_index(stub, "Uid", args[1], args[2])
	if Uid_resp.Status != shim.OK {
		return shim.Error("[ invalid account or password ]")
	}

	resp := query_authority(stub, &authority, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	if string(Uid_resp.Payload) != authority.Owner_uid {
		return shim.Error("[ restricting Authority ]")
	}

	new_cnt, _ := strconv.Atoi(authority.Cnt)
	cnt_mutex.Lock()
	authority.Cnt = strconv.Itoa(new_cnt + 1)
	cnt_mutex.Unlock()

	return save_authority(stub, authority, authority.Tables_name)
}

/*************************************
*
*query_authority_index ["tables,"index"]
*************************************/
func (t *AuthorityContract) query_authority_list_index(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var authority Thinkinfo_Authority

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_authority(stub, &authority, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	var index string
	object := reflect.ValueOf(&authority)
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
*query_authority_list ["tables"]
*************************************/
func (t *AuthorityContract) query_authority_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var authority Thinkinfo_Authority

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_authority(stub, &authority, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	json_bytes, err := json.Marshal(authority)
	if err != nil {
		return shim.Error("[ Failed to Marshal authority json ]")
	}

	return shim.Success(json_bytes)

}

/*************************************
*
*query_all_authority_list ["Authority"]
*************************************/
func (t *AuthorityContract) query_all_authority_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return query_all_authority(stub)
}

/*************************************
*
*user_authority ["tables","authority_account","flag",account",password]
*************************************/
func (t *AuthorityContract) user_authority(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var authority Thinkinfo_Authority

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments.")
	}

	uid_resp := query_account_index(stub, "Uid", args[3], args[4])
	if uid_resp.Status != shim.OK {
		return uid_resp
	}

	resp := query_authority(stub, &authority, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	if string(uid_resp.Payload) != authority.Owner_uid {
		return shim.Error("[ restricting Authority ]")
	}

	if "authority" == args[2] {
		authority.R_user_auth = authority.R_user_auth + "[" + args[1] + "]"
	} else if "cancel" == args[2] {
		authority.R_user_auth = strings.Replace(authority.R_user_auth, "["+args[1]+"]", "", -1)
	} else {
		return shim.Error("[Incorrect flag.]")
	}

	return save_authority(stub, authority, authority.Tables_name)
}

/*************************************
*
*admin_authority ["tables","authority_account","flag",,"admin",admin]
*************************************/
func (t *AuthorityContract) admin_authority(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var authority Thinkinfo_Authority

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments.")
	}

	uid_resp := query_account_index(stub, "Uid", args[3], args[4])
	if uid_resp.Status != shim.OK {
		return uid_resp
	}

	resp := query_authority(stub, &authority, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	if string(uid_resp.Payload) != "0" || authority.Admin_collocation != "OK" {
		return shim.Error("[ restricting Authority ]")
	}

	if "authority" == args[2] {
		authority.R_admin_auth = authority.R_admin_auth + "[" + args[1] + "]"
	} else if "cancel" == args[2] {
		authority.R_admin_auth = strings.Replace(authority.R_admin_auth, "["+args[1]+"]", "", -1)
	} else {
		return shim.Error("[Incorrect flag.]")
	}

	return save_authority(stub, authority, authority.Tables_name)
}

/*************************************
*
*Admin_collocation ["tables","flag",account","password"]
*************************************/
func (t *AuthorityContract) admin_collocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var authority Thinkinfo_Authority

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}

	uid_resp := query_account_index(stub, "Uid", args[2], args[3])
	if uid_resp.Status != shim.OK {
		return uid_resp
	}

	resp := query_authority(stub, &authority, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	if string(uid_resp.Payload) != authority.Owner_uid {
		return shim.Error("[ restricting Authority. ]")
	}
	if "collocation" == args[1] {
		authority.Admin_collocation = "OK"
		authority.R_user_auth = authority.R_user_auth + "[" + "admin" + "]"
	} else if "cancel" == args[1] {
		authority.Admin_collocation = ""
		authority.R_user_auth = strings.Replace(authority.R_user_auth, "["+"admin"+"]", "", -1)
	} else {
		return shim.Error("[Incorrect flag.]")
	}

	return save_authority(stub, authority, authority.Tables_name)
}

/*************************************
*
*
*
*************************************/
func (t *AuthorityContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

/*************************************
*
*
*
*************************************/
func (t *AuthorityContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "add_authority_list":
		return t.add_authority_list(stub, args)
	case "query_authority_list":
		return t.query_authority_list(stub, args)
	case "query_authority_list_index":
		return t.query_authority_list_index(stub, args)
	case "query_all_authority_list":
		return t.query_all_authority_list(stub, args)
	case "update_records_cnt":
		return t.update_records_cnt(stub, args)
	case "user_authority":
		return t.user_authority(stub, args)
	case "admin_authority":
		return t.admin_authority(stub, args)
	case "admin_collocation":
		return t.admin_collocation(stub, args)
	}
	return shim.Success(nil)
}

/*************************************
*
*
*
*************************************/
func main() {
	err := shim.Start(new(AuthorityContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
