package writer

import "os"

func WriteToFile(path string, data []byte) error {
	perms := os.FileMode(0o600)

	err := os.WriteFile(path, data, perms)
	if err != nil {
		return err
	}

	return os.Chmod(path, perms)
}
