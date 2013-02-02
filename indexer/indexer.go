// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indexer

import (
	"andyk/docs/model"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Index(repositoryPath string) map[int]model.Document {

	// check if the supplied repository path is set
	if strings.Trim(repositoryPath, " ") == "" {
		fmt.Print("The repository path cannot be null or empty.")
		return nil
	}

	// check if the supplied repository path exists
	if _, err := os.Stat(repositoryPath); err != nil {
		switch {
		case os.IsNotExist(err):
			fmt.Printf("The supplied repository path `%v` does not exist.", repositoryPath)
		default:
			fmt.Printf("An error occured while trying to access the supplied repository path `%v`.", repositoryPath)
		}

		return nil
	}

	items := FindAllRepositoryItems(repositoryPath)
	fmt.Printf("%#v", items)

	return nil
}

func FindAllRepositoryItems(repositoryPath string) map[int]*model.RepositoryItem {
	items := make(map[int]*model.RepositoryItem)
	itemIndex := 0

	// repository walker
	repositoryWalker := func(path string, _ os.FileInfo, _ error) error {

		fileName := filepath.Base(path)

		// check if the file a document
		isRepositoryItem := strings.EqualFold(strings.ToLower(fileName), "notes.md")
		if !isRepositoryItem {
			return nil
		}

		items[itemIndex] = model.NewRepositoryItem(path)
		itemIndex++

		return nil
	}

	// index the repository
	walkRepositoryError := filepath.Walk(repositoryPath, repositoryWalker)
	if walkRepositoryError != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v", repositoryPath, walkRepositoryError)
	}

	return items
}
