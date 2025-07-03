package pathparse

import (
  "fmt"
  "os"
  "strings"
  "path/filepath"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/owncheckers"
  "project1/tracepick/lib/fileparsers/lineparse"
)

func insertPath(root *owntypes.Folder, path string) {
  parts := strings.Split(strings.TrimPrefix(path, "/"), "/") // Split into folder names
  //parts = parts[:len(parts)-1]
	current := root
	for _, part := range parts {
		if current.Folders == nil {
			current.Folders = make(map[string]*owntypes.Folder)
		}

		// If folder doesn't exist, create it
		if _, exists := current.Folders[part]; !exists {
			current.Folders[part] = &owntypes.Folder{Name: part, Folders: make(map[string]*owntypes.Folder)}
		}

		// Move to the next level
		current = current.Folders[part]
	}
}

func BuildFolderStructure(paths owntypes.StringSet) *owntypes.Folder {
	root := &owntypes.Folder{Name: "/", Folders: make(map[string]*owntypes.Folder)}

	for path, _ := range paths {
		insertPath(root, path)
	}

	return root
}

func PrintFolderStructure(folder *owntypes.Folder, indent string) {
	fmt.Println(indent + folder.Name)
	for _, subfolder := range folder.Folders {
		PrintFolderStructure(subfolder, indent+"  ")
	}
}

func RecurseFolders(currentPath string, currentFolder *owntypes.Folder, accumulator []string, protectedFolders owntypes.StringSet) []string {
  splittedPath := currentPath //save the path from which multiple folders may be included
  for name, folder := range currentFolder.Folders {
    var addPaths []string = []string{}
    currentPath = splittedPath + "/" + name

    if protectedFolders[currentPath] {
      addPaths = RecurseFolders(currentPath, folder, accumulator, protectedFolders)
    } else {
      addPaths = []string{currentPath}
    }
    accumulator = append(accumulator, addPaths...)
  }
  return accumulator
}

func GetMinimalPaths(folderTree *owntypes.Folder) []string {
    var (
      programName string = "GetMinimalPaths"
      protectedFoldersFile string = "defaultdata/helperFiles/protected_paths"

      skipDirs []string = []string{"/sys", "/dev", "/tmp", "/boot"}
      includeDirs []string = []string{}

      result []string
      protectedFolders owntypes.StringSet
      err error
    )

    if protectedFolders, err = lineparse.ReadFileLines(programName, protectedFoldersFile); err != nil {
      fmt.Println(programName, ": Reading '", protectedFoldersFile, "' unsuccessful. No MinimalPaths returned.")
      return nil
    }

    for name1, folder1 := range folderTree.Folders {
        if !owncheckers.StringPartOfSliceElem(name1, skipDirs) {
          for name2, folder2 := range folder1.Folders {
            var addPaths []string = []string{}
            currentPath := "/" + name1 + "/" + name2
            currentFolder := folder2

            // fmt.Println(currentPath, "is protected:", protectedFolders[currentPath])
            // fmt.Println("kid folders: ", currentFolder.Folders, ". Length: ", len(currentFolder.Folders))
            if len(folder2.Folders) == 0 {
              result = append(result, currentPath)
              continue
            }

            if protectedFolders[currentPath] {
              addPaths = RecurseFolders(currentPath, currentFolder, []string{}, protectedFolders)
              } else {
                addPaths = []string{currentPath}
              }
              result = append(result, addPaths...)
            }
        }
    }

    return append(result, includeDirs...)
}

func GetMinimalPathsMap(folderProfiles map[string]*owntypes.Folder) map[string][]string {
  result := make(map[string][]string)
  for item, folderTree := range folderProfiles {
    result[item] = GetMinimalPaths(folderTree)
  }
  return result
}

func PrintMinimalFolderList(name string, list []string) {
	fmt.Println(name,": ")
	for _, path := range list {
		fmt.Println("- ", path)
	}
}

func GetFolderProfiles(path string) map[string]*owntypes.Folder {
  var (
    programName string = "GetFolderProfiles"
    linesFromFiles map[string]owntypes.StringSet

    result map[string]*owntypes.Folder = make(map[string]*owntypes.Folder)
  )

  linesFromFiles = lineparse.LineSetsFromFiles(programName, path)

  for item, paths := range linesFromFiles {
    result[item] = BuildFolderStructure(paths)
    // // Print the structured folder tree
    //
    // fmt.Println(item,":")
    // folder := BuildFolderStructure(paths)
    // PrintFolderStructure(folder, "  ")
  }
  return result
}

// BuildFolderStructureFromFilesystem builds a Folder structure from an actual filesystem path
func BuildFolderStructureFromFilesystem(rootPath string) (*owntypes.Folder, error) {
	root := &owntypes.Folder{
		Name:    filepath.Base(rootPath),
		Folders: make(map[string]*owntypes.Folder),
	}

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // Propagate the error up
		}

		if path == rootPath {
			// Skip the root path itself; it's already added
			return nil
		}

		if d.IsDir() {
			relativePath := strings.TrimPrefix(path, rootPath)
			relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator)) // Ensure no leading '/'
			insertPath(root, relativePath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return root, nil
}

// // Example function to show runtime folder profile
// func ShowFolderStructure() {
//   folder, _ := BuildFolderStructureFromFilesystem("artifacts/output/containers/nginx/myrootfs-nginx")
//   fmt.Println("runtime profile:")
//   PrintFolderStructure(folder, " ")
// }
