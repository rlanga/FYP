package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ElectoralRegisterChainCode struct{}

var logger = shim.NewLogger("mylogger")

type Voter struct {
	NIN          string `json:"nin"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	DOB          string `json:"DOB"`
	Constituency string `json:"constituency"`
	Hasvoted     bool   `json:"hasvoted"`
}

type Candidate struct {
	ID           string `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Constituency string `json:"constituency"`
	Position     string `json:"position"`
}

type Constituency struct {
	Name      string      `json:"name"`
	Candidate []Candidate `json:"candidates"`
}

/*Add a new voter to the ledger*/
func AddVoter(stub shim.ChaincodeStubInterface, args []string) error {
	logger.Debug("Adding voter details")

	if len(args) < 6 {
		logger.Error("Invalid number of args")
		return errors.New("Expected six arguments for voter addition")
	}

	var voter Voter
	voter.NIN = args[0]
	voter.Firstname = args[1]
	voter.Lastname = args[2]
	voter.DOB = args[3]
	voter.Constituency = args[4]
	b, fail := strconv.ParseBool(args[5])
	if fail == nil {
		voter.Hasvoted = b
	}

	// var voterId = args[0]
	// var voterDetails = args[1]

	// var voterDetails = []string{voter.Firstname, voter.Lastname, voter.DOB, voter.Constituency, args[5]}
	// var voterDetails = [voter.Firstname, voter.Lastname, voter.DOB, voter.Constituency, voter.Hasvoted]
	vtBytes, f := json.Marshal(&voter)
	if f != nil {
		logger.Error("Error marshalling voter details")
	}

	err := stub.PutState(voter.NIN, vtBytes)
	if err != nil {
		logger.Error("Could not save voter details to ledger", err)
		return err
	}

	// var customEvent = "{eventType: 'loanApplicationCreation', description:" + loanApplicationId + "' Successfully created'}"
	// err = stub.SetEvent("evtSender", []byte(customEvent))
	// if err != nil {
	// 	return nil, err
	// }
	logger.Info("Successfully added Voter")
	return err

}

/*Gets voter details*/
func GetVoterDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("getting voter info")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing National ID Number")
	}

	var NationalIdNo = args[0]
	bytes, err := stub.GetState(NationalIdNo)
	if err != nil {
		logger.Error("Could not fetch voter details with id "+NationalIdNo+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

/*Gets candidate details*/
func GetIndividualCandidateDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("getting candidate info")

	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing Candidate ID Number")
	}

	var CandidateId = args[0]
	bytes, err := stub.GetState(CandidateId)
	if err != nil {
		logger.Error("Could not fetch candidate details with id "+CandidateId+" from ledger", err)
		return nil, err
	}
	if bytes == nil {
		logger.Error("Candidate not found")
		return nil, fmt.Errorf("Candidate not found: %s", args[0])
	}
	return bytes, nil
}

/*Gets constituency candidate list*/
func GetConstituencyCandidateList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("getting candidate info")

	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing constituency name")
	}

	var constituency = args[0]
	bytes, err := stub.GetState(constituency)
	if err != nil {
		logger.Error("Could not load "+constituency+" candidate list from ledger", err)
		return nil, err
	}
	if bytes == nil {
		logger.Error("Constituency not found")
		return nil, fmt.Errorf("Constituency not found: %s", args[0])
	}
	return bytes, nil
}

func (t *ElectoralRegisterChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *ElectoralRegisterChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// if function == "GetLoanApplication" {
	// 	return GetLoanApplication(stub, args)
	// }

	switch function {
	case "GetVoterDetails":
		return GetVoterDetails(stub, args)
	case "GetCandidateDetails":
		return GetCandidateDetails(stub, args)
	}
	return nil, errors.New("unknown function")

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

func (t *ElectoralRegisterChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	// var result []byte
	var err error

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
	case "AddVoter":
		err = AddVoter(stub, args)
		if err != nil {
			return shim.Error("Failed to add voter")
		}
		return shim.Success([]byte("Voter added"))

	case "GetVoterDetails":
		result, err := GetVoterDetails(stub, args)
		if err != nil {
			return shim.Error("Could not get voter details")
		}
		return shim.Success(result)
	case "GetCandidateDetails":
		result, err := GetCandidateDetails(stub, args)
		if err != nil {
			return shim.Error("Could not get candidate details")
		}
		return shim.Success(result)
		// case "GetConstituencyCandidateDetails":
	}

	return shim.Error("Unknown function")
}

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(ElectoralRegisterChainCode))
	if err != nil {
		logger.Error("Could not start VoterChaincode")
	} else {
		logger.Info("VoterChaincode successfully started")
	}

}
