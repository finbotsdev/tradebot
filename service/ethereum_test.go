package service

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/test"
	"github.com/stretchr/testify/assert"
)

var ETHEREUM_DIR = "../test/ethereum/blockchain"
var ETHEREUM_KEYSTORE = fmt.Sprintf("%s/keystore", ETHEREUM_DIR)
var ETHEREUM_IPC = fmt.Sprintf("%s/geth.ipc", ETHEREUM_DIR)
var ETHEREUM_PASSPHRASE = "test"
var ETHEREUM_ETHERBASE = "0xe29e121a9da474a103ce39cc83d83a12ddee7073"
var ETHEREUM_BALANCE = "24000000000000000000"

func TestEthereumService_CreateDeleteAccount(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	service, err := NewEthereumService(ctx, ETHEREUM_IPC, ETHEREUM_KEYSTORE, userDAO)
	assert.Nil(t, err)

	err = service.Register("testuser", ETHEREUM_PASSPHRASE)
	assert.Equal(t, nil, err)
	assert.NotNil(t, ctx.GetUser().GetEtherbase())

	err = service.DeleteAccount(ETHEREUM_PASSPHRASE)
	assert.Equal(t, nil, err)

	test.CleanupIntegrationTest()
}

func TestEthereumService_Authenticate(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	service, err := NewEthereumService(ctx, ETHEREUM_IPC, ETHEREUM_KEYSTORE, userDAO)
	assert.Equal(t, nil, err)

	err = service.Authenticate(ctx.GetUser().GetEtherbase(), ETHEREUM_PASSPHRASE)
	assert.Nil(t, err)

	err = service.Authenticate(ctx.GetUser().GetEtherbase(), "nogood")
	assert.NotNil(t, err)
	assert.Equal(t, "could not decrypt key with given passphrase", err.Error())

	err = service.Login(ETHEREUM_PASSPHRASE)
	assert.Nil(t, err)

	err = service.Login("badpass")
	assert.NotNil(t, err)
	assert.Equal(t, "could not decrypt key with given passphrase", err.Error())

	err = service.DeleteAccount(ETHEREUM_PASSPHRASE)
	assert.Equal(t, nil, err)

	test.CleanupIntegrationTest()
}

func TestEthereumService_GetBalance(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	service, err := NewEthereumService(ctx, ETHEREUM_IPC, ETHEREUM_KEYSTORE, userDAO)
	assert.Equal(t, nil, err)

	err = service.Register("testuser", ETHEREUM_PASSPHRASE)
	assert.Equal(t, nil, err)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Mine some coins, press enter when you have a positive balance...")
	_, _ = reader.ReadString('\n')

	balance, err := service.GetBalance(ETHEREUM_ETHERBASE)
	expected := new(big.Int)
	expected.SetString(ETHEREUM_BALANCE, 10)
	assert.Equal(t, expected, balance)

	fmt.Printf("Balance: %d", balance)

	test.CleanupIntegrationTest()
}

func TestEthereumService_Register(t *testing.T) {
	ctx := test.NewIntegrationTestContext()

	userDAO := dao.NewUserDAO(ctx)
	service, err := NewEthereumService(ctx, ETHEREUM_IPC, ETHEREUM_KEYSTORE, userDAO)
	assert.Nil(t, err)

	err = service.Register("ethtest", ETHEREUM_PASSPHRASE)
	assert.Nil(t, err)

	assert.NotNil(t, ctx.GetUser().GetKeystore())
	assert.NotNil(t, ctx.GetUser().GetEtherbase())

	err = service.DeleteAccount(ETHEREUM_PASSPHRASE)
	assert.Nil(t, err)

	newUser, err := userDAO.GetByName("ethtest")
	assert.Nil(t, err)
	assert.Equal(t, true, newUser.GetId() > 0)
	assert.Equal(t, "ethtest", newUser.GetUsername())
	assert.NotNil(t, newUser.GetEtherbase())
	assert.NotNil(t, newUser.GetKeystore())

	test.CleanupIntegrationTest()
}

func getEthereumBalance(t *testing.T) *big.Int {
	bal, err := new(big.Int).SetString(ETHEREUM_BALANCE, 10)
	assert.Nil(t, err)
	return bal
}

func createEthereumSimulatedBackend() *backends.SimulatedBackend {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(150000000)}
	return backends.NewSimulatedBackend(alloc)
}