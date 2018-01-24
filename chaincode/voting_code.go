package chaincode

import {
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
}

type VotingChainCode struct {}

type Constituency struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Position  string `json:"position"`
}

type Voter struct {
	NIN       string `json:"nin"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	DOB       string `json:"DOB"`
	Constituency string `json:"constituency"`
	Hasvoted  bool   `json:"hasvoted"`
}