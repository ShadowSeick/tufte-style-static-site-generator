package domain

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	Name string
	Language Language
	Date *time.Time
	Type FileType
	Content string
}

func (f *File) DateString() string {
	if f.Date == nil {
		return ""
	}
	return f.Date.Format(time.DateOnly)
}

func (f *File) FullName() string {
	if date := f.DateString(); date != "" {
		return fmt.Sprintf("%s_%s", f.DateString(), f.Name)
	}
	return f.Name
}

func (f *File) LocalFilePath() string {
	path, err := filepath.Abs(filepath.Join(f.Type.LocalDirectory(), f.Language.String(), f.FullName()))
	if err != nil {
		panic("invalid filepath")
	}
	return fmt.Sprintf("%s.%s", path, Markdown.String())
}

func (f *File) RemoteFilePath() string {
	path, err := filepath.Abs(filepath.Join(f.Type.RemoteDirectory(), f.Language.String(), f.FullName()))
	if err != nil {
		panic("invalid filepath")
	}
	return fmt.Sprintf("%s.%s", path, Markdown.String())
}

func (f *File) DebugFilePath() string {
	path, err := filepath.Abs(filepath.Join(f.Type.DebugDirectory(), f.Language.String(), f.FullName()))
	if err != nil {
		panic("invalid filepath")
	}
	return fmt.Sprintf("%s.%s", path, Markdown.String())
}

func (f *File) Title() string {
	switch f.Type {
	case Article:
		name := strings.Join(strings.Split(f.Name, "-"), " ")
		if len(name) == 0 {
			panic("article name should be at least 1 character long")
		}
		return strings.ToUpper(string(name[0])) + name[1:]
	case Index:
		return "Nerd with a mouth - Blog"
	case ArticleIndex:
		return "Nerd with a mouth - Articles"
	case ContactIndex:
		return "Nerd with a mouth - Contact"
	default:
		panic("title not handled for file type "+f.Type.String())
	}
	}

func (f *File) Checksum() string {
	if f.Content == "" {
		panic("file does not have any content")
	}
	return strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(f.Content))))
}

func NewFile(language Language, fileBaseName string, fileType FileType) (File, error) {
	var file File

	cleanedFileName := fileBaseName
	switch fileType {
	case Article:
		fullName := strings.Split(fileBaseName, "_")
		if len(fullName) != 2 {
			return file, fmt.Errorf("wrong article name! The name should be [2026-01-01]_[article-name].md; name: %s", fileBaseName)
		}

		date, err := time.Parse(time.DateOnly, fullName[0])
		if err != nil {
			return file, fmt.Errorf("invalid date time: %w", err)
		}
		file.Date = &date

		cleanedFileName = fullName[1]
	}

	baseName := strings.Split(cleanedFileName, ".")
	if len(baseName) != 2 {
		return file, fmt.Errorf("invalid %s name: %s", fileType.String(), fileBaseName)
	}

	file.Name = baseName[0]
	file.Language = language
	file.Type = fileType

	return file, nil
}

type FileType uint8

const (
	Article FileType = iota
	Index
	ArticleIndex
	ContactIndex
	Image
	Gif
	FileTypeCount
)

var fileTypeStrings = [FileTypeCount]string{
	Article: "article",
	Index: "home page",
	ArticleIndex: "articles index page",
	ContactIndex: "contact info",
	Image: "image",
	Gif: "gif",
}

func (f FileType) String() string {
	if f >= FileTypeCount {
		panic("invalid file type")
	}
	return fileTypeStrings[f]
}

var remoteDirectory = [FileTypeCount]string{
	Article: "articles",
	Index: "",
	ArticleIndex: "articles",
	ContactIndex: "",
	Image: "assets/images",
	Gif: "assets/images",
}

var localDirectory = [FileTypeCount]string{
	Article: "articles",
	Index: "",
	ArticleIndex: "articles",
	ContactIndex: "",
	Image: "assets/images",
	Gif: "assets/images",
}

var debugDirectory = [FileTypeCount]string{
	Article: "debug/articles",
	Index: "debug",
	ArticleIndex: "debug/articles",
	ContactIndex: "debug",
	Image: "assets/images",
	Gif: "assets/images",
}

func (f FileType) RemoteDirectory() string {
	if f >= FileTypeCount {
		panic("invalid file type")
	}
	return remoteDirectory[f]
}

func (f FileType) LocalDirectory() string {
	if f >= FileTypeCount {
		panic("invalid file type")
	}
	return localDirectory[f]
}

func (f FileType) DebugDirectory() string {
if f >= FileTypeCount {
		panic("invalid file type")
	}
	return debugDirectory[f]
}
