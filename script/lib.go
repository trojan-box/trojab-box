package script

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func getNonce(address string) (string, error) {
	params := url.Values{}
	Url, _ := url.Parse("http://127.0.0.1:5577/api/v1/nonce")
	params.Set("address", address)
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	respResult := vo.Response{}
	err = json.Unmarshal(body, &respResult)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return respResult.Data.(string), nil
}

func playGame(playGameReq vo.PlayGameReq) {
	bytesData, err := json.Marshal(playGameReq)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post("http://127.0.0.1:5577/api/v1/game/play", "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func auth(authReq vo.UserAuthReq) {
	bytesData, err := json.Marshal(authReq)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post("http://127.0.0.1:5577/api/v1/user/auth", "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func generateEthAccount() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Panicf("%v", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Panicf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return fmt.Sprintf("%s,%s", address, privateKeyStr)
}
