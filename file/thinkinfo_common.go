package main

import (
	//	"bytes"
	//	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//	"reflect"
	//	"strconv"
	"sync"
)

type CommonContract struct {
}

var U_id uint = 1
var G_id uint = 1
var u_mutex sync.Mutex
var g_mutex sync.Mutex

/*************************************
*
*
*
*************************************/
func set_U_id() uint {
	u_mutex.Lock()
	U_id += 1
	u_mutex.Unlock()
	return U_id
}

/*************************************
*
*
*
*************************************/
func set_G_ip() uint {
	g_mutex.Lock()
	G_id += 1
	g_mutex.Unlock()
	return G_id
}

/*************************************
*
*
*
*************************************/
func (t *CommonContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "Get_Uid":
		return t.get_uid(stub, args)
	case "Get_Gid":
		return t.get_gid(stub, args)
	}
	return shim.Success(nil)
}

/*************************************
*
*
*
*************************************/
func (t *CommonContract) get_uid(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	id_s := fmt.Sprint(set_U_id())
	return shim.Success([]byte(id_s))
}

/*************************************
*
*
*
*************************************/
func (t *CommonContract) get_gid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	id_s := fmt.Sprint(set_G_ip())
	return shim.Success([]byte(id_s))
}

/*************************************
*
*
*
*************************************/
func (t *CommonContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

/*************************************
*
*
*
*************************************/
func main() {
	err := shim.Start(new(CommonContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
