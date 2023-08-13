package helper

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanDirectory(dirPath string) []string {
	var videoFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		// Check if the file has a video file extension
		if !info.IsDir() && isVideoFile(path) {
			videoFiles = append(videoFiles, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return videoFiles
}

func isVideoFile(path string) bool {
	// List of video file extensions
	videoExtensions := []string{".3g2",
		".3gp",
		".amv",
		".asf",
		".avi",
		".drc",
		".flv",
		".f4v",
		".f4p",
		".f4a",
		".f4b",
		".gifv",
		".m4v",
		".mkv",
		".mng",
		".mov",
		".qt",
		".mp4",
		".mpe",
		".mpeg",
		".mpg",
		".mpv",
		".mxf",
		".nsv",
		".ogv",
		".ogg",
		".rm",
		".rmvb",
		".roq",
		".svi",
		".vob",
		".webm",
		".wmv",
		".yuv"}

	// Get the file extension
	ext := strings.ToLower(filepath.Ext(path))

	// Check if the file extension is in the list of video extensions
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}

	return false
}

func GetFileNameFromPath(filePath string) string {
	_, fileName := filepath.Split(filePath)
	return fileName
}
