package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func UploadImageToAdminLocal(b64Data,fileName,adminUrl string,) (string, error){

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

	if err != nil || response.StatusCode != 200 {

		return "", err
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

	if !ok {

		return "", errors.New("failed to get storage path in admin panel")
	}

	return path, nil

}

