#!/usr/bin/bash

cryptogen generate --config=./crypto-config.yaml
export FABRIC_CFG_PATH=${PWD}
configtxgen -profile ECOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
export CHANNEL_NAME=elections  && configtxgen -profile VoteTallyChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
docker-compose -f docker-compose.yml up -d
