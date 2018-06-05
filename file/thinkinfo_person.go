package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"reflect"
	"strconv"
	"strings"
)

type PersonContract struct {
}

type Thinkinfo_Person struct {
	Uid    string `json:"U_id"`
	Name   string `json:"person_name"`
	ID     string `json:"person_id"`
	Age    string `json:"person_age"`
	School string `json:"person_school"`
	Mobile string `json:"person_mobile"`
	Email  string `json:"person_email"`
}

/*************************************
*
*************************************/
func query_person(stub shim.ChaincodeStubInterface, person_data *Thinkinfo_Person, uid string) pb.Response {
	var e_msg string

	json_bytes, err := stub.GetState(uid + "_" + "person")
	if err != nil {
		e_msg = fmt.Sprintf("[ %s ]", err.Error())
		return shim.Error(e_msg)
	}

	if person_data != nil {

		if json_bytes == nil {
			e_msg = fmt.Sprintf("[query_person %s Invalid account !]", uid)
			return shim.Error(e_msg)
		}

		err = json.Unmarshal(json_bytes, person_data)
		if err != nil {
			e_msg = fmt.Sprintf("[ Failed to Unmarshal account json ]")
			return shim.Error(e_msg)
		}

	} else {

		if json_bytes != nil {
			e_msg = fmt.Sprintf("[ %s Exited. ]", uid)
			return shim.Error(e_msg)
		}
	}
	return shim.Success(nil)

}

/*************************************
*
*************************************/
func save_person(stub shim.ChaincodeStubInterface, person_data Thinkinfo_Person, Uid string) pb.Response {
	//	var e_msg string

	json_bytes, err := json.Marshal(person_data)
	if err != nil {
		return shim.Error("[ Failed to Marshal authority_data json ]")
	}

	err = stub.PutState(Uid+"_"+"person", json_bytes)
	if err != nil {
		return shim.Error("[ Failed to save authority_data ]")
	}
	return shim.Success(nil)
}

/*************************************
*
*************************************/
func check_access(stub shim.ChaincodeStubInterface, tabels_name string, account string, password string) pb.Response {
	uid_resp := query_account_index(stub, "Uid", account, password)
	if uid_resp.Status != shim.OK {
		return shim.Error("[ Invalid account or password. ]")
	}

	user_group_resp := query_authority_index(stub, tabels_name, "R_user_auth")
	if user_group_resp.Status != shim.OK {
		return shim.Error("[ Invalid tables_name. ]")
	}

	admin_group_resp := query_authority_index(stub, tabels_name, "R_admin_auth")
	if admin_group_resp.Status != shim.OK {
		return shim.Error("[ Invalid tables_name. ]")
	}

	if false == strings.Contains(string(user_group_resp.Payload), "["+account+"]") && false == strings.Contains(string(admin_group_resp.Payload), "["+account+"]") {
		return shim.Error("[ restricting Authority.]")
	}
	return shim.Success(nil)
}

/*************************************
* checkout ["index", "account", "password"]
*************************************/
func query_account_index(stub shim.ChaincodeStubInterface, index string, account string, pwd string) pb.Response {
	person := [][]byte{[]byte("query_account_index"), []byte(account), []byte(index), []byte(account), []byte(pwd)}
	return stub.InvokeChaincode("account", person, "person1")
}

/*************************************
* checkout ["index", "account", "password"]
*************************************/
func update_records_cnt(stub shim.ChaincodeStubInterface, tables string, account string, pwd string) pb.Response {
	person := [][]byte{[]byte("update_records_cnt"), []byte(tables), []byte(account), []byte(pwd)}
	return stub.InvokeChaincode("authority", person, "person1")
}

/*************************************
* checkout ["index", "account", "password"]
*************************************/
func add_authority_list(stub shim.ChaincodeStubInterface, tables string, flag string, account string, pwd string) pb.Response {
	person := [][]byte{[]byte("add_authority_list"), []byte(tables), []byte(flag), []byte(account), []byte(pwd)}
	return stub.InvokeChaincode("authority", person, "person1")
}

/*************************************
* checkout ["tables_name","index",]
*************************************/
func query_authority_index(stub shim.ChaincodeStubInterface, tables string, index string) pb.Response {
	person := [][]byte{[]byte("query_authority_list_index"), []byte(tables), []byte(index)}
	return stub.InvokeChaincode("authority", person, "person1")
}

/*************************************
*
*
*add_person_list ["name","id","age","school","mobile","eamil","account","password"]
*************************************/
func (t *PersonContract) add_person_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var person Thinkinfo_Person

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := query_account_index(stub, "Uid", args[6], args[7])
	if resp.Status != shim.OK {
		return resp
	}

	person.Uid = string(resp.Payload)
	person.Name = args[0]
	person.ID = args[1]
	person.Age = args[2]
	person.School = args[3]
	person.Mobile = args[4]
	person.Email = args[5]

	resp = query_person(stub, nil, person.Uid)
	if resp.Status != shim.OK {
		return resp
	}

	resp = save_person(stub, person, person.Uid)
	if resp.Status != shim.OK {
		return resp
	}

	resp = add_authority_list(stub, person.Uid, "NG", args[6], args[7])
	if resp.Status != shim.OK {
		return resp
	}

	return shim.Success(nil)
}

/*************************************
*
*
*delete_person ["id" "account", "password" ]
*********************,****************/
func (t *PersonContract) delete_person_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)

}

/*************************************
*
*
*query_person ["tables_name", "account", "password"]
*************************************/
func (t *PersonContract) query_person_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var person Thinkinfo_Person

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := check_access(stub, args[0], args[1], args[2])
	if resp.Status != shim.OK {
		return resp
	}

	resp = query_person(stub, &person, args[0])
	if resp.Status != shim.OK {
		return resp
	}

	json_bytes, err := json.Marshal(person)
	if err != nil {
		return shim.Error("[ Failed to Marshal authority json ]")
	}

	return shim.Success(json_bytes)
}

/*************************************
*
*
*query_index ["id","index","account","password"]
*************************************/
func (t *PersonContract) query_index(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var person Thinkinfo_Person

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := check_access(stub, args[0], args[1], args[2])
	if resp.Status != shim.OK {
		return resp
	}

	resp = query_person(stub, &person, string(resp.Payload))
	if resp.Status != shim.OK {
		return resp
	}

	var index string
	object := reflect.ValueOf(&person)
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
*
*update_index ["id","index","value","account","password"]
*************************************/
func (t *PersonContract) update_index(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var person Thinkinfo_Person

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments.")
	}

	resp := check_access(stub, args[0], args[3], args[4])
	if resp.Status != shim.OK {
		return resp
	}

	resp = query_person(stub, &person, string(resp.Payload))
	if resp.Status != shim.OK {
		return resp
	}

	object := reflect.ValueOf(&person)
	ref := object.Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		if ref.Type().Field(i).Name == args[1] {
			switch field.Kind() {
			case reflect.String:
				ref.FieldByName(args[1]).Set(reflect.ValueOf(args[2]))
				break
			case reflect.Int:
				val_int, err := strconv.Atoi(args[2])
				if err != nil {
					return shim.Error(err.Error())
				}
				ref.FieldByName(args[1]).SetInt(int64(val_int))
				break
			}
		}
	}
	return save_person(stub, person, person.Uid)
}

/*************************************
*
*
*
*************************************/
func (t *PersonContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

/*************************************
*
*
*
*************************************/
func (t *PersonContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "add_person":
		return t.add_person_list(stub, args)
	case "delete_person":
		return t.delete_person_list(stub, args)
	case "query_person":
		return t.query_person_list(stub, args)
	case "query_index":
		return t.query_index(stub, args)
	case "update_index":
		return t.update_index(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query_person\"")
}

/*************************************
*
*
*
*************************************/
func main() {
	err := shim.Start(new(PersonContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
