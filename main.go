package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const GIT = "./tools/git/bin/git"
const PANDOC = "./tools/pandoc"
const MD_FOLDER = "./paolo-sernini/src/pages/scrittura/post/"
const JSON = "./paolo-sernini/src/pages/scrittura/posts.json"
const REPO = "./paolo-sernini"
const GIT_FROM_REPO = "./../tools/git/bin/git"

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

		date := time.Now()

		if err := addHeading(file, date); err != nil {
			log.Panic(err)
		}

		if err := editJSON(file, date); err != nil {
			log.Panic(err)
		}
	}

	if err := commit(); err != nil {
		log.Panic(err)
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
	cmd := exec.Command(PANDOC, "-f", "docx", "-t", "markdown", fileName+".docx", "-o", MD_FOLDER+fileName+".md")
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
func getToken() (string, error) {
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

//Add the Astro markdown heading at the top of a markdown file
func addHeading(fileName string, date time.Time) error {
	path := MD_FOLDER + fileName + ".md"
	// Read the file contents into a string
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)

	// Add new text at the top
	newText := "---\n\nlayout: ../../../layouts/Post.astro\n\ntitle: \"" + fileName + "\"\n\ndate: \"" + date.Format(time.RFC3339) + "\"\n\n---\n\n"
	content = newText + content

	// Write the modified content back to the file
	err = ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func editJSON(fileName string, date time.Time) error {
	// Define a struct type to represent an entry
	type Entry struct {
		Title string `json:"title"`
		Date  string `json:"date"`
	}

    fileBytes, err := os.ReadFile(JSON)
    if err != nil {
        return err
    }

	var entries []Entry
	err = json.Unmarshal(fileBytes, &entries)
	if err != nil {
		return err
	}

	found := false
	// Iterate over the entries and read their values
	for _, entry := range entries {
		if fileName == entry.Title {
			entry.Date = date.Format(time.RFC3339)
			found = true
		}
	}

	if !found {
		// Add a new entry to the array
		newEntry := Entry{Title: fileName, Date: date.Format(time.RFC3339)}
		entries = append(entries, newEntry)
	}

	// write updated data back to file
    fileBytes, err = json.Marshal(entries)
    if err != nil {
		return err
    }

    err = ioutil.WriteFile(JSON, fileBytes, 0644)
    if err != nil {
		return err
    }

	return nil
}

func commit() error {
	// set the working directory to the subfolder
    cmd := exec.Command(GIT_FROM_REPO, "config", "--global", "user.name", "Paolo Sernini")
    cmd.Dir = REPO
    err := cmd.Run()
    if err != nil {
        return err
    }

    cmd = exec.Command(GIT_FROM_REPO, "config", "--global", "user.email", "tommy.ornato@gmail.com")
    cmd.Dir = REPO
    err = cmd.Run()
    if err != nil {
        return err
    }

    // run git add
    cmd = exec.Command(GIT_FROM_REPO, "add", ".")
    cmd.Dir = REPO
    err = cmd.Run()
    if err != nil {
        return err
    }

    // run git commit
    cmd = exec.Command(GIT_FROM_REPO, "commit", "-m", "Updating posts")
    cmd.Dir = REPO
    err = cmd.Run()
    if err != nil {
        return err
    }

	token, err := getToken()
	if err != nil {
		return err
	}

    // set up authentication with GitHub token
    cmd = exec.Command(GIT_FROM_REPO, "remote", "set-url", "origin", token+"@github.com/ornato-t/paolo-sernini.git")
    cmd.Dir = REPO
    err = cmd.Run()
    if err != nil {
        return err
    }

	//Crashes sometimes after this point
	
	/*
		fatal: 'TOKEN@github.com/ornato-t/paolo-sernini.git' does not appear to be a git repository
		fatal: Could not read from remote repository.

		Please make sure you have the correct access rights
		and the repository exists.
	*/

	// run git push
    cmd = exec.Command(GIT_FROM_REPO, "push")
    cmd.Dir = REPO
    err = cmd.Run()
    if err != nil {
        return err
    }

	return nil
}
