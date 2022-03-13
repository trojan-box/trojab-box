package usecase

import (
	"errors"
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/pkg/jwt"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type AuthUseCase struct {
}

func newAuth() *AuthUseCase {
	return &AuthUseCase{}
}

func (u *AuthUseCase) VerifyTimestamp(loginTimestamp string) (bool, error) {
	loginTimestampInt, err := strconv.ParseInt(loginTimestamp, 10, 64)
	if err != nil {
		logger.WithError(err).Errorf("parse loginTimestamp occur error,loginTimestamp:%s", loginTimestamp)
		return false, err
	}
	loginTimestampTime := time.Unix(loginTimestampInt, 0)

	timeNow := time.Now()
	if loginTimestampTime.After(timeNow.Add(5 * time.Minute)) {
		logger.Infof("login timestamp after now 5 mins")
		return false, nil
	}
	if loginTimestampTime.Add(5 * time.Minute).Before(timeNow) {
		logger.Infof("login timestamp before now 5 mins")
		return false, nil
	}
	return true, nil
}
func (u *AuthUseCase) VerifyNonce(nonce string, address string) (bool, error) {
	log := logger.WithField("nonce", nonce).WithField("address", address)

	nonceUseCase := Svc.Nonce()
	nonceModal, err := nonceUseCase.FindByAddressAndNonce(address, nonce)
	if err != nil {
		log.WithError(err).Errorf("find nonce by address and nonceCode occur error")
		return false, err
	}
	if nonceModal == nil {
		log.Infof("not found nonce")
		return false, nil
	}
	if time.Unix(nonceModal.CreatedAt, 0).Add(5 * time.Minute).Before(time.Now()) {
		log.Infof("nonce create time before now 5 mins")
		return false, nil
	}
	return true, nil
}

func (u *AuthUseCase) VerifyLoginSignature(req vo.UserAuthReq) (bool, error) {

	msg := strings.ReplaceAll(constant.LoginSignMsg, "${address}", req.Address)
	msg = strings.ReplaceAll(msg, "${nonce}", req.Nonce)
	msg = strings.ReplaceAll(msg, "${timestamp}", req.Timestamp)
	msgBytes := []byte(msg)
	msgHash := crypto.Keccak256Hash(msgBytes)
	signature, err := hexutil.Decode(req.SignedMsg)
	if err != nil {
		logger.WithError(err).Errorf("decode signedmsg occur error")
		return false, err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(msgHash.Bytes(), signature)
	if err != nil {
		logger.WithError(err).Errorf("SigToPub signature occur error")
		return false, err
	}

	sigAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA).Hex()
	if sigAddress == req.Address {
		return true, nil
	} else {
		return false, nil
	}

}

func (u *AuthUseCase) GenerateToken(address string) (string, error) {
	return jwt.GenToken(address, []byte(app.Conf.JwtSecret))
}

func (u *AuthUseCase) VerifyLoginRequest(req vo.UserAuthReq) (bool, error) {
	result, err := u.VerifyNonce(req.Nonce, req.Address)
	if err != nil {
		logger.WithError(err).Errorf("verify nonce occur error,nonce:%s", req.Nonce)
		return false, errors.New("verify nonce incorrect")
	}
	if !result {
		return false, errors.New("verify nonce incorrect")
	}
	result, err = u.VerifyTimestamp(req.Timestamp)
	if err != nil {
		logger.WithError(err).Errorf("verify timestamp occur error,timestamp:%s", req.Timestamp)
		return false, errors.New("verify timestamp incorrect")
	}
	if !result {
		return false, errors.New("verify timestamp incorrect")
	}

	result, err = u.VerifyLoginSignature(req)
	if err != nil {
		logger.WithError(err).Errorf("verify login signature occur error")
		return false, errors.New("verify signature incorrect")
	}
	if !result {
		return false, errors.New("verify signature incorrect")
	}
	return true, nil

}

func (u *AuthUseCase) VerifyPlayGameSignature(req vo.PlayGameReq) (bool, error) {

	msg := strings.ReplaceAll(constant.PlayGameSignMsg, "${address}", req.Address)
	msg = strings.ReplaceAll(msg, "${nonce}", req.Nonce)
	msg = strings.ReplaceAll(msg, "${timestamp}", req.Timestamp)
	msg = strings.ReplaceAll(msg, "${cards}", fmt.Sprintf("%v", req.Cards))
	msg = strings.ReplaceAll(msg, "${chosen}", fmt.Sprintf("%v", req.ChosenIndex))
	logger.Infof("sign msg:%s", msg)
	msgBytes := []byte(msg)
	msgHash := crypto.Keccak256Hash(msgBytes)
	signature, err := hexutil.Decode(req.SignedMsg)
	if err != nil {
		logger.WithError(err).Errorf("decode signedmsg occur error")
		return false, err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(msgHash.Bytes(), signature)
	if err != nil {
		logger.WithError(err).Errorf("SigToPub signature occur error")
		return false, err
	}

	sigAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA).Hex()
	if sigAddress == req.Address {
		return true, nil
	} else {
		return false, nil
	}

}

func (u *AuthUseCase) VerifyPlayGameRequest(req vo.PlayGameReq) (bool, error) {
	result, err := u.VerifyNonce(req.Nonce, req.Address)
	if err != nil {
		logger.WithError(err).Errorf("verify nonce occur error,nonce:%s", req.Nonce)
		return false, errors.New("verify nonce incorrect")
	}
	if !result {
		return false, errors.New("verify nonce incorrect")
	}
	result, err = u.VerifyTimestamp(req.Timestamp)
	if err != nil {
		logger.WithError(err).Errorf("verify timestamp occur error,timestamp:%s", req.Timestamp)
		return false, errors.New("verify timestamp incorrect")
	}
	if !result {
		return false, errors.New("verify timestamp incorrect")
	}

	result, err = u.VerifyPlayGameSignature(req)
	if err != nil {
		logger.WithError(err).Errorf("verify playGame signature occur error")
		return false, errors.New("verify signature incorrect")
	}
	if !result {
		return false, errors.New("verify signature incorrect")
	}
	return true, nil

}

func (u *AuthUseCase) VerifyWithdrawBonusApplySignature(req vo.WithdrawBonusApplyReq) (bool, error) {

	msg := strings.ReplaceAll(constant.WithdrawBonusApplySignMsg, "${address}", req.Address)
	msg = strings.ReplaceAll(msg, "${nonce}", req.Nonce)
	msg = strings.ReplaceAll(msg, "${timestamp}", req.Timestamp)
	msg = strings.ReplaceAll(msg, "${bonus}", fmt.Sprintf("%v", req.Bonus))
	logger.Infof("sign msg:%s", msg)
	msgBytes := []byte(msg)
	msgHash := crypto.Keccak256Hash(msgBytes)
	signature, err := hexutil.Decode(req.SignedMsg)
	if err != nil {
		logger.WithError(err).Errorf("decode signedmsg occur error")
		return false, err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(msgHash.Bytes(), signature)
	if err != nil {
		logger.WithError(err).Errorf("SigToPub signature occur error")
		return false, err
	}

	sigAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA).Hex()
	if sigAddress == req.Address {
		return true, nil
	} else {
		return false, nil
	}

}

func (u *AuthUseCase) VerifyWithdrawBonusApplyRequest(req vo.WithdrawBonusApplyReq) (bool, error) {
	result, err := u.VerifyNonce(req.Nonce, req.Address)
	if err != nil {
		logger.WithError(err).Errorf("verify nonce occur error,nonce:%s", req.Nonce)
		return false, errors.New("verify nonce incorrect")
	}
	if !result {
		return false, errors.New("verify nonce incorrect")
	}
	result, err = u.VerifyTimestamp(req.Timestamp)
	if err != nil {
		logger.WithError(err).Errorf("verify timestamp occur error,timestamp:%s", req.Timestamp)
		return false, errors.New("verify timestamp incorrect")
	}
	if !result {
		return false, errors.New("verify timestamp incorrect")
	}

	result, err = u.VerifyWithdrawBonusApplySignature(req)
	if err != nil {
		logger.WithError(err).Errorf("verify WithdrawBonusApply signature occur error")
		return false, errors.New("verify signature incorrect")
	}
	if !result {
		return false, errors.New("verify signature incorrect")
	}
	return true, nil

}
