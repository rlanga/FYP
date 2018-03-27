package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type TallyingChainCode struct{}

type Vote struct {
	ObjectType        string `json:"docType"`
	CandidateID       int    `json:"candidateid"`
	CandidatePosition string `json:"candidateposition"`
	Constituency      string `json:"constituency"`
}

type Candidate struct {
	ObjectType   string `json:"docType"`
	ID           int    `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Party        string `json:"party"`
	Constituency string `json:"constituency"`
	Position     string `json:"position"`
}

type CandidateTally struct {
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Position   string `json:"position"`
	Totalvotes int    `json:"totalvotes"`
}

var logger = shim.NewLogger("tallyinglogger")

/*show votes sorted per candidate in each constituency*/
func (t *TallyingChainCode) CountConstituencyCandidateVotes(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Debug("counting candidate votes in constituency")

	if len(args) != 1 {
		logger.Error("Invalid number of arguments")
		return shim.Error("Missing constituency name")
	}
	var constituency = args[0]
	var votes []Vote
	var candidates []Candidate
	votesQueryString := fmt.Sprintf("{\"selector\":{\"data.docType\":\"vote\",\"data.constituency\":\"%s\"}}", constituency)
	candidatesQueryString := fmt.Sprintf("{\"selector\":{\"data.docType\":\"candidate\",\"data.constituency\":\"%s\"}}", constituency)

	voteResultsIterator, err := stub.GetQueryResult(votesQueryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer voteResultsIterator.Close()
	for voteResultsIterator.HasNext() {
		votesQueryResponse, er := voteResultsIterator.Next()
		if er != nil {
			return shim.Error(er.Error())
		}
		var v Vote
		voteEr := json.Unmarshal(votesQueryResponse.Value, v)
		if voteEr != nil {
			return shim.Error(voteEr.Error())
		}
		logger.Info(v)
		votes = append(votes, v)
	}

	candidatesIterator, e := stub.GetQueryResult(candidatesQueryString)
	if e != nil {
		return shim.Error(e.Error())
	}
	defer candidatesIterator.Close()
	for candidatesIterator.HasNext() {
		candidatesQueryResponse, errr := candidatesIterator.Next()
		if errr != nil {
			return shim.Error(errr.Error())
		}
		var c Candidate
		candEr := json.Unmarshal(candidatesQueryResponse.Value, c)
		if candEr != nil {
			return shim.Error(candEr.Error())
		}
		logger.Info(c)
		candidates = append(candidates, c)
	}

	var tallies []*CandidateTally
	for _, cdt := range candidates {
		var totalVotes int
		for _, vt := range votes {
			if vt.CandidateID == cdt.ID {
				totalVotes++
			}
		}
		tallies = append(tallies, &CandidateTally{cdt.Firstname, cdt.Lastname, cdt.Position, totalVotes})
	}

	results, i := json.Marshal(tallies)
	if i != nil {
		shim.Error(i.Error())
	}

	return shim.Success(results)
}

/*show total number of votes for presidential candidates nationwide*/
func (t *TallyingChainCode) CountPresidentialCandidateVotes(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Debug("counting presidential candidate votes")

	var votes []Vote
	var candidates []Candidate
	votesQueryString := "{\"selector\":{\"docType\":\"vote\",\"candidateposition\":\"president\"}}"
	candidatesQueryString := "{\"selector\":{\"docType\":\"candidate\",\"position\":\"president\"}}"

	voteResultsIterator, err := stub.GetQueryResult(votesQueryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer voteResultsIterator.Close()
	for voteResultsIterator.HasNext() {
		votesQueryResponse, er := voteResultsIterator.Next()
		if er != nil {
			return shim.Error(er.Error())
		}
		var v Vote
		fmt.Println("response value: " + string(votesQueryResponse.Value))
		voteEr := json.Unmarshal(votesQueryResponse.Value, &v)
		if voteEr != nil {
			return shim.Error(voteEr.Error())
		}
		logger.Info(v)
		votes = append(votes, v)
	}

	candidatesIterator, e := stub.GetQueryResult(candidatesQueryString)
	if e != nil {
		return shim.Error(e.Error())
	}
	defer candidatesIterator.Close()
	for candidatesIterator.HasNext() {
		candidatesQueryResponse, errr := candidatesIterator.Next()
		if errr != nil {
			return shim.Error(errr.Error())
		}
		var c Candidate
		fmt.Println("response value: " + string(candidatesQueryResponse.Value))
		candEr := json.Unmarshal(candidatesQueryResponse.Value, &c)
		if candEr != nil {
			return shim.Error(candEr.Error())
		}
		logger.Info(c)
		candidates = append(candidates, c)
	}

	var tallies []*CandidateTally
	for _, cdt := range candidates {
		var totalVotes int
		for _, vt := range votes {
			if vt.CandidateID == cdt.ID {
				totalVotes++
			}
		}
		tallies = append(tallies, &CandidateTally{cdt.Firstname, cdt.Lastname, cdt.Position, totalVotes})
	}

	results, i := json.Marshal(tallies)
	if i != nil {
		shim.Error(i.Error())
	}

	return shim.Success(results)
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
	// return shim.Success(buffer.Bytes())
}

func (t *TallyingChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {

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

func (t *TallyingChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
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
	case "CountConstituencyCandidateVotes":
		return t.CountConstituencyCandidateVotes(stub, args)
	case "CountPresidentialCandidateVotes":
		return t.CountPresidentialCandidateVotes(stub)

	}

	return shim.Error("Unknown function")
}

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	shim.SetLoggingLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(TallyingChainCode))
	if err != nil {
		logger.Error("Could not start TallyingChainCode" + err.Error())
	} else {
		logger.Info("TallyingChainCode successfully started")
	}

}
