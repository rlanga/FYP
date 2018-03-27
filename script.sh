#!/bin/bash
# Copyright London Stock Exchange Group All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
set -e
# This script expedites the chaincode development process by automating the
# requisite channel create/join commands

#wait for a bit while other containers set themselves up
sleep 10
# We use a pre-generated orderer.block and channel transaction artifact (myc.tx),
# both of which are created using the configtxgen tool
# ls /etc/hyperledger/msp/peers/peer0.BallotMachines.ec.ug/tls/
# first we create the channel against the specified configuration in myc.tx
# this call returns a channel configuration block - myc.block - to the CLI container
peer channel create -c elections -f channel-artifacts/elections.tx -o orderer.ec.ug:7050 --tls --cafile /etc/hyperledger/msp/orderer/msp/tlscacerts/tlsca.ec.ug-cert.pem

# now we will join the channel and start the chain with myc.block serving as the
# channel's first block (i.e. the genesis block)
peer channel join -b elections.block
CORE_PEER_ADDRESS=peer1.BallotMachines.ec.ug:7051 CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/msp/peers/peer1.BallotMachines.ec.ug/tls/ca.crt peer channel join -b elections.block
CORE_PEER_ADDRESS=peer2.BallotMachines.ec.ug:7051 CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/msp/peers/peer2.BallotMachines.ec.ug/tls/ca.crt peer channel join -b elections.block
CORE_PEER_ADDRESS=peer3.BallotMachines.ec.ug:7051 CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/msp/peers/peer3.BallotMachines.ec.ug/tls/ca.crt peer channel join -b elections.block
CORE_PEER_ADDRESS=peer4.BallotMachines.ec.ug:7051 CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/msp/peers/peer4.BallotMachines.ec.ug/tls/ca.crt peer channel join -b elections.block

CORE_PEER_ADDRESS=peer0.BallotMachines.ec.ug:7051
CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/msp/peers/peer0.BallotMachines.ec.ug/tls/ca.crt
peer channel update -c elections -f channel-artifacts/BallotMachinesMSPanchors.tx -o orderer.ec.ug:7050 --tls --cafile /etc/hyperledger/msp/orderer/msp/tlscacerts/tlsca.ec.ug-cert.pem

# peer channel list
# Now the user can proceed to build and start chaincode in one terminal
# And leverage the CLI container to issue install instantiate invoke query commands in another

#we should have bailed if above commands failed.
#we are here, so they worked
sleep 600000
exit 0
