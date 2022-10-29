package main

import "go_aptos/aptos"

func main() {
	//创建账号
	aptos.CreateAccount()

	//测试签名
	aptos.TestSign()

	//查询余额
	aptos.Balance()

	// apt本币转账
	aptos.TransferApt()

	// apt链的token币转账 发布的新代币moon
	aptos.TransferToken()

	// 注册代币
	aptos.RegisterCoin()

	// 给别人mint代币
	aptos.MintCoin()

	// burn代币
	aptos.BurnCoin()
}
