package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func main() {
	fileList, err := findDocs()
	if err != nil {
		log.Panic(err)
	}

	for _, file := range fileList {
		if err := convert(file); err != nil {
			log.Panic(err)
		}
	}

}

func findDocs() ([]string, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var docxFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".docx") {
			purged := strings.Split(file.Name(), ".")

			docxFiles = append(docxFiles, purged[0])
		}
	}

	return docxFiles, nil
}

func convert(fileName string) error {
	cmd := exec.Command("./tools/pandoc", "-f", "docx", "-t", "markdown", fileName+".docx", "-o", fileName+".md")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
