#!/usr/bin/bash

cryptogen generate --config=./crypto-config.yaml
export FABRIC_CFG_PATH=${PWD}
configtxgen -profile ECOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
export CHANNEL_NAME=elections  && configtxgen -profile VoteTallyChannel -outputCreateChannelTx ./channel-artifacts/elections.tx -channelID $CHANNEL_NAME
configtxgen -profile VoteTallyChannel -outputAnchorPeersUpdate ./channel-artifacts/BallotMachinesMSPanchors.tx -channelID $CHANNEL_NAME -asOrg BallotMachinesMSP
# docker-compose -f docker-compose.yml up -d
