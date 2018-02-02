package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type VotingChainCode struct{}

var logger = shim.NewLogger("mylogger")

// type Constituency struct {
// 	Name string `json:"name"`
// 	Candidates Candidate `json:"candidates"`
// }

// type Voter struct {
// 	NIN       string `json:"nin"`
// 	Constituency string `json:"constituency"`
// 	Hasvoted  bool   `json:"hasvoted"`
// }

type Candidate struct {
	ObjectType   string `json:"docType"`
	ID           int    `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Party        string `json:"party"`
	Constituency string `json:"constituency"`
	Position     string `json:"position"`
}

type Vote struct {
	ObjectType        string `json:"docType"`
	CandidateID       int    `json:"candidateid"`
	CandidatePosition string `json:"candidateposition"`
	Constituency      string `json:"constituency"`
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

/*Gets candidate details*/
func (t *VotingChainCode) GetIndividualCandidateDetails(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Debug("getting candidate info")

	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return shim.Error("Missing Candidate ID Number")
	}

	var CandidateId = args[0]
	bytes, err := stub.GetState(CandidateId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error("Candidate with ID " + args[0] + " not found")
	}
	fmt.Printf("CandidateDetails: %s", bytes)
	return shim.Success(bytes)
}

/*Gets constituency candidate list*/
func (t *VotingChainCode) GetConstituencyCandidateList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Debug("getting candidate info")

	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return shim.Error("Missing constituency name")
	}

	var constituency = args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"candidate\",\"constituency\":\"%s\"}}", constituency)
	// bytes, err := stub.GetState(constituency)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	// if bytes == nil {
	// 	logger.Error("Constituency not found")
	// 	return nil, fmt.Errorf("Constituency not found: %s", args[0])
	// }
	return shim.Success(queryResults)
}

/*Add a vote to the ledger*/
func (t *VotingChainCode) CastVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Debug("Casting vote")

	if len(args) != 2 {
		logger.Error("Invalid number of args")
		return shim.Error("Expected two arguments to cast vote")
	}

	// var vote Vote
	// voter.NIN = args[0]
	// voter.Firstname = args[1]
	// voter.Lastname = args[2]
	// voter.DOB = args[3]
	// voter.Constituency = args[4]
	// b, fail := strconv.ParseBool(args[5])
	// if fail == nil {
	// 	voter.Hasvoted = b
	// }

	vote := &Vote{"vote", args[0], args[1]}
	voteBytes, f := json.Marshal(vote)
	if f != nil {
		shim.Error(f.Error())
	}

	err := stub.PutState(args[0], voteBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *VotingChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	candidates := []*Candidate{
		&Candidate{"candidate", 1, "John", "Smith", "constituencyA", "president"},
		&Candidate{"candidate", 2, "Jane", "Doe", "constituencyB", "president"},
		&Candidate{"candidate", 3, "Tom", "Sawyer", "constituencyC", "president"}}

	for _, x := range candidates {
		candidateBytes, f := json.Marshal(x)
		if f != nil {
			shim.Error(f.Error())
		}

		err := stub.PutState(strconv.Itoa(x.ID), candidateBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

// func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
// 	logger.Debug("Entering GetCertAttribute")
// 	attr, err := stub.ReadCertAttribute(attributeName)
// 	if err != nil {
// 		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
// 	}
// 	attrString := string(attr)
// 	return attrString, nil
// }

func (t *VotingChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	// var result []byte
	// var err error

	// if fn == "AddVoter" {
	// username, _ := GetCertAttribute(stub, "username")
	// role, _ := GetCertAttribute(stub, "role")
	// if role == "Registrar" {
	// 	return AddVoter(stub, args)
	// } else {
	// 	return nil, errors.New(username + " with role " + role + " does not have access to add a new voter")
	// }
	// 	err = AddVoter(stub, args)
	// 	if err != nil {
	// 		return shim.Error("Failed to add voter")
	// 	}
	// 	return shim.Success([]byte("Voter added"))
	// }

	switch fn {
	case "CastVote":
		return t.CastVote(stub, args)
	case "GetIndividualCandidateDetails":
		return t.GetIndividualCandidateDetails(stub, args)
	case "GetConstituencyCandidateList":
		return t.GetConstituencyCandidateList(stub, args)
	}

	return shim.Error("Unknown function")
}

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(VotingChainCode))
	if err != nil {
		logger.Error("Could not start VotingChaincode")
	} else {
		logger.Info("VotingChaincode successfully started")
	}

}
