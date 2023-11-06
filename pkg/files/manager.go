package files

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/jszwec/csvutil"
	"github.com/weeaa/nft/pkg/utils"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadCSV[T any](filePath string) ([]T, error) {
	var err error
	var Rows []T
	fileContent, err := os.ReadFile(utils.ExecPath + filePath)
	if err != nil {
		return nil, err
	}
	if err = csvutil.Unmarshal(fileContent, &Rows); err != nil {
		return nil, err
	}
	return Rows, nil
}

func AppendCSV(filePath string, params []string) error {
	f, _ := os.OpenFile(utils.ExecPath+filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	w := csv.NewWriter(f)
	if err := w.Write(params); err != nil {
		return err
	}
	w.Flush()
	return f.Close()
}

func ReadJSON[T any](filePath string) (T, error) {
	var data T
	file, err := os.ReadFile(filePath)
	if err != nil {
		return data, err
	}
	if err = json.Unmarshal(file, &data); err != nil {
		return data, err
	}
	return data, nil
}

func WriteJSON(filePath string, data any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0777)
}

func CreateCSV(filePath string, keys [][]string) {
	if _, err := os.Stat(utils.ExecPath + filePath); err != nil {
		f, _ := os.OpenFile(utils.ExecPath+filePath, os.O_CREATE|os.O_WRONLY, 0666)
		csv.NewWriter(f).WriteAll(keys)
		f.Close()
	}
}

func CreateFolder(folderPath string) {
	if _, err := os.Stat(utils.ExecPath + folderPath); err != nil {
		os.MkdirAll(utils.ExecPath+folderPath, 0777)
	}
}

func CreateFile(filePath string) {
	if _, err := os.Stat(utils.ExecPath + filePath); err != nil {
		f, _ := os.OpenFile(utils.ExecPath+filePath, os.O_CREATE|os.O_WRONLY, 0666)
		f.Close()
	}
}

func CreateYAML[T any](filePath string, dataEncoded T) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	err = yaml.NewEncoder(file).Encode(&dataEncoded)
	return err
}

func ReadYAML[T any](filePath string) (T, error) {
	var dataDecoded T

	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return dataDecoded, err
	}

	err = yaml.NewDecoder(f).Decode(&dataDecoded)
	return dataDecoded, err
}
