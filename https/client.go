package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// 读取自签名证书
	cert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		fmt.Println("Failed to read cert file:", err)
		return
	}

	// 创建证书池并添加自签名证书
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(cert)

	// 创建 HTTP 客户端，使用自定义证书
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: rootCAs},
		},
	}

	// 访问 HTTPS 服务器
	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}
