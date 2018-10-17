package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type DeliveredTo struct {
	Quantity  string `json:quantity`
	NewOwner  string `json:quantity`
	RequestId string `json:reqeustId`
	JsonObj   string `json:jsonObj`
}
type ApproveReject struct {
	Decision []Request `json:decision`
}
type Request struct {
	RequestedFrom string `json: requestedfrom`
	RequestId     string `json: requestid`
	JsonObj       string `json: jsonObj`
	Id            string `json:id`
	Quantity      string `json:quantity`
	Requester     string `json:requester`
	Status        string `json:status`
}
type Transaction struct {
	Owner        string        `json:owner`
	Id           string        `json:id`
	JsonObj      string        `json:jsonObj`
	Quantity     string        `json:quantity`
	Transferrer  string        `json:transferrer`
	Receiver     string        `json:receiver`
	Dispatchedto []DeliveredTo `json: dispatchedto`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createProduct" { //create a new marble
		return t.createProduct(stub, args)
	} else if function == "transferProduct" { //change owner of a specific marble
		return t.transferProduct(stub, args)
	} else if function == "readProduct" { //read a marble
		return t.readProduct(stub, args)
	} else if function == "getHistoryForProduct" { //get history of values for a marble
		return t.getHistoryForProduct(stub, args)
	} else if function == "requestProduct" { //get history of values for a marble
		return t.requestProduct(stub, args)
	} else if function == "getproductByRange" { //get history of values for a marble
		return t.getproductByRange(stub, args)
	} else if function == "queryonCompositeKey" {
		return t.queryonCompositeKey(stub, args)
	} else if function == "deleteRequest" {
		return t.deleteRequest(stub, args)

	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initMarble - create a new product, store into chaincode state
// ============================================================
func (t *SimpleChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//   0       1       2
	// "id", "jsonobj", "quantity"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	Id := args[0]
	Jsonobj := args[1]
	quantity := args[2]

	// ==== Check if product already exists ====
	ProductAsBytes, err := stub.GetState(Id)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if ProductAsBytes != nil {
		fmt.Println("This product with the product-ID " + Id + "already exists: ")
		return shim.Error("This product with the product-ID " + Id + "already exists: ")
	}

	// ==== Create marble object and marshal to JSON ====

	Product := &Transaction{Owner: "Manufacturer", Id: Id, JsonObj: Jsonobj, Quantity: quantity, Transferrer: "Manufacturer", Receiver: "Distributor"}

	ProductJSONasBytes, err := json.Marshal(Product)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save marble to state ===
	err = stub.PutState(Id, ProductJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "transferrer~receiver"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{Product.Transferrer, Product.Receiver})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(colorNameIndexKey, value)

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init marble")
	messageAsbytes, _ := json.Marshal("product created successfully")
	return shim.Success(messageAsbytes)
}

// ===============================================
// readMarble - read a marble from chaincode state
// ===============================================
func (t *SimpleChaincode) readProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var id, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + id + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===========================================================
// transfer a marble by setting a new owner name on the marble
// ===========================================================
func (t *SimpleChaincode) transferProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0      1        2          3                     4           5         6
	// "id", "jsonobj","Quantity" "decision"   "requestid"    "requestedFrom"  "newOwner"
	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	var requester string
	id := args[0]
	jsonobj := args[1]
	quantity := args[2]
	decision := args[3]
	requestid := args[4]
	RequestedFrom := args[5]
	newOwner := args[6]
	if decision == "approved" {
		RequestAsbytes, err := stub.GetState(RequestedFrom)
		if err != nil {
			return shim.Error("Failed to get requested proposal:" + err.Error())
		}
		requestArray := ApproveReject{}
		err = json.Unmarshal(RequestAsbytes, &requestArray) //unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}
		for i := 0; i < len(requestArray.Decision); i++ {

			if requestArray.Decision[i].RequestId == requestid {
				requester = requestArray.Decision[i].Requester
				requestArray.Decision[i].Status = "Approved"
				break
			}
		}
		requestasBytes, _ := json.Marshal(requestArray)
		err = stub.PutState(RequestedFrom, requestasBytes) //rewrite the marble
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("- changed the status of request (success)")

		fmt.Println("- starting transfer of products ", id, newOwner, quantity)

		productAsBytes, err := stub.GetState(id)
		if err != nil {
			return shim.Error("Failed to get marble:" + err.Error())
		} else if productAsBytes == nil {
			return shim.Error("Marble does not exist")
		}

		productToTransfer := Transaction{}
		err = json.Unmarshal(productAsBytes, &productToTransfer) //unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}
		if newOwner == "distributor" {
			deliveredto := DeliveredTo{quantity, newOwner, requestid, jsonobj}
			//change the owner
			intQuantity, err := strconv.Atoi(productToTransfer.Quantity)
			if err != nil {
				return shim.Error(err.Error())
			}
			requestedQuantity, err := strconv.Atoi(quantity)
			newQuantity := intQuantity - requestedQuantity
			if newQuantity < 0 {
				return shim.Error("reduce the quantity only" + productToTransfer.Quantity + "left")
			}
			stringQuantity := strconv.Itoa(newQuantity)
			productToTransfer.Transferrer = newOwner
			productToTransfer.Receiver = requester
			productToTransfer.Quantity = stringQuantity
			productToTransfer.JsonObj = jsonobj
			productToTransfer.Dispatchedto = append(productToTransfer.Dispatchedto, deliveredto)
			//	productToTransfer.Owner = newOwner

		} else if newOwner == "retailer" {
			deliveredto := DeliveredTo{quantity, newOwner, requestid, jsonobj}
			//change the owner
			intQuantity, err := strconv.Atoi(productToTransfer.Dispatchedto[0].Quantity)
			if err != nil {
				return shim.Error(err.Error())
			}
			requestedQuantity, err := strconv.Atoi(quantity)
			newQuantity := intQuantity - requestedQuantity
			if newQuantity < 0 {
				return shim.Error("reduce the quantity only" + productToTransfer.Quantity + "left")
			}
			stringQuantity := strconv.Itoa(newQuantity)
			productToTransfer.Transferrer = newOwner
			productToTransfer.Receiver = requester
			productToTransfer.Quantity = stringQuantity
			productToTransfer.JsonObj = jsonobj
			productToTransfer.Dispatchedto = append(productToTransfer.Dispatchedto, deliveredto)
			//	productToTransfer.Owner = newOwner
		} else if newOwner == "enduser" {
			deliveredto := DeliveredTo{quantity, newOwner, requestid, jsonobj}
			//change the owner
			intQuantity, err := strconv.Atoi(productToTransfer.Dispatchedto[1].Quantity)
			if err != nil {
				return shim.Error(err.Error())
			}
			requestedQuantity, err := strconv.Atoi(quantity)
			newQuantity := intQuantity - requestedQuantity
			if newQuantity < 0 {
				return shim.Error("reduce the quantity only" + productToTransfer.Quantity + "left")
			}
			stringQuantity := strconv.Itoa(newQuantity)
			productToTransfer.Transferrer = newOwner
			productToTransfer.Receiver = requester
			productToTransfer.Quantity = stringQuantity
			productToTransfer.JsonObj = jsonobj
			productToTransfer.Dispatchedto = append(productToTransfer.Dispatchedto, deliveredto)
			//	productToTransfer.Owner = newOwner
		}
		productJSONasBytes, _ := json.Marshal(productToTransfer)
		err = stub.PutState(id, productJSONasBytes) //rewrite the marble
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("- end transferMarble (success)")
		return shim.Success(nil)

	} else if decision == "rejected" {
		RequestAsbytes, err := stub.GetState(RequestedFrom)
		if err != nil {
			return shim.Error("Failed to get requested proposal:" + err.Error())
		}
		requestArray := ApproveReject{}
		err = json.Unmarshal(RequestAsbytes, &requestArray.Decision) //unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}
		for i := 0; i < len(requestArray.Decision); i++ {

			if requestArray.Decision[i].RequestId == requestid {
				requestArray.Decision[i].Status = decision
				break
			}
		}
		requestasBytes, _ := json.Marshal(requestArray)
		err = stub.PutState(RequestedFrom, requestasBytes) //rewrite the marble
		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("- changed the status of request (success)")
	} else {
		fmt.Println("not a valid option")
		return shim.Success(nil)
	}
	return shim.Success(nil)
}

// ===========================================================================================
// getMarblesByRange performs a range query based on the start and end keys provided.
// Read-only function results are not typically submitted to ordering. If the read-only
// results are submitted to ordering, or if the query is used in an update transaction
// and submitted to ordering, then the committing peers will re-execute to guarantee that
// result sets are stable between endorsement time and commit time. The transaction is
// invalidated by the committing peers if the result set has changed between endorsement
// time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (t *SimpleChaincode) getproductByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//==============get history for a particular product===============================================//
func (t *SimpleChaincode) getHistoryForProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	marbleName := args[0]

	fmt.Printf("- start getHistoryForMarble: %s\n", marbleName)

	resultsIterator, err := stub.GetHistoryForKey(marbleName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//=================request product==============================================//
func (t *SimpleChaincode) requestProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var requestFrom, id, jsonResp, quantity, requestid, requester string
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}
	requestFrom = args[0]
	id = args[1]
	quantity = args[2]
	requestid = args[3]
	requester = args[4]
	valAsbytes, err := stub.GetState(id) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"product with id: " + id + " does not exist:\"}"
		return shim.Error(jsonResp)

	}
	productObject := Transaction{}
	err = json.Unmarshal(valAsbytes, &productObject)

	ArrayAsBytes, err := stub.GetState(requestFrom)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	}
	requestArray := ApproveReject{}
	err = json.Unmarshal(ArrayAsBytes, &requestArray)

	request := Request{RequestedFrom: requestFrom, JsonObj: productObject.JsonObj, Id: id, RequestId: requestid, Quantity: quantity, Requester: requester, Status: "initiated"}
	requestArray.Decision = append(requestArray.Decision, request)

	requestasBytes, err := json.Marshal(requestArray)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save marble to state ===//
	err = stub.PutState(requestFrom, requestasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//==================partial composite keys=============================================================//
func (t *SimpleChaincode) queryonCompositeKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "transferer", "receiver"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	transferer := args[0]
	receiver := args[1]
	fmt.Println("- query product based on transferer and receiver ", transferer, receiver)

	// Query the color~name index by color
	// This will execute a key range query on all keys starting with 'color'
	coloredMarbleResultsIterator, err := stub.GetStateByPartialCompositeKey("transferer~receiver", []string{transferer})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer coloredMarbleResultsIterator.Close()

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int
	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := coloredMarbleResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedColor := compositeKeyParts[0]
		returnedMarbleName := compositeKeyParts[1]
		returnedColorAsbytes, _ := json.Marshal(returnedColor)
		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, returnedColor, returnedMarbleName)
		return shim.Success(returnedColorAsbytes)
	}
	return shim.Success(nil)

}

//=================================delete request============================================================//
func (t *SimpleChaincode) deleteRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//id            //requestid      //userType
	if len(args) < 3 {
		return shim.Error("invalid number of args")
	}
	id := args[0]
	Requestid := args[1]
	userType := args[2]

	RequestAsbytes, err := stub.GetState(userType)

	if err != nil {
		return shim.Error("Failed to get requests for : " + userType + " " + err.Error())
	}
	requestArray := ApproveReject{}
	err = json.Unmarshal(RequestAsbytes, &requestArray) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	for i := 0; i < len(requestArray.Decision); i++ {

		if requestArray.Decision[i].RequestId == Requestid && requestArray.Decision[i].Id == id {

			requestArray.Decision[i] = requestArray.Decision[len(requestArray.Decision)-1] // Copy last element to index i
			requestArray.Decision[len(requestArray.Decision)-1] = Request{}                // Erase last element (write zero value)
			requestArray.Decision = requestArray.Decision[:len(requestArray.Decision)-1]
			break
		}
	}
	requestasBytes, _ := json.Marshal(requestArray)
	err = stub.PutState(userType, requestasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}
