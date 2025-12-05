package cmd

import (
	// "fmt"
	"github.com/spf13/cobra"
	// "github.com/udaycmd/rdv/utils"
	// "google.golang.org/api/drive/v3"
	// "io"
	// "os"
	// "path/filepath"
)

var fetchCmd = &cobra.Command{
	Use:                   "fetch [-f file | -d directory] [-o outputPath]",
	Short:                 "fetch a file or directory",
	Long:                  "fetch a file or directory",
	DisableFlagsInUseLine: true,
	Aliases:               []string{"get"},
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

	// 	if err = utils.ClearScreen(); err != nil {
	// 		utils.ExitOnError(err.Error())
	// 	}
	// },
	// Run: func(cmd *cobra.Command, args []string) {
	// 	var message string
	// 	var driveFName string

	// 	if isDir {
	// 		driveFName = dir
	// 		message = fmt.Sprintf("Fetching directory %s", dir)
	// 	} else {
	// 		driveFName = file
	// 		message = fmt.Sprintf("Fetching file %s", file)
	// 	}

	// 	_, err := utils.Spinner(func() (struct{}, error) {
	// 		return struct{}{}, fetch(driveFName)
	// 	}, message)

	// 	if err != nil {
	// 		utils.ExitOnError("Unable to complete the fetch operation: " + err.Error())
	// 	}

	// 	utils.ExitOnSuccess("Fetch operation for " + driveFName + " successfully completed!")
	// },
}

func init() {
	// fetchCmd.Flags().StringVarP(&file, "file", "f", "", "file to fetch")
	// fetchCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to fetch")
	// fetchCmd.Flags().StringVarP(&outDir, "outPath", "o", ".", "path for placing the fetched file/directory")
	rootCmd.AddCommand(fetchCmd)
}

// func fetch(driveFileName string) error {
// 	var err error

// 	if isDir {
// 		err = fetchDriveFolder(driveFileName)
// 	} else {
// 		err = fetchDriveFile(driveFileName)
// 	}

// 	return err
// }

// func download(driveFileId string, driveFileName string) error {
// 	res, err := srv.Files.Get(driveFileId).Download()
// 	if err != nil {
// 		return err
// 	}

// 	defer res.Body.Close()

// 	out := filepath.Join(outDir, driveFileName)
// 	outFile, err := os.Create(out)
// 	if err != nil {
// 		return err
// 	}
// 	defer outFile.Close()

// 	_, err = io.Copy(outFile, res.Body)
// 	return err
// }

// func fetchDriveFile(f string) error {
// 	q := fmt.Sprintf("name='%s' and trashed=false", f)
// 	nextPageToken := ""
// 	var err error

// 	for {
// 		req := srv.Files.List().Q(q).PageSize(5).Fields("nextPageToken, files(id, name, mimeType)")
// 		var res *drive.FileList

// 		if nextPageToken != "" {
// 			req = req.PageToken(nextPageToken)
// 		}

// 		res, err = req.Do()
// 		if err != nil {
// 			break
// 		}

// 		if len(res.Files) == 0 {
// 			err = fmt.Errorf("no artifacts found with the name `%s` in the drive", f)
// 		} else {
// 			if len(res.Files) > 1 {
// 				fmt.Println()
// 				utils.Warn("Found more than one matching artifacts in the drive, attempting to fetch the first one!")
// 				for _, df := range res.Files {
// 					fmt.Printf("Name: %s Id: [%s] MimeType: %s\n", df.Name, df.Id, df.MimeType)
// 				}
// 				err = download(res.Files[0].Id, res.Files[0].Name)
// 			} else {
// 				err = download(res.Files[0].Id, res.Files[0].Name)
// 			}
// 		}

// 		if err != nil {
// 			break
// 		}

// 		if res.NextPageToken == "" {
// 			break
// 		}
// 		nextPageToken = res.NextPageToken
// 	}

// 	return err
// }

// func fetchDriveFolder(fol string) error {
// 	panic("unimplemented")
// }
