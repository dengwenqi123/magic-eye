package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ sdk.Account = (*MgyAccount)(nil)

type MgyAccount struct {
	auth.BaseAccount
	Name string `json:"name"`
}

func (acc MgyAccount) GetName() string      { return acc.Name }
func (acc *MgyAccount) SetName(name string) { acc.Name = name }

//Get the AccountDecoder function for the custom MgyAcount
//获取自定义MgyAccount的AccountDecoder函数

func GetAccountDecoder(cdc *wire.Codec) sdk.AccountDecoder {
	return func(accBytes []byte) (res sdk.Account, err error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}
		acct := new(MgyAccount)
		err = cdc.UnmarshalBinaryBare(accBytes, &acct)
		if err != nil {
			panic(err)
		}
		return acct, err
	}

}

//_________________________________________________

//State to Unmarshal
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

//GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name    string      `json:"name"`
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

func NewGenesisAccount(aa *MgyAccount) *GenesisAccount {
	return &GenesisAccount{
		Name:    aa.Name,
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// convert GenesisAccount to MgyAccount
func (ga *GenesisAccount) ToMGyAccount() (acc *MgyAccount, err error) {
	baseAcc := auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
	return &MgyAccount{
		BaseAccount: baseAcc,
		Name:        ga.Name,
	}, nil
}
