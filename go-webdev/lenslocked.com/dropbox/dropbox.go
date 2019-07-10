package dropbox

import (
	"fmt"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	dbxfiles "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type File struct {
	Name string
	Path string
}

type Folder struct {
	Name string
	Path string
}

func List(accessToken, path string) ([]Folder, []File, error) {
	dbxClient := dbxfiles.New(dropbox.Config{Token: accessToken})
	res, err := dbxClient.ListFolder(&dbxfiles.ListFolderArg{Path: path})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list path %s: %s", path, err)
	}

	files := []File{}
	folders := []Folder{}
	for _, e := range res.Entries {
		switch v := e.(type) {
		case *dbxfiles.FileMetadata:
			files = append(files, File{Name: v.Id, Path: v.PathLower})
		case *dbxfiles.FolderMetadata:
			folders = append(folders, Folder{Name: v.Id, Path: v.PathLower})
		}
	}

	return folders, files, nil
}
