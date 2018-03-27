#!/bin/bash

peer node start --peer-chaincodedev=true -o orderer.ec.ug:7050

sleep 2

peer channel join -b elections.block -o orderer.ec.ug:7050 --tls --cafile /etc/hyperledger/msp/orderer/msp/tlscacerts/tlsca.ec.ug-cert.pem