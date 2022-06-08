package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	rdb "github.com/reddio-com/red-adapter/pkg/db"
	"github.com/reddio-com/red-adapter/types"
	"github.com/reddio-com/starkex-contracts-source/source/deposits"
)

// ContractType checks contract type
func GetContractType(db *rdb.DB, contractAddress string) types.TokenType {
	info, _ := db.GetContractInfo(contractAddress)
	return info.Type
}

// IsERC20 check is contract is ERC20
func IsERC20(db *rdb.DB, contractAddress string) bool {
	typ := GetContractType(db, contractAddress)
	return typ == types.ERC20
}

// IsERC721 check is contract is ERC721
func IsERC721(db *rdb.DB, contractAddress string) bool {
	typ := GetContractType(db, contractAddress)
	return typ == types.ERC721
}

type ContractName = string

const (
	Deposit = ContractName("Deposit")
)

const DepositContractAddress = "0x471bda7f420de34282ab8af1f5f3daf2a4c09746"

type ContractMeta struct {
	ContractAddress common.Address
	ABI             *abi.ABI
}

func init() {
	depositABI, _ := deposits.DepositsMetaData.GetAbi()
	depositAddress := common.HexToAddress(DepositContractAddress)
	contracts[Deposit] = ContractMeta{
		ContractAddress: depositAddress,
		ABI:             depositABI,
	}
}

var contracts = make(map[ContractName]ContractMeta)

func GetContractMeta(name ContractName) ContractMeta {
	return contracts[name]
}

func (meta *ContractMeta) ToBoundContract() *bind.BoundContract {
	return bind.NewBoundContract(meta.ContractAddress, *meta.ABI, nil, nil, nil)
}

func (meta *ContractMeta) ToBoundContractWithCaller(client *ethclient.Client) *bind.BoundContract {
	return bind.NewBoundContract(meta.ContractAddress, *meta.ABI, client, client, client)
}
