package configParser

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// DownloadFile will download an url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// source : https://golangcode.com/download-a-file-from-a-url/
func downloadFile(dst io.Writer, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read imported file : %v", err)
	}
	return nil
}

func resolveExternalTemplate(tpl importTemplate) (string, error) {

	// Create the temp folder
	if err := os.MkdirAll("/tmp/workigo", os.ModePerm); err != nil {
		return "", err
	}

	// Create the file
	outfile, err := ioutil.TempFile("/tmp/workigo", "external-wf-")
	if err != nil {
		return "", err
	}
	defer outfile.Close()

	err = downloadFile(outfile, tpl.Url)
	if err != nil {
		return "", fmt.Errorf("fail to fetch external template : %v", err)
	}
	return outfile.Name(), nil
}

func ResolveExternalTemplates(templates []importTemplate) (map[string]string, error) {
	importedTemplates := make(map[string]string)

	if len(templates) > 0 {
		log.Println("[WARNING] : Using external templates can expose your system to several risks.")
		for _, externalTemplate := range templates {
			templateFile, err := resolveExternalTemplate(externalTemplate)
			if err != nil {
				return nil, fmt.Errorf("fail to import template : %v", err)
			}
			importedTemplates[externalTemplate.Name] = templateFile
		}
	}

	return importedTemplates, nil
}
