package builder

import (
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/pkg/common"
)

func (o *Output) GetScript() []byte {
	return o.script
}

func (o *Output) HandleAddressInfo(params *chaincfg.Params) error {
	info := common.GetBTCAddressInfo(o.Address)
	if info == nil || info.GetChainConfig().Net != params.Net {
		return errors.New("address not valid")
	}
	o.addressInfo = info
	o.script = info.GetPayToAddrScript()
	return nil
}
