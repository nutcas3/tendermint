package test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	"testing"
)

// MockBankKeeper is a mock of the bank keeper for testing
type MockBankKeeper struct {
	MintedCoins     sdk.Coins
	SentCoins       sdk.Coins
	SentCoinsToAddr sdk.AccAddress
}

// NewMockBankKeeper returns a new mock bank keeper
func NewMockBankKeeper() *MockBankKeeper {
	return &MockBankKeeper{}
}

// MintCoins implements the BankKeeper interface
func (k *MockBankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	k.MintedCoins = amt
	return nil
}

// SendCoinsFromModuleToAccount implements the BankKeeper interface
func (k *MockBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	k.SentCoins = amt
	k.SentCoinsToAddr = recipientAddr
	return nil
}

// MockStakingKeeper is a mock of the staking keeper for testing
type MockStakingKeeper struct {
	DelegatorStakes map[string]sdk.Int
	TotalBonded     sdk.Int
}

// NewMockStakingKeeper returns a new mock staking keeper
func NewMockStakingKeeper() *MockStakingKeeper {
	return &MockStakingKeeper{
		DelegatorStakes: make(map[string]sdk.Int),
		TotalBonded:     sdk.ZeroInt(),
	}
}

// GetDelegatorStake implements the StakingKeeper interface
func (k *MockStakingKeeper) GetDelegatorStake(ctx sdk.Context, delegator sdk.AccAddress) sdk.Int {
	return k.DelegatorStakes[delegator.String()]
}

// SetDelegatorStake sets the delegator stake for testing
func (k *MockStakingKeeper) SetDelegatorStake(delegator sdk.AccAddress, stake sdk.Int) {
	k.DelegatorStakes[delegator.String()] = stake
}

// GetTotalBondedTokens implements the StakingKeeper interface
func (k *MockStakingKeeper) GetTotalBondedTokens(ctx sdk.Context) sdk.Int {
	return k.TotalBonded
}

// SetTotalBondedTokens sets the total bonded tokens for testing
func (k *MockStakingKeeper) SetTotalBondedTokens(bonded sdk.Int) {
	k.TotalBonded = bonded
}

// MockPosKeeper is a mock of the proof of service keeper for testing
type MockPosKeeper struct {
	ServiceScores map[string]sdk.Int
	TotalScore    sdk.Int
}

// NewMockPosKeeper returns a new mock proof of service keeper
func NewMockPosKeeper() *MockPosKeeper {
	return &MockPosKeeper{
		ServiceScores: make(map[string]sdk.Int),
		TotalScore:    sdk.ZeroInt(),
	}
}

// GetServiceScore implements the ProofOfServiceKeeper interface
func (k *MockPosKeeper) GetServiceScore(ctx sdk.Context, addr string) sdk.Int {
	return k.ServiceScores[addr]
}

// SetServiceScore sets the service score for testing
func (k *MockPosKeeper) SetServiceScore(addr string, score sdk.Int) {
	k.ServiceScores[addr] = score
}

// GetTotalServiceScore implements the ProofOfServiceKeeper interface
func (k *MockPosKeeper) GetTotalServiceScore(ctx sdk.Context) sdk.Int {
	return k.TotalScore
}

// SetTotalServiceScore sets the total service score for testing
func (k *MockPosKeeper) SetTotalServiceScore(score sdk.Int) {
	k.TotalScore = score
}

// MakeTestEncodingConfig creates an EncodingConfig for testing
func MakeTestEncodingConfig() TestEncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := codec.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return TestEncodingConfig{
		Marshaler:         marshaler,
		Amino:             cdc,
		InterfaceRegistry: interfaceRegistry,
	}
}

// TestEncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type TestEncodingConfig struct {
	Marshaler         codec.Codec
	Amino             *codec.LegacyAmino
	InterfaceRegistry codec.InterfaceRegistry
}

func initParamsKeeper(
	cdc codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key storetypes.StoreKey,
	tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(cdc, legacyAmino, key, tkey)

	return paramsKeeper
}

func initKVStore(t *testing.T, storeKey storetypes.StoreKey) storetypes.KVStore {
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	err := stateStore.LoadLatestVersion()
	if err != nil {
		t.Fatal(err)
	}
	return stateStore.GetKVStore(storeKey)
}
