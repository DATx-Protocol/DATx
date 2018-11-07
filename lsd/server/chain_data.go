package server

type ProducerScheduleParams struct {
	Limit int32 `json:"limit"`

	LowerBound int32 `json:"lower_bound"`

	JSON string `json:"json"`
}

type TableParams struct {
	Scope string `json:"scope"`
	Code  string `json:"code"`
	Table string `json:"table"`
	JSON  string `json:"json"`
	Lower int32  `json:"lower"`
	Upper int32  `json:"upper"`
	Limit int32  `json:"limit"`
}

type Producers struct {
	ProducerName    string `json:"producer_name"`
	BlockSigningKey string `json:"block_signing_key"`
}

type ActionsParams struct {
	AccountName string
	Pos         int32
	Offset      int32
}

type ProducerSchedule struct {
	Active struct {
		Version   int         `json:"version"`
		Producers []Producers `json:"producers"`
	} `json:"active"`
	Pending  interface{} `json:"pending"`
	Proposed interface{} `json:"proposed"`
}

type ChainInfo struct {
	ServerVersion            string `json:"server_version"`
	ChainID                  string `json:"chain_id"`
	HeadBlockNum             int    `json:"head_block_num"`
	LastIrreversibleBlockNum int    `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  string `json:"last_irreversible_block_id"`
	HeadBlockID              string `json:"head_block_id"`
	HeadBlockTime            string `json:"head_block_time"`
	HeadBlockProducer        string `json:"head_block_producer"`
	VirtualBlockCPULimit     int    `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit     int    `json:"virtual_block_net_limit"`
	BlockCPULimit            int    `json:"block_cpu_limit"`
	BlockNetLimit            int    `json:"block_net_limit"`
	ServerVersionString      string `json:"server_version_string"`
}

// type Action struct {
// 	Name           string `json:"name"`
// 	TrxID          string `json:"trx_id"`
// 	ActionNum      int    `json:"action_num"`
// 	HandlerAccount string `json:"handler_account"`
// 	Authorization  []struct {
// 		Actor      string `json:"actor"`
// 		Permission string `json:"permission"`
// 	} `json:"authorization"`
// 	Expiration int `json:"expiration"`
// 	Data       struct {
// 		From     string `json:"from"`
// 		To       string `json:"to"`
// 		Quantity string `json:"quantity"`
// 		Memo     string `json:"memo"`
// 	} `json:"data"`
// }

// type AccountActions struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// 	Data    struct {
// 		Page    int      `json:"page"`
// 		PerPage int      `json:"per_page"`
// 		Actions []Action `json:"data"`
// 		Total   int      `json:"total"`
// 		HasNext bool     `json:"has_next"`
// 		HasPrev bool     `json:"has_prev"`
// 	} `json:"data"`
// 	Error interface{} `json:"error"`
// }

type Block struct {
	Timestamp         string        `json:"timestamp"`
	Producer          string        `json:"producer"`
	Confirmed         int           `json:"confirmed"`
	Previous          string        `json:"previous"`
	TransactionMroot  string        `json:"transaction_mroot"`
	ActionMroot       string        `json:"action_mroot"`
	ScheduleVersion   int           `json:"schedule_version"`
	NewProducers      interface{}   `json:"new_producers"`
	HeaderExtensions  []interface{} `json:"header_extensions"`
	ProducerSignature string        `json:"producer_signature"`
	Transactions      []interface{} `json:"transactions"`
	BlockExtensions   []interface{} `json:"block_extensions"`
	ID                string        `json:"id"`
	BlockNum          int           `json:"block_num"`
	RefBlockPrefix    int           `json:"ref_block_prefix"`
}

type ExtractAction struct {
	GlobalActionSeq  int    `json:"global_action_seq"`
	AccountActionSeq int    `json:"account_action_seq"`
	BlockNum         int    `json:"block_num"`
	BlockTime        string `json:"block_time"`
	ActionTrace      struct {
		Receipt struct {
			Receiver       string          `json:"receiver"`
			ActDigest      string          `json:"act_digest"`
			GlobalSequence int             `json:"global_sequence"`
			RecvSequence   int             `json:"recv_sequence"`
			AuthSequence   [][]interface{} `json:"auth_sequence"`
			CodeSequence   int             `json:"code_sequence"`
			AbiSequence    int             `json:"abi_sequence"`
		} `json:"receipt"`
		Act struct {
			Account       string `json:"account"`
			Name          string `json:"name"`
			Authorization []struct {
				Actor      string `json:"actor"`
				Permission string `json:"permission"`
			} `json:"authorization"`
			Data struct {
				From     string `json:"from"`
				To       string `json:"to"`
				Quantity string `json:"quantity"`
				Memo     string `json:"memo"`
			} `json:"data"`
			HexData string `json:"hex_data"`
		} `json:"act"`
		Elapsed       int    `json:"elapsed"`
		CPUUsage      int    `json:"cpu_usage"`
		Console       string `json:"console"`
		TotalCPUUsage int    `json:"total_cpu_usage"`
		TrxID         string `json:"trx_id"`
		InlineTraces  []struct {
			Receipt struct {
				Receiver       string          `json:"receiver"`
				ActDigest      string          `json:"act_digest"`
				GlobalSequence int             `json:"global_sequence"`
				RecvSequence   int             `json:"recv_sequence"`
				AuthSequence   [][]interface{} `json:"auth_sequence"`
				CodeSequence   int             `json:"code_sequence"`
				AbiSequence    int             `json:"abi_sequence"`
			} `json:"receipt"`
			Act struct {
				Account       string `json:"account"`
				Name          string `json:"name"`
				Authorization []struct {
					Actor      string `json:"actor"`
					Permission string `json:"permission"`
				} `json:"authorization"`
				Data struct {
					From     string `json:"from"`
					To       string `json:"to"`
					Quantity string `json:"quantity"`
					Memo     string `json:"memo"`
				} `json:"data"`
				HexData string `json:"hex_data"`
			} `json:"act"`
			Elapsed       int           `json:"elapsed"`
			CPUUsage      int           `json:"cpu_usage"`
			Console       string        `json:"console"`
			TotalCPUUsage int           `json:"total_cpu_usage"`
			TrxID         string        `json:"trx_id"`
			InlineTraces  []interface{} `json:"inline_traces"`
		} `json:"inline_traces"`
	} `json:"action_trace"`
}

//ExtractActions
type ExtractActions struct {
	Actions               []ExtractAction `json:"actions"`
	LastIrreversibleBlock int             `json:"last_irreversible_block"`
}

//Transaction
type Transaction struct {
	ID  string `json:"id"`
	Trx struct {
		Receipt struct {
			Status        string        `json:"status"`
			CPUUsageUs    int           `json:"cpu_usage_us"`
			NetUsageWords int           `json:"net_usage_words"`
			Trx           []interface{} `json:"trx"`
		} `json:"receipt"`
		Trx struct {
			Expiration         string        `json:"expiration"`
			RefBlockNum        int           `json:"ref_block_num"`
			RefBlockPrefix     int           `json:"ref_block_prefix"`
			MaxNetUsageWords   int           `json:"max_net_usage_words"`
			MaxCPUUsageMs      int           `json:"max_cpu_usage_ms"`
			DelaySec           int           `json:"delay_sec"`
			ContextFreeActions []interface{} `json:"context_free_actions"`
			Actions            []struct {
				Account       string `json:"account"`
				Name          string `json:"name"`
				Authorization []struct {
					Actor      string `json:"actor"`
					Permission string `json:"permission"`
				} `json:"authorization"`
				Data struct {
					Hash     string `json:"hash"`
					From     string `json:"from"`
					To       string `json:"to"`
					Quantity string `json:"quantity"`
					Memo     string `json:"memo"`
				} `json:"data"`
				HexData string `json:"hex_data"`
			} `json:"actions"`
			TransactionExtensions []interface{} `json:"transaction_extensions"`
			Signatures            []string      `json:"signatures"`
			ContextFreeData       []interface{} `json:"context_free_data"`
		} `json:"trx"`
	} `json:"trx"`
	BlockTime             string `json:"block_time"`
	BlockNum              int    `json:"block_num"`
	LastIrreversibleBlock int    `json:"last_irreversible_block"`
	Traces                []struct {
		Receipt struct {
			Receiver       string          `json:"receiver"`
			ActDigest      string          `json:"act_digest"`
			GlobalSequence int             `json:"global_sequence"`
			RecvSequence   int             `json:"recv_sequence"`
			AuthSequence   [][]interface{} `json:"auth_sequence"`
			CodeSequence   int             `json:"code_sequence"`
			AbiSequence    int             `json:"abi_sequence"`
		} `json:"receipt"`
		Act struct {
			Account       string `json:"account"`
			Name          string `json:"name"`
			Authorization []struct {
				Actor      string `json:"actor"`
				Permission string `json:"permission"`
			} `json:"authorization"`
			Data struct {
				Hash     string `json:"hash"`
				From     string `json:"from"`
				To       string `json:"to"`
				Quantity string `json:"quantity"`
				Memo     string `json:"memo"`
			} `json:"data"`
			HexData string `json:"hex_data"`
		} `json:"act"`
		Elapsed       int    `json:"elapsed"`
		CPUUsage      int    `json:"cpu_usage"`
		Console       string `json:"console"`
		TotalCPUUsage int    `json:"total_cpu_usage"`
		TrxID         string `json:"trx_id"`
		InlineTraces  []struct {
			Receipt struct {
				Receiver       string          `json:"receiver"`
				ActDigest      string          `json:"act_digest"`
				GlobalSequence int             `json:"global_sequence"`
				RecvSequence   int             `json:"recv_sequence"`
				AuthSequence   [][]interface{} `json:"auth_sequence"`
				CodeSequence   int             `json:"code_sequence"`
				AbiSequence    int             `json:"abi_sequence"`
			} `json:"receipt"`
			Act struct {
				Account       string `json:"account"`
				Name          string `json:"name"`
				Authorization []struct {
					Actor      string `json:"actor"`
					Permission string `json:"permission"`
				} `json:"authorization"`
				Data string `json:"data"`
			} `json:"act"`
			Elapsed       int           `json:"elapsed"`
			CPUUsage      int           `json:"cpu_usage"`
			Console       string        `json:"console"`
			TotalCPUUsage int           `json:"total_cpu_usage"`
			TrxID         string        `json:"trx_id"`
			InlineTraces  []interface{} `json:"inline_traces"`
		} `json:"inline_traces"`
	} `json:"traces"`
}

//SystemProducers ...
type SystemProducers struct {
	Rows []struct {
		Owner         string `json:"owner"`
		TotalVotes    string `json:"total_votes"`
		ProducerKey   string `json:"producer_key"`
		IsActive      int    `json:"is_active"`
		URL           string `json:"url"`
		UnpaidBlocks  int    `json:"unpaid_blocks"`
		LastClaimTime int    `json:"last_claim_time"`
		Location      int    `json:"location"`
	} `json:"rows"`
	TotalProducerVoteWeight string `json:"total_producer_vote_weight"`
	More                    string `json:"more"`
}

//ChargeExpirationTable ...
type ChargeExpirationTable struct {
	Rows []struct {
		ID       int    `json:"id"`
		Trxid    string `json:"trxid"`
		From     string `json:"from"`
		To       string `json:"to"`
		Blocknum int    `json:"blocknum"`
		Quantity string `json:"quantity"`
		Category string `json:"category"`
		Memo     string `json:"memo"`
		Data     string `json:"data"`
		Count    int    `json:"count"`
	} `json:"rows"`
	More bool `json:"more"`
}
