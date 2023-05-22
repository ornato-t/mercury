package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

const GIT = "./tools/git/bin/git"
const PANDOC = "./tools/pandoc"
const REPO = "./paolo-sernini/src/pages/scrittura/post/"

func main() {
	if err := downloadRepo(); err != nil {
		log.Panic(err)
	}

	fileList, err := findDocs()
	if err != nil {
		log.Panic(err)
	}

	for _, file := range fileList {
		if err := convert(file); err != nil {
			log.Panic(err)
		}
	}

	if err := deleteRepo(); err != nil {
		log.Panic(err)
	}
}

//Return a slice of all the .docx documents in the calling directory
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

//Convert a slice of .docx documents to .md ones with the same name
func convert(fileName string) error {
	cmd := exec.Command(PANDOC, "-f", "docx", "-t", "markdown", fileName+".docx", "-o", REPO+fileName+".md")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

//Clone a git repo
func downloadRepo() error {
	cmd := exec.Command(GIT, "clone", "https://github.com/ornato-t/paolo-sernini")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

//Reads a github token from the .env file
func getToken() (string, error){
    // Load the .env file
    err := godotenv.Load()
    if err != nil {
		return "", err
    }

    // Access a value using os.Getenv
    value := os.Getenv("TOKEN")
    return value, err
}

//Delete the repo once all changes have been made
func deleteRepo() error {
	// Set the directory to delete
	dir := "paolo-sernini"

	// Delete the directory and its contents
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}
