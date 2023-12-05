package workfile

import (
	"encoding/gob"
	"os"
)

func SaveToGobFile(data interface{}, direct string) error {
	file, err := os.Create(direct)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromGobFile(result interface{}, direct string) error {
	file, err := os.Open(direct)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(result)
	if err != nil {
		return err
	}

	return nil
}
