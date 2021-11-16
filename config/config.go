package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"text/template"
	"time"
)

type FileStruct struct {
	Name     string
	Env      map[string]string
	Commands []string
	Image    string
	Workdir  string
}

type TemplateFile struct {
	Data           *FileStruct
	TemplateName   string
	ConfigTemplate *template.Template
	name           string
}

func ReadConfFile(filename string) *TemplateFile {

	configContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Fail to open file %s : %v", filename, err)
	}

	cnf := &FileStruct{}

	err = yaml.Unmarshal(configContent, &cnf)
	if err != nil {
		log.Fatalf("Fail to load file %v", err)
	}

	templateName := strconv.FormatInt(time.Now().Unix(), 16)

	tmpl := template.New(templateName)

	return &TemplateFile{
		Data:           cnf,
		TemplateName:   templateName,
		ConfigTemplate: tmpl,
	}
}

func (config *TemplateFile) ResolvePattern(value string) (string, error) {

	tmpl, err := template.New("config").Parse(value)

	buf := &bytes.Buffer{}

	err = tmpl.ExecuteTemplate(buf, "config", config.Data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (config *TemplateFile) Name() (string, error) {
	if config.name == "" {
		name, err := config.ResolvePattern(config.Data.Name)
		if err != nil {
			return "", err
		}
		config.name = name
	}
	return config.name, nil
}
