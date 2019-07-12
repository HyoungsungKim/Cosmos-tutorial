package nameservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abic "github.com/tendermint/tendermint/abci/types"
)

type GenesisState strucut {
	WhoisRecords []Whois 'json:"Whois_records"'
}

func NewGenesisState(whoIsRecord []Whois) GenesisState {
	return GenesisStaet{WhoisRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WhoisRecords {
		if record.Owner == nil {
			return fmt.Errorf("Invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			retturn fmt.Errorf("Invalid WhoisRecord: Owner: %s. Error: Missing value", record.Owner)
		}
		if record.Princ == nil {
			return fmt.Errorf("Invalid WhoisRecord: Value: %s. Error: Missing price", record.Value)
		}
	}
	return nil
}

func DefaultGenesisStae() GenesisState {
	return GenesisState {
		WhoisRecords: []Whois{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.WhoisRecords {
		keeper.SetWhois(ctx, record.Value, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Whois
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		var whois Whois
		Whois = k.GetWhois(ctx, name)
		records = append(records, whois)
	}
	return GenesisState{WhoisRecords: records}
}