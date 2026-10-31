package controllers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const ImageDir = "static/img"

const MaxImageSize = 10 << 20 // 10 MB

// Allowed image extensions and the content type their contents must sniff as.
// SVG is left out on purpose, it can contain scripts that would run on this origin.
var imageTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".ico":  "image/x-icon",
}

var (
	ErrImageType   = errors.New("only jpg, png, gif, webp and ico images are allowed")
	ErrImageExists = errors.New("an image with this name already exists")
	ErrImageName   = errors.New("invalid file name")

	unsafeNameChars = regexp.MustCompile(`[^A-Za-z0-9._-]+`)
)

// ListImages returns the web paths of all images in ImageDir, e.g. /static/img/me.jpg.
func ListImages() ([]string, error) {
	entries, err := os.ReadDir(ImageDir)
	if err != nil {
		return nil, err
	}
	var images []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, ok := imageTypes[strings.ToLower(filepath.Ext(entry.Name()))]; ok {
			images = append(images, "/"+ImageDir+"/"+entry.Name())
		}
	}
	sort.Strings(images)
	return images, nil
}

// SaveImage stores an uploaded image in ImageDir under a sanitized version of name and returns its web path.
// It never overwrites an existing file.
func SaveImage(name string, file io.Reader) (string, error) {
	name = unsafeNameChars.ReplaceAllString(filepath.Base(strings.TrimSpace(name)), "-")
	name = strings.Trim(name, ".-")
	ext := strings.ToLower(filepath.Ext(name))
	if name == "" || name == ext {
		return "", ErrImageName
	}
	contentType, ok := imageTypes[ext]
	if !ok {
		return "", ErrImageType
	}

	head := make([]byte, 512)
	n, err := io.ReadFull(file, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}
	head = head[:n]
	if http.DetectContentType(head) != contentType {
		return "", ErrImageType
	}

	if err := os.MkdirAll(ImageDir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(ImageDir, name)
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", ErrImageExists
		}
		return "", err
	}

	_, err = out.Write(head)
	if err == nil {
		_, err = io.Copy(out, file)
	}
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(path)
		return "", err
	}
	return "/" + ImageDir + "/" + name, nil
}
