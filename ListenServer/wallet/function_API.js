//钱包信息，在多个页面可以共用
//请求样例
//curl http://127.0.0.1:8880/wallet_info -d '{"category":"DATX","address":"datx24dotcavazu"}'
//响应结构体
type WalletInfo struct {
	Category  string `json:"category"` // DATX或BTC或ETH，三种钱包
	Address   string `json:"address"`  // DATX填账号
	TokenInfo struct {
		Symbol  string  `json:"symbol"`  // 代币，DATX，DBTC，DETH，DEOS，BTC，ETH
		Balance float64 `json:"balance"` // 数量
		Price   float64 `json:"price"`   // 单价（USD）
	} `json:"tokeninfo"`
}

//代币信息，在多个页面可以共用
//请求样例
//curl http://127.0.0.1:8880/token_info -d '{"symbol":"DATX"}'
//响应结构体
type TokenInfo struct {
	Symbol  string  `json:"symbol"`  // 代币，DATX，DBTC，DETH，DEOS，BTC，ETH
	Balance float64 `json:"balance"` // 数量
	Price   float64 `json:"price"`   // 单价（USD）
}

//矿工费
//请求样例
//curl http://127.0.0.1:8880/token_fee -d '{"symbol":"BTC"}'
//响应结构体
type TokenFee struct {
	Symbol   string  `json:"symbol"`   // 代币，BTC，ETH
	Regular  float64 `json:"regular"`  // 六个块内确认
	Priority float64 `json:"priority"` // 最快确认
}

//转出
//请求样例
//curl http://127.0.0.1:8880/transfer_out -d '{"category":"DATX","address":"datx24dotcavazu","to":"0x5f1g2a144654asd54f34gh3gdsqfad13543514s3d4f1343g4","quantity":0.1765,"price":8,,"limit":10,"memo":"test"}'
//请求结构体
type TransferOutInfo struct {
	Category string  `json:"category"` // DATX或BTC或ETH
	Address  string  `json:"address"`  // 相当于From
	To       string  `json:"to"`       // 收款地址
	Quantity float64 `json:"quantity"` // 转账金额
	Price    float64 `json:"price"`    // 矿工费，单价
	Limit    float64 `json:"limit"`    // ETH专用
	Memo     string  `json:"memo"`     // DATX专用
}

//转入
//页面生成地址二维码

//提币
//请求样例
//curl http://127.0.0.1:8880/draw -d '{"category":"DATX","address":"datx24dotcavazu","quantity":0.1765,"password":"test"}'
//请求结构体
type DrawInfo struct {
	Category string  `json:"category"`
	Address  string  `json:"address"` // 相当于To
	Quantity float64 `json:"quantity"`
	Password string  `json:"password"`
}

//充值
//请求样例
//curl http://127.0.0.1:8880/charge -d '{"category":"DATX","address":"datx24dotcavazu","quantity":0.1765,"password":"test","memo":"datx24dotcavazu"}'
//请求结构体
type ChargeInfo struct {
	Category string  `json:"category"`
	Address  string  `json:"address"` // 相当于From
	Quantity float64 `json:"quantity"`
	Memo     string  `json:"memo"` // 相当于From
	Password string  `json:"password"`
}

//导入钱包，这里提供了接口，但是我感觉不需要访问服务器，只要本地记录私钥密码等就可以了
//请求样例
//curl http://127.0.0.1:8880/import_wallet -d '{"category":"DATX","privatekey":"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3","memo":"wink job cat scheme tuna happy hawk clump notable small spider pupil","keystore":"太长省略","password":"test"}'
//请求结构体
type ImportWalletInfo struct {
	Category   string `json:"category"`   // BTC或ETH
	PrivateKey string `json:"privatekey"` // 导入私钥
	Memo       string `json:"memo"`       // 导入助记词
	Keystore   string `json:"keystore"`   // 导入Keystore
	Password   string `json:"password"`   // 设置密码
}
