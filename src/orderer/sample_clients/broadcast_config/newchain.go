/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/foglink/fnkcore/src/common/localmsp"
	"github.com/foglink/fnkcore/src/common/tools/configtxgen/encoder"
	genesisconfig "github.com/foglink/fnkcore/src/common/tools/configtxgen/localconfig"
	cb "github.com/foglink/fnkcore/src/protos/common"
)

func newChainRequest(consensusType, creationPolicy, newChannelId string) *cb.Envelope {
	env, err := encoder.MakeChannelCreationTransaction(newChannelId, localmsp.NewSigner(), nil, genesisconfig.Load(genesisconfig.SampleSingleMSPChannelProfile))
	if err != nil {
		panic(err)
	}
	return env
}
