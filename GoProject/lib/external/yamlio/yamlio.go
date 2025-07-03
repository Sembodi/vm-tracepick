package yamlio

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Encode and write object as yaml file.
func WriteYaml[T any](path string, object T) error {
	var (
		buf []byte
		err error
	)
	if buf, err = yaml.Marshal(object); err != nil {
		return err
	}
	if err = os.WriteFile(path, buf, 0644); err != nil {
		return err
	}
	return nil
}

// Read and decode yaml file to object.
// Internally it uses a pointer of object so beware!
func ReadYaml[T any](path string) (T, error) {
	var (
		buf    []byte
		err    error
		object T
	)

	if buf, err = os.ReadFile(path); err != nil {
		return object, err
	}
	if err = yaml.Unmarshal(buf, &object); err != nil {
		return object, err
	}
	return object, nil
}
