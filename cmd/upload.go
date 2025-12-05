package cmd

import (
	// "fmt"
	"github.com/spf13/cobra"
	// "github.com/udaycmd/rdv/utils"
	// "google.golang.org/api/drive/v3"
	// "os"
	// "path/filepath"
)

var uploadCmd = &cobra.Command{
	Use:                   "upload [-f file | -d directory]",
	Short:                 "uploads a file or directory to the drive",
	Long:                  "uploads a file or directory to the drive",
	Aliases:               []string{"set"},
	DisableFlagsInUseLine: true,
	// PreRun: func(cmd *cobra.Command, args []string) {
	// 	err := auth()
	// 	if err != nil {
	// 		utils.ExitOnError("Authorization error: " + err.Error())
	// 	}

	// 	if (file == "" && dir == "") || (file != "" && dir != "") {
	// 		utils.ExitOnError("You must provide either -f (file) or -d (directory), but not both")
	// 	}

	// 	if file == "" {
	// 		isDir = true
	// 	} else {
	// 		isDir = false
	// 	}

	// 	if parent != "" {
	// 		q := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and 'root' in parents and trashed=false", parent)
	// 		res, err := srv.Files.List().Q(q).Fields("files(id, name)").Do()
	// 		if err != nil || len(res.Files) == 0 {
	// 			utils.ExitOnError("Unable to find the parent folder: " + parent)
	// 		}

	// 		parentID = res.Files[0].Id
	// 	}

	// 	if err = utils.ClearScreen(); err != nil {
	// 		utils.ExitOnError(err.Error())
	// 	}
	// },
	// Run: func(cmd *cobra.Command, args []string) {
	// 	var path string
	// 	var message string

	// 	if isDir {
	// 		path = dir
	// 		message = fmt.Sprintf("Uploading directory %s", filepath.Base(dir))
	// 	} else {
	// 		path = file
	// 		message = fmt.Sprintf("Uploading file %s", filepath.Base(file))
	// 	}

	// 	_, err := utils.Spinner(func() (struct{}, error) {
	// 		return struct{}{}, upload(path, isDir, parentID)
	// 	}, message)

	// 	if err != nil {
	// 		utils.ExitOnError("Unable to complete the upload operation: " + err.Error())
	// 	}

	// 	utils.ExitOnSuccess("Upload operation for " + filepath.Base(path) + " successfully completed!")
	// },
}

func init() {
	// uploadCmd.Flags().StringVarP(&file, "file", "f", "", "file to upload")
	// uploadCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to upload")
	// uploadCmd.Flags().StringVarP(&parent, "parent", "p", "", "finds the first folder named as specified and uploads the file or directory inside it")
	rootCmd.AddCommand(uploadCmd)
}

// func upload(path string, isDir bool, parentDirId string) error {
// 	var err error

// 	if isDir {
// 		err = uploadDriveFolderRec(path, parentDirId)
// 	} else {
// 		err = uploadDriveFile(path, parentDirId)
// 	}

// 	return err
// }

// func uploadDriveFile(file string, parentId string) error {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	filename := filepath.Base(filepath.Clean(file))
// 	if filename == "/" || filename == "\\" {
// 		return fmt.Errorf("invalid file name or path")
// 	}

// 	driveFile := &drive.File{Name: filename}
// 	if parentId != "" {
// 		driveFile.Parents = []string{parentId}
// 	}

// 	_, err = srv.Files.Create(driveFile).Media(f).Do()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func createDriveFolder(dir string, parentId string) (string, error) {
// 	dirname := filepath.Base(filepath.Clean(dir))

// 	driveFolder := &drive.File{
// 		Name:     dirname,
// 		MimeType: "application/vnd.google-apps.folder",
// 	}
// 	if parentId != "" {
// 		driveFolder.Parents = []string{parentId}
// 	}
// 	createdFolder, err := srv.Files.Create(driveFolder).Do()

// 	if err != nil {
// 		return "", err
// 	}

// 	return createdFolder.Id, nil
// }

// func uploadDriveFolderRec(dir string, parentId string) error {
// 	entries, err := os.ReadDir(dir)
// 	if err != nil {
// 		return err
// 	}

// 	folderId, err := createDriveFolder(dir, parentId)
// 	if err != nil {
// 		return err
// 	}

// 	for _, entry := range entries {
// 		fullpath := filepath.Join(dir, entry.Name())
// 		if entry.IsDir() {
// 			err = uploadDriveFolderRec(fullpath, folderId)
// 		} else {
// 			err = uploadDriveFile(fullpath, folderId)
// 		}
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
