package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
)

func testBlockchain() {
	zilliqa := provider.NewProvider("https://dev-api.zilliqa.com/")

	// These are set by the core protocol, and may vary per-chain.
	// You can manually pack the bytes according to chain id and msg version.
	// For more information: https://apidocs.zilliqa.com/?shell#getnetworkid

	const chainID = 333  // chainId of the developer testnet
	const msgVersion = 1 // current msgVersion
	VERSION := util.Pack(chainID, msgVersion)

	// Populate the wallet with an account
	const privateKey = "3375F915F3F9AE35E6B301B7670F53AD1A5BE15D8221EC7FD5E503F21D3450C8"

	user := account.NewWallet()
	user.AddByPrivateKey(privateKey)
	user.SetDefault("8254b2c9acdf181d5d6796d63320fbb20d4edd12")
	addr := keytools.GetAddressFromPrivateKey(util.DecodeHex(privateKey))
	fmt.Println("My account address is:", user.DefaultAccount.Address)
	fmt.Println("Converting from private key gives:", addr)
	bech, _ := bech32.ToBech32Address(user.DefaultAccount.Address)
	fmt.Println("The bech32 address is:", bech)

	//testing Transaction methods
	bal, _ := zilliqa.GetBalance(user.DefaultAccount.Address).Result.(map[string]interface{})["balance"]
	gas := zilliqa.GetMinimumGasPrice().Result

	fmt.Println("The balance for account ", user.DefaultAccount.Address, " is: ", bal)
	fmt.Println("The blockchain reports minimum gas price: ", gas)

	//Begin importing the contract and init parameters.
	init, err := ioutil.ReadFile("./init.json")
	if err != nil {
		panic(err.Error())
	}
	var initArray []contract.Value
	_ = json.Unmarshal(init, &initArray)
	code, _ := ioutil.ReadFile("./FungibleToken.scilla")

	fmt.Println("Attempting to deploy Fungible Token smart contract...")

	hello := contract.Contract{
		Code:     string(code),
		Init:     initArray,
		Signer:   user,
		Provider: zilliqa,
	}
	nonce, err := zilliqa.GetBalance(string(user.DefaultAccount.Address)).Result.(map[string]interface{})["nonce"].(json.Number).Int64()
	if err != nil {
		fmt.Println("Nonce response error thrown: ", err)
	}
	deployParams := contract.DeployParams{
		Version:      strconv.FormatInt(int64(VERSION), 10),
		Nonce:        strconv.FormatInt(nonce+1, 10),
		GasPrice:     "1000000000",
		GasLimit:     "40000",
		SenderPubKey: string(user.DefaultAccount.PublicKey),
	}
	deployTx, err := DeployWith(&hello, deployParams, "8254B2C9ACDF181D5D6796D63320FBB20D4EDD12")

	if err != nil {
		fmt.Println("Contract deployment failed with error: ", err)
	}

	deployTx.Confirm(deployTx.ID, 1000, 10, zilliqa)

	//verify that the contract is deployed

}

func DeployWith(c *contract.Contract, params contract.DeployParams, pubkey string) (*transaction.Transaction, error) {
	if c.Code == "" || c.Init == nil || len(c.Init) == 0 {
		return nil, errors.New("Cannot deploy without code or initialisation parameters.")
	}

	tx := &transaction.Transaction{
		ID:           params.ID,
		Version:      params.Version,
		Nonce:        params.Nonce,
		Amount:       "0",
		GasPrice:     params.GasPrice,
		GasLimit:     params.GasLimit,
		Signature:    "",
		Receipt:      transaction.TransactionReceipt{},
		SenderPubKey: params.SenderPubKey,
		ToAddr:       "0000000000000000000000000000000000000000",
		Code:         strings.ReplaceAll(c.Code, "/\\", ""),
		Data:         c.Init,
		Status:       0,
	}

	err2 := c.Signer.SignWith(tx, pubkey, *c.Provider)
	if err2 != nil {
		return nil, err2
	}

	rsp := c.Provider.CreateTransaction(tx.ToTransactionPayload())

	if rsp.Error != nil {
		return nil, errors.New(rsp.Error.Message)
	}

	result := rsp.Result.(map[string]interface{})
	hash := result["TranID"].(string)
	contractAddress := result["ContractAddress"].(string)

	tx.ID = hash
	tx.ContractAddress = contractAddress
	return tx, nil

}

func addOperator() {
	host := "https://dev-api.zilliqa.com/"
	privateKey := "3375F915F3F9AE35E6B301B7670F53AD1A5BE15D8221EC7FD5E503F21D3450C8"
	chainID := 333
	msgVersion := 1

	publickKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(publickKey)
	pubkey := util.EncodeHex(publickKey)
	provider := provider.NewProvider(host)

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(privateKey)

	contract := contract.Contract{
		Address:  "zil15e20r8mz6zwqxa7mvg2a72pazvdevcuguafxfp",
		Signer:   wallet,
		Provider: provider,
	}

	//Begin importing the contract and init parameters.
	init, err := ioutil.ReadFile("./operatorAdd.json")
	if err != nil {
		panic(err.Error())
	}
	var args []contract.Value
	_ = json.Unmarshal(init, &args)

	nonce, _ := provider.GetBalance("9bfec715a6bd658fcb62b0f8cc9bfa2ade71434a").Result.(map[string]interface{})["nonce"].(json.Number).Int64()
	n := nonce + 1
	gasPrice := provider.GetMinimumGasPrice().Result.(string)

	params := contract.CallParams{
		Nonce:        strconv.FormatInt(n, 10),
		Version:      strconv.FormatInt(int64(util.Pack(chainID, msgVersion)), 10),
		GasPrice:     gasPrice,
		GasLimit:     "1000",
		SenderPubKey: pubkey,
		Amount:       "0",
	}

	err, tx := contract.Call("Transfer", args, params, true, 1000, 3)
	if err != nil {
		fmt.Printf(err.Error())
	}

	tx.Confirm(tx.ID, 1000, 3, provider)

}

func main() {
	//testBlockchain()
	addOperator()
}
