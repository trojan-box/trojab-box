package thirdparty

import (
	"encoding/base64"
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/app"
	ipfsApi "github.com/ipfs/go-ipfs-http-client" // v0.1.0
	"net/http"
	"os"
)

var IpfsClient *ipfsApi.HttpApi

func InitIpfs(config app.Config) {

	ipfsConfig := config.Ipfs

	httpClient := &http.Client{}
	httpApi, err := ipfsApi.NewURLApiWithClient(ipfsConfig.Url, httpClient)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	httpApi.Headers.Add("Authorization", "Basic "+basicAuth(ipfsConfig.ProjectId, ipfsConfig.ProjectSecret))
	IpfsClient = httpApi
	//
	//content := strings.NewReader("Infura IPFS - Getting started demo.")
	//p, err := httpApi.Unixfs().Add(context.Background(), ipfsFiles.NewReaderFile(content))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//fmt.Printf("Data successfully stored in IPFS: %v\n", p.Cid().String())
}

func basicAuth(projectId, projectSecret string) string {
	auth := projectId + ":" + projectSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
