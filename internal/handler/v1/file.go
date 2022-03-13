package v1

import (
	"context"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/thirdparty"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
	ipfsFiles "github.com/ipfs/go-ipfs-files"
	logger "github.com/sirupsen/logrus"
	"io"
)

// UploadFile godoc
// @Summary UploadFile
// @Description uploadFile
// @Tags file
// @Produce json
// @Param Authorization header string true "accessToken"
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {object} vo.Response{data=string}
// @Router /file/upload [post]
func UploadFile(ctx *gin.Context) {
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	file, _ := ctx.FormFile("file")
	src, err := file.Open()
	if err != nil {
		logger.WithError(err).Error("open file occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	defer src.Close()

	buffer, err := io.ReadAll(src)
	if err != nil {
		logger.WithError(err).Error("read all image err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		logger.WithError(err).Error("convert to webp occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 20, Height: 1500, Embed: true})
	if err != nil {
		logger.WithError(err).Error("convert to webp occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	p, err := thirdparty.IpfsClient.Unixfs().Add(context.Background(), ipfsFiles.NewBytesFile(processed))
	if err != nil {
		logger.WithError(err).Errorln("add file to ipfs occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	logger.Infof("Data successfully stored in IPFS: %v\n", p.Cid().String())
	ipfsHash := p.Cid().String()
	fileUseCase := usecase.Svc.File()
	link, err := fileUseCase.Save(address, ipfsHash)
	if err != nil {
		logger.WithError(err).Errorf("save file occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, link)
}
