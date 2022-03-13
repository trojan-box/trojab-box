package script

import (
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBatchPlayGame(t *testing.T) {

	keys := []string{
		"0x9866c5dB69592ade1F1179D5110fDc7Be46F6f25,ae62e62adb70f4684bf9cc463eb4b4c5dab932da959e71f2a0c09c3ac5488c8b",
		//"0x92561002B6BB157B7839E05953094cF47d8c7B61,954b3d334bcb0447018c33f00665a4f0ca7241b6e991e8107e4ef32f6f13a780",
		//"0x7c7A5e0A0A691D8913A24EDC6A367D24D93831dD,d425ba17713cc45e13d201f0fbc4ec0adb3fcfcb4188be044e285eb6281a8cd3",
		//"0x1170e45A16fdf7410beb68C913DF94c647366e09,4a18fbefa70befdb5b6aee6a2861fd438d77558c5e753293f5e6e4320024fc94",
		//"0xF6275476ce251366C2ad06d6137Cae6b0ad5C4EB,4aa778bc83ba2fb72bb15bd02dd29471fdbf1ed695f989abf592c8d3bf136b35",
		//"0x7A536907262802112D8937410EDf397EA8d06d49,6b25f952eb64299696602a62caab193150a5e058a32a2c769ad095118d8ad26c",
		//"0xd6069AeC76C354465A7F1C2F9c815A5BCC80A75b,823440e4320f870be13e4b174625fa5e1fee7310fbb60994d37eaea67cffdc6a",
		//"0x2aa55220d5494BB890e53C20e5b18472aAe57878,6591c0c35711bdf6b389b0cad24b30e6a2ded9342e823c1040eca351672f36b6",
		//"0xd9389a489CF197F60aA1eDd24047963b80aA8879,ebc69b92582ba03c920ea02dbd66ee6d2b551629df1736b40a183e3daa313ed8",
		//"0x868F3398050F394196b4eae96cB436E6d5Aab9d4,0f88333ea82217284291da6d05dc2abc8741791db184a9ec21915e7a86af270b",
		//"0xCA156b8163121c772e1c84F2CE34f32680A898e8,7835e4affeab5ef9734e22687d3631240faf36b82d22b816e4e4136b02a7b942",
	}
	for _, keyTemp := range keys {
		keysArr := strings.Split(keyTemp, ",")
		address := keysArr[0]
		privateKeyStr := keysArr[1]

		nonce, _ := getNonce(address)
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		cards := "[1 2 3 4 5 6 7 8 9]"
		chosen := "4"

		msg := strings.ReplaceAll(constant.PlayGameSignMsg, "${address}", address)
		msg = strings.ReplaceAll(msg, "${nonce}", nonce)
		msg = strings.ReplaceAll(msg, "${timestamp}", timestamp)
		msg = strings.ReplaceAll(msg, "${cards}", cards)
		msg = strings.ReplaceAll(msg, "${chosen}", chosen)

		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("sign msg :", msg)
		data := []byte(msg)
		msgHash := crypto.Keccak256Hash(data)
		fmt.Println("msgHash", msgHash.Hex())
		signature, err := crypto.Sign(msgHash.Bytes(), privateKey)
		if err != nil {
			log.Fatal(err)
		}
		signatureMsg := hexutil.Encode(signature)
		fmt.Println("signature msg", signatureMsg)

		playGameReq := vo.PlayGameReq{
			Address:     address,
			Nonce:       nonce,
			Timestamp:   timestamp,
			Cards:       []constant.God{1, 2, 3, 4, 5, 6, 7, 8, 9},
			ChosenIndex: 4,
			SignedMsg:   signatureMsg,
		}

		playGame(playGameReq)

	}

}

func TestBatchWaitNewAccountPlayGame(t *testing.T) {
	for i := 0; i < 1; i++ {
		for j := 0; j < 100; j++ {
			newAccountPlayGame()
			//time.Sleep(time.Duration(500) * time.Millisecond)
		}
		//time.Sleep(time.Duration(1) * time.Minute)
	}
}

func TestGoBatchNewAccountPlayGame(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 50)
	for i := 0; i < 1000; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			newAccountPlayGame()
			<-ch
		}(i)
	}
	wg.Wait()

	time.Sleep(60 * time.Second)
}

func newAccountPlayGame() {
	keyTemp := generateEthAccount()
	keysArr := strings.Split(keyTemp, ",")
	address := keysArr[0]
	privateKeyStr := keysArr[1]

	nonce, _ := getNonce(address)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	cards := "[1 2 3 4 5 6 7 8 9]"
	chosenIndex := rand.Intn(9)

	msg := strings.ReplaceAll(constant.PlayGameSignMsg, "${address}", address)
	msg = strings.ReplaceAll(msg, "${nonce}", nonce)
	msg = strings.ReplaceAll(msg, "${timestamp}", timestamp)
	msg = strings.ReplaceAll(msg, "${cards}", cards)
	msg = strings.ReplaceAll(msg, "${chosen}", fmt.Sprintf("%d", chosenIndex))

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sign msg :", msg)
	data := []byte(msg)
	msgHash := crypto.Keccak256Hash(data)
	fmt.Println("msgHash", msgHash.Hex())
	signature, err := crypto.Sign(msgHash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	signatureMsg := hexutil.Encode(signature)
	fmt.Println("signature msg", signatureMsg)

	playGameReq := vo.PlayGameReq{
		Address:     address,
		Nonce:       nonce,
		Timestamp:   timestamp,
		Cards:       []constant.God{1, 2, 3, 4, 5, 6, 7, 8, 9},
		ChosenIndex: chosenIndex,
		SignedMsg:   signatureMsg,
	}
	playGame(playGameReq)
}
