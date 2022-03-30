package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var _ auth.Account = (*AppAccount)(nil)

// AppAccount is a custom extension for this application. It is an example of
// extending auth.BaseAccount with custom fields. It is compatible with the
// stock auth.AccountKeeper, since auth.AccountKeeper uses the flexible go-amino
// library.
type AppAccount struct {
	auth.BaseAccount
	Name string `json:"name"`

	/* TODO: parse to native macaddr type using net.ParseMAC
	/* macAddress       net.HardwareAddr */
	MacAddress string `json:"macAddress"`
	/* TODO: use float32 to not loose decimal precision during calculation */
	/* TODO: Use struct to represent additional fields (unit ...) */
	ElectricityPrice string          `json:"price"`
	HsmInfo          ccrypto.HsmInfo `json:"hsmInfo"`
}

// nolint
func (acc AppAccount) GetName() string      { return acc.Name }
func (acc *AppAccount) SetName(name string) { acc.Name = name }

func (acc AppAccount) GetMacAddress() string            { return acc.MacAddress }
func (acc *AppAccount) SetMacAddress(macAddress string) { acc.MacAddress = macAddress }

func (acc AppAccount) GetElectricityPrice() string       { return acc.ElectricityPrice }
func (acc *AppAccount) SetElectricityPrice(price string) { acc.ElectricityPrice = price }

func (acc AppAccount) GetHsmInfo() ccrypto.HsmInfo         { return acc.HsmInfo }
func (acc *AppAccount) SetHsmInfo(hsmInfo ccrypto.HsmInfo) { acc.HsmInfo = hsmInfo }

// NewAppAccount returns a reference to a new AppAccount given a name and an
// auth.BaseAccount.
func NewAppAccount(name string, macAddress string, electricityPrice string, hsmInfo ccrypto.HsmInfo, baseAcct auth.BaseAccount) *AppAccount {
	return &AppAccount{BaseAccount: baseAcct, Name: name, MacAddress: macAddress, ElectricityPrice: electricityPrice, HsmInfo: hsmInfo}
}

// GetAccountDecoder returns the AccountDecoder function for the custom
// AppAccount.
func GetAccountDecoder(cdc *codec.Codec) auth.AccountDecoder {
	return func(accBytes []byte) (auth.Account, error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}

		acct := new(AppAccount)
		err := cdc.UnmarshalBinaryBare(accBytes, &acct)
		if err != nil {
			panic(err)
		}

		return acct, err
	}
}

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount reflects a genesis account the application expects in it's
// genesis state.
type GenesisAccount struct {
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
	Price      string `json:"price"`

	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount returns a reference to a new GenesisAccount given an
// AppAccount.
func NewGenesisAccount(aa *AppAccount) *GenesisAccount {
	return &GenesisAccount{
		Name:       aa.Name,
		MacAddress: aa.MacAddress,
		Price:      aa.ElectricityPrice,

		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// ToAppAccount converts a GenesisAccount to an AppAccount.
func (ga *GenesisAccount) ToAppAccount() (acc *AppAccount, err error) {
	return &AppAccount{
		Name:             ga.Name,
		MacAddress:       ga.MacAddress,
		ElectricityPrice: ga.Price,

		BaseAccount: auth.BaseAccount{
			Address: ga.Address,
			Coins:   ga.Coins.Sort(),
		},
	}, nil
}
