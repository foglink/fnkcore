/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"github.com/foglink/fnkcore/src/common/channelconfig"
	"github.com/foglink/fnkcore/src/common/flogging"
	"github.com/foglink/fnkcore/src/common/genesis"
	"github.com/foglink/fnkcore/src/common/tools/configtxgen/encoder"
	genesisconfig "github.com/foglink/fnkcore/src/common/tools/configtxgen/localconfig"
	cb "github.com/foglink/fnkcore/src/protos/common"
	mspproto "github.com/foglink/fnkcore/src/protos/msp"
	pb "github.com/foglink/fnkcore/src/protos/peer"
	"github.com/foglink/fnkcore/src/protos/utils"
)

var logger = flogging.MustGetLogger("common/configtx/test")

// MakeGenesisBlock creates a genesis block using the test templates for the given chainID
func MakeGenesisBlock(chainID string) (*cb.Block, error) {
	profile := genesisconfig.Load(genesisconfig.SampleDevModeSoloProfile)
	channelGroup, err := encoder.NewChannelGroup(profile)
	if err != nil {
		logger.Panicf("Error creating channel config: %s", err)
	}

	return genesis.NewFactoryImpl(channelGroup).Block(chainID)
}

// MakeGenesisBlockWithMSPs creates a genesis block using the MSPs provided for the given chainID
func MakeGenesisBlockFromMSPs(chainID string, appMSPConf, ordererMSPConf *mspproto.MSPConfig,
	appOrgID, ordererOrgID string) (*cb.Block, error) {
	profile := genesisconfig.Load(genesisconfig.SampleDevModeSoloProfile)
	profile.Orderer.Organizations = nil
	channelGroup, err := encoder.NewChannelGroup(profile)
	if err != nil {
		logger.Panicf("Error creating channel config: %s", err)
	}

	ordererOrg := cb.NewConfigGroup()
	ordererOrg.ModPolicy = channelconfig.AdminsPolicyKey
	ordererOrg.Values[channelconfig.MSPKey] = &cb.ConfigValue{
		Value:     utils.MarshalOrPanic(channelconfig.MSPValue(ordererMSPConf).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}

	applicationOrg := cb.NewConfigGroup()
	applicationOrg.ModPolicy = channelconfig.AdminsPolicyKey
	applicationOrg.Values[channelconfig.MSPKey] = &cb.ConfigValue{
		Value:     utils.MarshalOrPanic(channelconfig.MSPValue(appMSPConf).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}
	applicationOrg.Values[channelconfig.AnchorPeersKey] = &cb.ConfigValue{
		Value:     utils.MarshalOrPanic(channelconfig.AnchorPeersValue([]*pb.AnchorPeer{}).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}

	channelGroup.Groups[channelconfig.OrdererGroupKey].Groups[ordererOrgID] = ordererOrg
	channelGroup.Groups[channelconfig.ApplicationGroupKey].Groups[appOrgID] = applicationOrg

	return genesis.NewFactoryImpl(channelGroup).Block(chainID)
}
