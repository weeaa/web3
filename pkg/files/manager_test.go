package files

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

const perm = 0666

const (
	tempFilePathCSV  = "temp.csv"
	tempFilePathJSON = "temp.json"
)

var expectedValuesCSV = []string{"bob", "foo", "43", "FR"}

type csvContent struct {
	FirstName  string `csv:"firstName"`
	LastName   string `csv:"lastName"`
	Age        int    `csv:"age"`
	CountryISO string `csv:"countryISO"`
}

type jsonContent struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Age        int    `json:"age"`
	CountryISO string `json:"countryISO"`
}

func TestCSV(t *testing.T) {
	CreateCSV(tempFilePathCSV, [][]string{
		{
			"firstName",
			"lastName",
			"age",
			"countryISO",
		},
	})

	defer os.Remove(tempFilePathCSV)

	if err := AppendCSV(tempFilePathCSV, expectedValuesCSV); err != nil {
		assert.Error(t, err)
	}

	rows, err := ReadCSV[csvContent](tempFilePathCSV)
	if err != nil {
		assert.Error(t, err)
	}

	if rows[0].CountryISO != "FR" || rows[0].FirstName != "bob" || rows[0].LastName != "foo" || rows[0].Age != 43 {
		assert.Errorf(t, errors.New("data not matching"), "expected %v, got %v", expectedValuesCSV, rows)
	}

	assert.NoError(t, nil)
}

func TestJSON(t *testing.T) {
	var err error

	values := jsonContent{
		FirstName:  expectedValuesCSV[0],
		LastName:   expectedValuesCSV[1],
		Age:        43,
		CountryISO: expectedValuesCSV[3],
	}

	if err = WriteJSON(tempFilePathJSON, values); err != nil {
		assert.Error(t, err)
	}

	defer os.Remove(tempFilePathJSON)

	data, err := ReadJSON[jsonContent](tempFilePathJSON)
	if err != nil {
		assert.Error(t, err)
	}

	if data.FirstName != "bob" {
		assert.Errorf(t, errors.New("data not matching"), "expected %v, got %v", expectedValuesCSV, data)
	}

	assert.NoError(t, nil)
}

func TestYAML(t *testing.T) {
	var err error
	var buf bytes.Buffer

	values := map[string]any{
		"foo": true,
		"bar": 61.77,
		"baz": "hey",
	}

	if err = yaml.NewEncoder(&buf).Encode(values); err != nil {
		assert.Error(t, err)
	}

}

func TestCreateFolder(t *testing.T) {
	tempFolderPath := "tempfolder"

	CreateFolder(tempFolderPath)
	defer os.Remove(tempFolderPath)

	_, err := os.Stat(tempFolderPath)
	if err != nil {
		assert.Errorf(t, err, "failed to create test folder")
	}

	assert.NoError(t, nil)
}

func TestCreateFile(t *testing.T) {
	tempFilePath := "tempfile.txt"

	CreateFile(tempFilePath)
	defer os.Remove(tempFilePath)

	_, err := os.Stat(tempFilePath)
	if err != nil {
		assert.Errorf(t, err, "failed to create test file")
	}

	assert.NoError(t, nil)
}
