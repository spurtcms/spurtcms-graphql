package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func UploadImageToAdminLocal(b64Data, fileName, adminUrl string) (string, error) {

	client := http.Client{}

	paramSetter := url.Values{}

	paramSetter.Add("imgFile", b64Data)

	paramSetter.Add("imgFileName", fileName)

	method := "POST"

	uploadReq, err := http.NewRequest(method, adminUrl, bytes.NewBufferString(paramSetter.Encode()))

	if err != nil {

		return "", err
	}

	response, err := client.Do(uploadReq)

	if err != nil {

		return "", err
	}

	if response.StatusCode != 200{

		fmt.Println("response status",response.StatusCode)

		return  "", errors.New("failed to store file in admin panel local")
	}

	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)

	if err != nil {

		return "", err
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(responseBytes, &responseData)

	if err != nil {

		return "", err
	}

	path, ok := responseData["StoragePath"].(string)

	fmt.Printf("path %v\n", path)

	if !ok {

		return "", errors.New("failed to get storage path")
	}

	return path, nil

}
