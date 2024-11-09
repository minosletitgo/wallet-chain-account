package cosmos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/protobuf/ptypes"
)

const (
	NetWork   = "mainnet"
	ChainName = "Cosmos"
)

type ChainAdaptor struct {
	client CosmosClient
	conf   *config.Config
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	cosmosClient, err := DialCosmosClient(context.Background(), conf)
	if err != nil {
		log.Error("new chain adaptor error (%w)", err)
		return nil, err
	}
	return &ChainAdaptor{
		client: *cosmosClient,
		conf:   conf,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	supportList := []string{"stake", "cosmos", "atom"}

	checkIf := func(s string) bool {
		for _, v := range supportList {
			if strings.EqualFold(v, s) {
				return true
			}
		}
		return false
	}

	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: checkIf(req.Chain),
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	addr := c.client.GetAddressFromPubKey([]byte(req.PublicKey))

	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: addr,
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	_, err := types2.AccAddressFromBech32(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   err.Error(),
			Valid: false,
		}, err
	}
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "valid address success",
		Valid: true,
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.client.GetBlock(req.Height)
	if err != nil {
		log.Error("get block by number error (%w)", err)
		return nil, err
	}

	for _, txData := range block.Block.Txs {
		// hash ok
		txHash := fmt.Sprintf("%x", string(txData.Hash()))
		fmt.Printf("tx0=%s ,tx hash:%s \n", txData.String(), fmt.Sprintf("%x", string(txData.Hash())))
		result, _ := c.client.Tx(txHash, true)
		for _, event := range result.TxResult.Events {
			if event.Type == "transfer" {
				fmt.Printf("event type: %s, \n", event.Type)
				for _, attr := range event.GetAttributes() {
					fmt.Printf("attribute key: %s, value: %s , index: %v \n", attr.GetKey(), attr.GetValue(), attr.GetIndex())
				}
			}
		}

		c.client.ParseTx(txData)
		//fmt.Printf("tx1=%s \n", []byte(txData))
		//fmt.Printf("tx2=%x \n", txData)
		//
		//decodeBytes, _ := base64.StdEncoding.DecodeString(string(txData))
		//fmt.Printf("decodeBytes=%s \n", string(decodeBytes))
		//encodeStr := base64.StdEncoding.EncodeToString(txData)
		//fmt.Printf("encodeStr=%s \n", encodeStr)
		//
		//if err != nil {
		//	log.Error("failed to decode base64 transaction: %v", err)
		//}
		//response, err := c.client.TxDecode(txData)
		//if err != nil {
		//	log.Error("get block by number error (%w)", err)
		//	continue
		//}
		//
		//log.Info("response=%s", response.GetTx().String())
		//decodeBytes1, _ := base64.StdEncoding.DecodeString(response.GetTx().String())
		//fmt.Printf("decodeBytes1=%s \n", string(decodeBytes1))
		//encodeStr1 := base64.StdEncoding.EncodeToString([]byte(response.GetTx().String()))
		//fmt.Printf("encodeStr1=%s \n", encodeStr1)
	}

	return nil, nil
	//
	//transactions, err := c.client.DecodeBlockTx(c.conf.WalletNode.Cosmos.DataApiUrl, block)
	//
	//if err != nil {
	//	log.Error("decode block tx error (%w)", err)
	//	return nil, err
	//}
	//blockHeight, err := strconv.ParseInt(block.Block.Header.Height, 10, 64)
	//if err != nil {
	//	log.Error("parse block height error (%w)", err)
	//	return nil, err
	//}
	//// BaseFee
	//return &account.BlockResponse{
	//	Code:         common2.ReturnCode_SUCCESS,
	//	Msg:          "get block by number success",
	//	Height:       blockHeight,
	//	Hash:         block.BlockId.Hash,
	//	Transactions: transactions,
	//}, nil
	return nil, nil
}

// error
func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//	hexBytes, err := hex.DecodeString(req.GetHash())
	//	if err != nil {
	//		log.Error("get block header by hash decode hash error (%w)", err)
	//		return nil, err
	//	}
	//	hashBlock, err := c.client.GetBlockByHash(hexBytes)
	//=======
	//func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//	block, err := c.client.GetBlockByHash([]byte(req.Hash))
	//>>>>>>> main
	//	if err != nil {
	//		log.Error("get block by hash error (%w)", err)
	//		return nil, err
	//	}
	//<<<<<<< HEAD
	//	block, err := c.client.GetBlock(c.conf.WalletNode.Cosmos.RestUrl, hashBlock.Block.Height)
	//	if err != nil {
	//		log.Error("get block by number error (%w)", err)
	//		return nil, err
	//	}
	//
	//	transactions, err := c.client.DecodeBlockTx(c.conf.WalletNode.Cosmos.RestUrl, block)
	//	if err != nil {
	//		log.Error("decode block tx error (%w)", err)
	//		return nil, err
	//	}
	//
	//	height, err := strconv.ParseInt(block.Block.Header.Height, 10, 64)
	//	if err != nil {
	//		log.Error("decode block tx parse height  error (%w)", err)
	//		return nil, err
	//	}
	//
	//=======
	//
	//	blockResponse, err := c.client.GetBlock(c.conf.WalletNode.Cosmos.DataApiUrl, block.Block.Header.Height)
	//	log.Info("block tx : %s", blockResponse.Block.Data.Txs[0])
	//>>>>>>> main
	//	// BaseFee
	//	return &account.BlockResponse{
	//		Code:         common2.ReturnCode_SUCCESS,
	//		Msg:          "get block by hash success",
	//<<<<<<< HEAD
	//		Height:       height,
	//		Hash:         block.BlockId.Hash,
	//		Transactions: transactions,
	//=======
	//		Height:       block.Block.Height,
	//		Hash:         string(block.Block.Hash()),
	//		Transactions: nil,
	//>>>>>>> main
	//	}, nil
	return nil, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	header, err := c.client.GetHeaderByHash(req.GetHash())
	if err != nil {
		log.Error("get block header by hash error (%w)", err)
		return nil, err
	}
	// todo gas field
	blockHeader := &account.BlockHeader{
		Hash:       header.Header.Hash().String(),
		TxHash:     header.Header.DataHash.String(),
		ParentHash: header.Header.AppHash.String(),
		Number:     strconv.FormatInt(header.Header.Height, 10),
		Time:       uint64(header.Header.Time.Unix()),
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by hash success",
		BlockHeader: blockHeader,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	header, err := c.client.GetHeaderByHeight(req.Height)
	if err != nil {
		log.Error("get block header by number error (%w)", err)
		return nil, err
	}
	// todo gas field
	blockHeader := &account.BlockHeader{
		Hash:       header.Header.Hash().String(),
		TxHash:     header.Header.DataHash.String(),
		ParentHash: header.Header.AppHash.String(),
		Number:     strconv.FormatInt(header.Header.Height, 10),
		Time:       uint64(header.Header.Time.Unix()),
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by number success",
		BlockHeader: blockHeader,
	}, nil
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	response, err := c.client.GetAccount(req.Address)
	if err != nil {
		log.Error("get account error (%w)", err)
		return nil, err
	}

	authAccount := new(authv1beta1.BaseAccount)
	if err := ptypes.UnmarshalAny(response.Account, authAccount); err != nil {
		log.Error("get account error (%w)", err)
		return nil, err
	}
	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account success",
		Network:       NetWork,
		AccountNumber: strconv.FormatUint(authAccount.AccountNumber, 10),
		Sequence:      strconv.FormatUint(authAccount.Sequence, 10),
		Balance:       strconv.FormatUint(authAccount.AccountNumber, 10),
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	ret, err := c.client.BroadcastTx([]byte(req.RawTx))
	if err != nil {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "BroadcastTx fail",
		}, err
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: ret.TxResponse.TxHash,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO 需接第三方
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	//<<<<<<< HEAD
	//	txResponse, err := c.client.GetTxByHash(c.conf.WalletNode.Cosmos.RestUrl, req.Hash)
	//=======
	//	txResponse, err := c.client.GetTxByHash(c.conf.WalletNode.Cosmos.DataApiUrl, req.Hash)
	//>>>>>>> main
	//	if err != nil {
	//		return &account.TxHashResponse{
	//			Code: common2.ReturnCode_ERROR,
	//			Msg:  "get tx by hash fail",
	//		}, err
	//	}
	//
	//	index := int64(0)
	//	fromAddr, toAddr, amount := "", "", ""
	//	for _, v := range txResponse.Response.Events {
	//		if v.Type == "transfer" {
	//			for _, attr := range v.Attributes {
	//				if attr.Key == "recipient" {
	//					toAddr = attr.Value
	//				}
	//				if attr.Key == "sender" {
	//					fromAddr = attr.Value
	//				}
	//				if attr.Key == "amount" {
	//					amount = attr.Value
	//				}
	//				if attr.Key == "msg_index" {
	//					index, _ = strconv.ParseInt(attr.Value, 10, 32)
	//				}
	//			}
	//		}
	//	}
	//	log.Info("tx hash: %s, amount: %s", req.GetHash(), amount)
	//
	//	return &account.TxHashResponse{
	//<<<<<<< HEAD
	//		Code: common2.ReturnCode_SUCCESS,
	//		Msg:  "get tx by hash success",
	//=======
	//>>>>>>> main
	//		Tx: &account.TxMessage{
	//			Hash:            req.Hash,
	//			Index:           uint32(index),
	//			Froms:           []*account.Address{{Address: fromAddr}},
	//			Tos:             []*account.Address{{Address: toAddr}},
	//			Values:          nil,
	//			Fee:             txResponse.Response.GasUsed,
	//			Status:          account.TxStatus_Success,
	//			Type:            0,
	//			Height:          txResponse.Response.Height,
	//			ContractAddress: "",
	//			Datetime:        txResponse.Response.Timestamp,
	//		},
	//	}, nil
	return nil, nil
}

// fail
func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	minHeight, err := strconv.ParseInt(req.GetStart(), 10, 64)
	if err != nil {
		log.Error("min height invalid", err)
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block range fail ! start height err",
		}, err
	}
	maxHeight, err := strconv.ParseInt(req.GetEnd(), 10, 64)
	if err != nil {
		log.Error("max height invalid", err)
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block range fail ! start height err",
		}, err
	}
	blockInfo, err := c.client.BlockchainInfo(minHeight, maxHeight)
	log.Info("block metas len: %d", len(blockInfo.BlockMetas))
	return &account.BlockByRangeResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get block by range success",
	}, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	return &account.UnSignTransactionResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "create unsigned transaction success",
	}, nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	return &account.SignedTransactionResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "build signed transaction success",
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "decode transaction success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify signed transaction success",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}
