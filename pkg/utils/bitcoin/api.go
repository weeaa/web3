package bitcoin

import (
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"io"
	"net/url"
)

type MempoolClient struct {
	baseURL string
}

func NewClient(netParams *chaincfg.Params) *MempoolClient {
	var baseURL string

	if netParams.Net == wire.MainNet {
		baseURL = "https://mempool.space/api"
	} else if netParams.Net == wire.TestNet3 {
		baseURL = "https://mempool.space/testnet/api"
	} else if netParams.Net == chaincfg.SigNetParams.Net {
		baseURL = "https://mempool.space/signet/api"
	} else {
		//log.Fatal("ERROR MemPool !=s netParams")
	}

	return &MempoolClient{
		baseURL: baseURL,
	}
}

func (c *MempoolClient) request(method, subPath string, requestBody io.Reader) ([]byte, error) {
	return Request(method, c.baseURL, subPath, requestBody)
}

func Request(method, baseURL, subPath string, requestBody io.Reader) ([]byte, error) {
	URL, err := url.Parse(fmt.Sprintf("%s%s", baseURL, subPath))
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: method,
		URL:    URL,
		Body:   io.NopCloser(requestBody),
		Header: http.Header{
			"content-type": {"application/json"},
			"accept":       {"application/json"},
		},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type UTXO struct {
	TxID   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int64  `json:"block_time"`
	} `json:"status"`
	Value int64 `json:"value"`
}

type UnspentOutput struct {
	Outpoint *wire.OutPoint
	Output   *wire.TxOut
}

// UTXOs is a slice of UTXO
type UTXOs []UTXO

func (c *MempoolClient) ListUnspent(address btcutil.Address) ([]*UnspentOutput, error) {
	res, err := c.request(http.MethodGet, fmt.Sprintf("/address/%s/utxo", address.EncodeAddress()), nil)
	if err != nil {
		return nil, err
	}

	var UTXOS UTXOs
	if err = json.Unmarshal(res, &UTXOS); err != nil {
		return nil, err
	}

	unspentOutputs := make([]*UnspentOutput, 0)
	for _, utxo := range UTXOS {
		var txHash *chainhash.Hash

		txHash, err = chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}
		unspentOutputs = append(unspentOutputs, &UnspentOutput{
			Outpoint: wire.NewOutPoint(txHash, uint32(utxo.Vout)),
			Output:   wire.NewTxOut(utxo.Value, address.ScriptAddress()),
		})
	}

	return unspentOutputs, nil
}
