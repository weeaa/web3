package files

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jszwec/csvutil"
	"github.com/weeaa/nft/pkg/utils"
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
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0777)
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
