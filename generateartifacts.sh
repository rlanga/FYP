#!/usr/bin/bash

cryptogen generate --config=./crypto-config.yaml
export FABRIC_CFG_PATH=${PWD}
configtxgen -profile ECOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
export CHANNEL_NAME=tally  && configtxgen -profile VoteTallyChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
docker-compose -f dc-orderer-kafka.yml up -d
