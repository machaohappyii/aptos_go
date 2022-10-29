package aptos

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/tyler-smith/go-bip39"
	"log"
	"strings"
)

// devNet的节点地址
var restUrl = "https://fullnode.devnet.aptoslabs.com"

// mainNet的节点地址
//var restUrl = "https://fullnode.mainnet.aptoslabs.com"

// apt 主币合约地址
//aptAddress := "0x1::aptos_coin::AptosCoin"
// 发布的moon代币合约地址
//tokenAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin"

// 钱包信息
type AccountInfo struct {
	FromAddress string
	PrivateKey  string
	PublicKey   string
}

// 合约币信息
type CoinInfo struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int64  `json:"decimals"`
}

//创建账号 包含生成助记词  私钥  地址为64位
//aptos 钱包使用的是bip39，bip44协议，椭圆曲线使用的是ed25519
func CreateAccount() {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}
	// 生成助记词
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		log.Fatal(err)
	}
	//助记词下面的私钥   44/637  和 以太坊有区别 【eth的path是 m/44'/60'/0'/0/】
	// 修改最后面的 m/44'/637'/0'/0'/1 为第2个私钥
	path := "m/44'/637'/0'/0'/0'"
	key, err := derivation.DeriveForPath(path, seed)
	if err != nil {
		log.Fatal(err)
	}
	account := aptosaccount.NewAccount(key.Key)
	FromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	PrivateKey := hex.EncodeToString(account.PrivateKey[:32])
	PublicKey := hex.EncodeToString(account.PublicKey[:])
	fmt.Println("mnemonic-->", mnemonic)
	fmt.Println("path-->", path)
	fmt.Println("fromAddress-->", FromAddress)
	fmt.Println("PrivateKey-->", PrivateKey)
	fmt.Println("PublicKey-->", PublicKey)
	//return AccountInfo{FromAddress, PrivateKey, PublicKey}
}

//测试签名 交易的时候需要用到
func TestSign() {
	// Import account with private key
	privateKey, err := hex.DecodeString("5372f4c8309ac2875a9c90d4c04a8b7cc668e8795aaafa75e08d119bf82bba73")
	if err != nil {
		fmt.Println(err)
	}
	account := aptosaccount.NewAccount(privateKey)
	fmt.Println("account", account)
	// Get private key, public key, address
	fmt.Printf("privateKey = %x\n", account.PrivateKey[:32])
	fmt.Printf(" publicKey = %x\n", account.PublicKey)
	fmt.Printf("  address = %x\n", account.AuthKey)

	var data = []byte("test")
	// Sign data
	signedData := account.Sign(data, "")
	fmt.Println("signedData", signedData)
}

//查询账户余额
func Balance() {
	client, err := aptosclient.Dial(context.Background(), restUrl)
	if err != nil {
		fmt.Println(err)
	}

	//查询的账户地址
	fromAddress := "0x1bed1a509263cfedc7f9ba5385c9ac2917638e4d415a74e27347a90214df2557"

	//查询本币apt余额 本币地址 0x1::aptos_coin::AptosCoin
	aptAddress := "0x1::aptos_coin::AptosCoin"
	balance, err := client.BalanceOf(fromAddress, aptAddress)
	aptCoinInfo := getCoinInfo(aptAddress)
	decimals := aptCoinInfo.Decimals
	fmt.Println("balance-->", balance)
	fmt.Println("本币apt的balance && decimals", balance, decimals)

	//查询代币moon余额 代币地址 0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin
	tokenAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin"
	tokenBalance, err := client.BalanceOf(fromAddress, tokenAddress)
	tokenCoinInfo := getCoinInfo(aptAddress)
	tokenDecimals := tokenCoinInfo.Decimals
	fmt.Println("代币的balance && tokenDecimals:", tokenBalance, tokenDecimals)

}

// 获取币的 symbol decimals name 信息
func getCoinInfo(address string) CoinInfo {
	addressArr := strings.Split(address, "::")
	user := addressArr[0]
	resourceAddress := "0x1::coin::CoinInfo<" + address + ">"

	client, err := aptosclient.Dial(context.Background(), restUrl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("user:%s\n", user)
	fmt.Printf("user:%s\n", resourceAddress)

	rel, err := client.GetAccountResource(user, resourceAddress, 0)

	aptTokenCoin := CoinInfo{Name: rel.Data["name"].(string), Symbol: rel.Data["symbol"].(string), Decimals: int64(rel.Data["decimals"].(float64))}
	fmt.Printf("aptTokenCoin=%+v\n\r", aptTokenCoin)
	return aptTokenCoin
}

// apt 代币账户交易 转账
func TransferApt() {
	// true 为转本币  false 为其他币
	coinAddress := "0x1::aptos_coin::AptosCoin"
	privateKey, err := hex.DecodeString("2e939a3cb7f1697e8c9c4d1eeb1d04e6818d34bf50fb0826dcbe500a3976a19e")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("privateKey", privateKey)
	account := aptosaccount.NewAccount(privateKey)
	fmt.Println("account", account)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	fmt.Println("fromAddress", fromAddress)
	//0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f
	toAddress := "0x8bdc529d93c6ec5a48070af1c7a703d6a20d2b9f457f565b3fa4b7c087e06ca6"
	amount := "10002000"
	client, err := aptosclient.Dial(context.Background(), restUrl)

	// Get Sender's account data and ledger info
	accountData, err := client.GetAccount(fromAddress)
	ledgerInfo, err := client.LedgerInfo()
	// Get gas price
	gasPrice, err := client.EstimateGasPrice()

	// 该方法 支持 apt 和 代币转账 唯一的缺点就是 对方没有registerCoin 无法转账
	// Build paylod
	//payload := &aptostypes.Payload{
	//	Type:          "entry_function_payload",
	//	Function:      "0x1::coin::transfer",
	//	TypeArguments: []string{coinAddress},
	//	Arguments: []interface{}{
	//		toAddress, amount,
	//	},
	//}
	fmt.Println(coinAddress)

	// 该方法 目前支持 apt转账 优点是对方没有注册是  会自动注册 apt交易优先使用这个方法
	payload := &aptostypes.Payload{
		Type:          "entry_function_payload",
		Function:      "0x1::aptos_account::transfer",
		TypeArguments: []string{},
		Arguments: []interface{}{
			toAddress, amount,
		},
	}

	// Build transaction
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	// Get signing message from remote server
	// Note: Later we will implement the local use of BCS encoding to create signing messages
	signingMessage, err := client.CreateTransactionSigningMessage(transaction)

	// Sign message and complete transaction information
	signatureData := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}
	// Submit transaction
	newTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		fmt.Println("交易失败", err)
	}

	fmt.Printf("TransferApt hash = %v\n", newTx.Hash)
}

//	其他代币 账户交易 转账  主要以 moon币为例
//  moon代币的发布用户地址为: 0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f
//	这个有点尴尬 因为官方介绍 用户的 source资源是第一公民 不允许被别人私自空投  所以想转账的话  转账的地址必须先注册 registerCoin
//	后期想想怎么绕过这个步骤
func TransferToken() {
	// true 为转本币  false 为其他币
	coinAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin"
	privateKey, err := hex.DecodeString("2e939a3cb7f1697e8c9c4d1eeb1d04e6818d34bf50fb0826dcbe500a3976a19e")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("privateKey", privateKey)

	account := aptosaccount.NewAccount(privateKey)
	fmt.Println("account", account)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	fmt.Println("fromAddress", fromAddress)
	//
	//0x8bdc529d93c6ec5a48070af1c7a703d6a20d2b9f457f565b3fa4b7c087e06ca6
	toAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f"
	amount := "200000"

	client, err := aptosclient.Dial(context.Background(), restUrl)

	// Get Sender's account data and ledger info
	accountData, err := client.GetAccount(fromAddress)
	ledgerInfo, err := client.LedgerInfo()

	// Get gas price
	gasPrice, err := client.EstimateGasPrice()

	//该方法 支持 apt 和 代币转账 唯一的缺点就是 对方没有registerCoin 无法转账
	//Build paylod
	payload := &aptostypes.Payload{
		Type:          "entry_function_payload",
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{coinAddress},
		Arguments: []interface{}{
			toAddress, amount,
		},
	}
	fmt.Println(coinAddress)

	// Build transaction
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	// Get signing message from remote server
	// Note: Later we will implement the local use of BCS encoding to create signing messages
	signingMessage, err := client.CreateTransactionSigningMessage(transaction)

	// Sign message and complete transaction information
	signatureData := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}
	// Submit transaction
	newTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		fmt.Println("交易失败", err)
	}
	fmt.Printf("TransferToken hash = %v\n", newTx.Hash)
}

//	注册代币  主要以 moon币为例  记得注册前先领取下水龙头 当做gas
//  moon代币的发布用户地址为: 0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f
//	这个有点尴尬 因为官方介绍 用户的 source资源是第一公民 不允许被别人私自空投  所以想转账的话  转账的地址必须先注册 registerCoin
//	后期想想怎么绕过这个步骤  目前就是注册这一步
func RegisterCoin() {
	privateKey, err := hex.DecodeString("efa6ff46edca862f3e712f6e2645586d5c28c72b93355fb34b933fb7e0dd1454")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("privateKey", privateKey)

	account := aptosaccount.NewAccount(privateKey)
	fmt.Println("account", account)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	fmt.Println("fromAddress", fromAddress)

	coinAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin"
	isRegister := checkRegister(fromAddress, coinAddress)
	if isRegister {
		fmt.Println("用户已注册过代币，无需再注册")
		return
	}

	client, err := aptosclient.Dial(context.Background(), restUrl)

	// Get Sender's account data and ledger info
	accountData, err := client.GetAccount(fromAddress)
	ledgerInfo, err := client.LedgerInfo()

	// Get gas price
	gasPrice, err := client.EstimateGasPrice()

	//Build paylod
	payload := &aptostypes.Payload{
		Type:          "entry_function_payload",
		Function:      "0x1::managed_coin::register",
		TypeArguments: []string{coinAddress},
		Arguments:     []interface{}{},
	}
	fmt.Println(coinAddress)

	// Build transaction
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	// Get signing message from remote server
	// Note: Later we will implement the local use of BCS encoding to create signing messages
	signingMessage, err := client.CreateTransactionSigningMessage(transaction)

	// Sign message and complete transaction information
	signatureData := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}
	// Submit transaction
	newTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		fmt.Println("regisertCoin交易失败", err)
	}
	fmt.Printf("regisertCoin tx hash = %v\n", newTx.Hash)
}

// 检查是否注册代币
func checkRegister(user string, coinAddress string) bool {
	resourceAddress := "0x1::coin::CoinStore<" + coinAddress + ">"
	client, err := aptosclient.Dial(context.Background(), restUrl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("user:%s\n", user)
	fmt.Printf("user:%s\n", resourceAddress)

	rel, err := client.GetAccountResource(user, resourceAddress, 0)
	if err != nil {
		fmt.Println(err)
	}
	value, ok := rel.Data["coin"]
	if ok {
		fmt.Printf("用户已注册= %+v", value)
		return true
	} else {
		fmt.Printf("用户未注注册= %+v", value)
		return false
	}
}

//	mint测试币 主要以 moon币为例  必须是发布者才能调用mintCoin  记得注册前先领取下水龙头 当做gas
//  moon代币的发布用户地址为: 0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f
//  领取完就可以测试下 代币的转账了
func MintCoin() {
	privateKey, err := hex.DecodeString("56cce9b64e140c55d065d1150452129aa128bbbf445190552df8a42011beb8b6")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("privateKey", privateKey)

	account := aptosaccount.NewAccount(privateKey)
	fmt.Println("account", account)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	fmt.Println("fromAddress", fromAddress)

	coinAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin"
	isRegister := checkRegister(fromAddress, coinAddress)
	if !isRegister {
		fmt.Println("请先注册在领取")
		return
	}
	client, err := aptosclient.Dial(context.Background(), restUrl)

	// Get Sender's account data and ledger info
	accountData, err := client.GetAccount(fromAddress)
	ledgerInfo, err := client.LedgerInfo()

	// Get gas price
	gasPrice, err := client.EstimateGasPrice()

	toAddress := "0x8bdc529d93c6ec5a48070af1c7a703d6a20d2b9f457f565b3fa4b7c087e06ca6"
	//Build paylod
	payload := &aptostypes.Payload{
		Type:          "entry_function_payload",
		Function:      "0x1::managed_coin::mint",
		TypeArguments: []string{coinAddress},
		Arguments: []interface{}{
			toAddress, "100000",
		},
	}
	fmt.Println(coinAddress)

	// Build transaction
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	// Get signing message from remote server
	// Note: Later we will implement the local use of BCS encoding to create signing messages
	signingMessage, err := client.CreateTransactionSigningMessage(transaction)

	// Sign message and complete transaction information
	signatureData := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}
	// Submit transaction
	newTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		fmt.Println("MintCoin交易失败", err)
	}
	fmt.Printf("MintCoin tx hash = %v\n", newTx.Hash)
}

//	mint测试币 主要以 moon币为例  必须是发布者才能调用mintCoin  记得注册前先领取下水龙头 当做gas
//  销毁代币
func BurnCoin() {
	privateKey, err := hex.DecodeString("56cce9b64e140c55d065d1150452129aa128bbbf445190552df8a42011beb8b6")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("privateKey", privateKey)

	account := aptosaccount.NewAccount(privateKey)
	fmt.Println("account", account)

	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	fmt.Println("fromAddress", fromAddress)

	coinAddress := "0x39edef225b4d840416209012f7553d216f9ad62eed04f428059a1e1215df4d2f::moon_coin::MoonCoin"
	isRegister := checkRegister(fromAddress, coinAddress)
	if !isRegister {
		fmt.Println("请先注册在领取")
		return
	}
	client, err := aptosclient.Dial(context.Background(), restUrl)

	// Get Sender's account data and ledger info
	accountData, err := client.GetAccount(fromAddress)
	ledgerInfo, err := client.LedgerInfo()

	// Get gas price
	gasPrice, err := client.EstimateGasPrice()

	//Build paylod
	payload := &aptostypes.Payload{
		Type:          "entry_function_payload",
		Function:      "0x1::managed_coin::burn",
		TypeArguments: []string{coinAddress},
		Arguments: []interface{}{
			"100000",
		},
	}
	fmt.Println(coinAddress)

	// Build transaction
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	// Get signing message from remote server
	// Note: Later we will implement the local use of BCS encoding to create signing messages
	signingMessage, err := client.CreateTransactionSigningMessage(transaction)

	// Sign message and complete transaction information
	signatureData := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}
	// Submit transaction
	newTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		fmt.Println("MintCoin交易失败", err)
	}
	fmt.Printf("MintCoin tx hash = %v\n", newTx.Hash)
}
