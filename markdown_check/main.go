package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const referenceDir string = "./content/references/"
const slipboxDir string = "./content/slipbox/"
const fleetingDir string = "./content/fleeting/"

func checkSyncConflictsAllDirectories() error {
	filesWithConflicts := make([]string, 0)

	var dirNames [3]string = [3]string{fleetingDir, referenceDir, slipboxDir}
	for _, dirName := range dirNames {
		conflictFiles, err := checkSyncConflicts(dirName)
		if err != nil {
			return err
		}
		if len(conflictFiles) > 0 {
			filesWithConflicts = append(filesWithConflicts, conflictFiles...)
		}
	}

	if len(filesWithConflicts) > 0 {
		fmt.Println("Files with sync conflicts:")
		for _, file := range filesWithConflicts {
			fmt.Println(file)
		}
		return errors.New("sync conflicts are not allowed")
	}

	return nil
}

func checkSyncConflicts(directoryPath string) ([]string, error) {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	filesWithConflicts := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "sync-conflict") {
			filesWithConflicts = append(filesWithConflicts, fmt.Sprintf("%s%s", directoryPath, file.Name()))
		}
	}
	return filesWithConflicts, nil
}

func reportUnusedReferences(directoryPath string) error {
	r, err := regexp.Compile("title: .+[^\n]")
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	const matchString string = "used_yet: false\n---\n"

	unusedReferenceFiles := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileContents, err := os.ReadFile(fmt.Sprintf("%s%s", directoryPath, file.Name()))
		if err != nil {
			return err
		}

		if strings.Contains(string(fileContents), matchString) {
			title := strings.TrimLeft(r.FindString(string(fileContents)), "title:")
			unusedReferenceFiles = append(unusedReferenceFiles, fmt.Sprintf("%s - %s", file.Name(), title))
		}
	}

	if len(unusedReferenceFiles) > 0 {
		fmt.Println("Unused reference files:")
		for _, file := range unusedReferenceFiles {
			fmt.Println(file)
		}
	}
	return nil
}

func main() {
	if err := checkSyncConflictsAllDirectories(); err != nil {
		log.Fatal(err)
	}

	if err := reportUnusedReferences(referenceDir); err != nil {
		log.Fatal(err)
	}
}
