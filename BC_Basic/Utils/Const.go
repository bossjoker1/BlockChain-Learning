package Utils

const (
	DBNAME         = "D:\\go\\src\\BlockChain-Learning\\bolt-db\\bc_%s.db"
	BLOCKTABLENAME = "blocks"
	UTXOTABLENAME  = "utxos"
	LATEST_HASH    = "latest"
	MINEAWARD      = 10 // 挖矿奖励常量
	VERSION        = byte(0x00)
	CHECKSUMLEN    = 4
	WALLETFILEPATH = "D:\\go\\src\\BlockChain-Learning\\bolt-db\\Wallets_%s.dat"
	PROTOCOL       = "tcp"
	VERSION_NUM    = "version"
	NODE_VERSION   = 1
	CMDLENGTH      = 12

	CMD_VERSION   = "version"
	CMD_GETBLOCKS = "getblocks"
	CMD_INV       = "inv"
	CMD_GETDATA   = "getdata"
	CMD_BLOCK     = "block"

	BLOCK_TYPE = "block"
)
