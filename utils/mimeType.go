package utils

var mimeList = map[string]string{
	// Images
	"bmp":  "image/bmp",
	"gif":  "image/gif",
	"jpe":  "image/jpg",
	"jpeg": "image/jpg",
	"jpg":  "image/jpg",
	"svg":  "image/svg+xml",
	"png":  "image/png",
	"ico":  "image/x-icon",

	// Text
	"txt":  "text/plain",
	"htm":  "text/html",
	"html": "text/html",
	"css":  "text/css",

	// Application
	"json": "application/json",
	"pdf":  "application/pdf",
	"tgz":  "application/x-compressed",
	"gz":   "application/x-gzip",
	"js":   "application/x-javascript",
	"zip":  "application/zip",
}

func GetMime(ext string) string {
	mime, exists := mimeList[ext]

	if !exists {
		return ""
	}
	return mime
}
