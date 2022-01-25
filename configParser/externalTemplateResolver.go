package configParser

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

const templateStorage = "./templates"

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// source : https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ResolveExternalTemplate(tpl importTemplate) error {
	fileName := fmt.Sprintf("%s.yml", tpl.Name)
	filePath := path.Join(templateStorage, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create folder
		err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("fail to create the template directory : %v", err)
		}

		// Import file
		fileUrl := tpl.Url
		err = DownloadFile(filePath, fileUrl)
		if err != nil {
			return fmt.Errorf("fail to import the template for %s : %v", tpl.Name, err)
		}
	}
	return nil
}
