package usecase

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type FileUseCase struct {
	FileRepo repository.FileRepo
}

func newFile(svc *useCase) *FileUseCase {
	return &FileUseCase{
		FileRepo: repository.NewFile(svc.db),
	}
}

func (u *FileUseCase) Save(address, ipfsHash string) (string, error) {

	link := "https://ipfs.io/ipfs/" + ipfsHash
	file := model.UploadFile{
		Address:  address,
		IpfsHash: ipfsHash,
		Link:     link,
	}

	_, err := u.FileRepo.Save(file)
	if err != nil {
		logger.WithError(err).Errorln("save file to db occur err")
		return "", err
	}
	return link, nil
}
