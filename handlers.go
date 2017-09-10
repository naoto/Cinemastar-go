package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"

	"log"
	"regexp"
	"strings"
        "time"
)

type FileInfos []os.FileInfo
type ByName struct{ FileInfos }

var latest string
var last = time.Now()

func (fi ByName) Len() int {
	return len(fi.FileInfos)
}
func (fi ByName) Swap(i, j int) {
	fi.FileInfos[i], fi.FileInfos[j] = fi.FileInfos[j], fi.FileInfos[i]
}
func (fi ByName) Less(i, j int) bool {
	return fi.FileInfos[j].ModTime().Unix() < fi.FileInfos[i].ModTime().Unix()
}

func ignoreFile(f string) bool {
	ignore := []string{".html", ".htm", ".jpg", ".jpeg", ".png", ".yml", ".yaml", ".txt", ".rss"}
	pos := strings.LastIndex(f, ".")

	if pos != -1 {
		extract := f[pos:]

		for _, v := range ignore {
			if v == extract {
				return true
			}
		}
	}

	ignoreDir := []string{"thumbnail", "bin", "html", "latest", "module", "lost+found"}

	for _, v2 := range ignoreDir {
		if v2 == f {
			return true
		}
	}

	return false
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "index.html")
}

func MovieContent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	basePath := ps.ByName("filepath")
	arg := base + basePath
	log.Print(arg)
	http.ServeFile(w, r, arg)
}

func Latest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

        duration := time.Since(last).Seconds()
        if duration < 10 || duration > 1800 {
            last = time.Now()
	    cmd := "find " + base + " -path " + base + "thumbnail -prune -o -path " + base + "t -prune -o -type f -not -name '*.html' -print0 | xargs -0 ls --full-time | sort -k6,7 -r | head -n 100 | rev | cut -f1 -d ' ' | rev"
            log.Print(cmd)
	    out, err := exec.Command("sh", "-c", cmd).Output()
	    if err != nil {
		fmt.Errorf("Latest Dead: %s\n", err)
		os.Exit(1)
	    }
            latest = string(out)
        }
	var movies Movies
	rep := regexp.MustCompile(`\.[^\.]+$`)
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(latest, -1) {
		fullpath := string(v)
		if ignoreFile(fullpath) {
			continue
		}

		file, err := os.Open(fullpath)
		if err == nil {
			fileinfo, _ := file.Stat()
			fileName := fileinfo.Name()
			name := strings.Replace(rep.ReplaceAllString(fileName, ""), "_", " ", -1)
			modTime := fileinfo.ModTime()
			isDir := fileinfo.IsDir()
			filePath := strings.Replace(fullpath, base, "/", 1)

			thumbnail := "/static/t/" + filePath + ".jpg"
			movies = append(movies, Movie{Name: name, Directory: isDir, Due: modTime, Path: filePath, Thumbnail: thumbnail})
		}
	}
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		panic(err)
	}
}

func MovieIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	basePath := ps.ByName("filepath")

	arg := base + basePath
	log.Printf(arg)
	fileInfos, err := ioutil.ReadDir(arg)
	if err != nil {
		fmt.Errorf("Directory cannot read %s\n", err)
		os.Exit(1)
	}

	var movies Movies
	rep := regexp.MustCompile(`\.[^\.]+$`)
	sort.Reverse(ByName{fileInfos})
	for _, fileInfo := range fileInfos {
		fileName := (fileInfo).Name()
		if ignoreFile(fileName) {
			continue
		}

		name := strings.Replace(rep.ReplaceAllString(fileName, ""), "_", " ", -1)
		modTime := (fileInfo).ModTime()
		isDir := (fileInfo).IsDir()
		filePath := basePath + "/" + fileName

		if !isDir {
			thumbnail := "/static/t/" + filePath + ".jpg"
			movies = append(movies, Movie{Name: name, Directory: isDir, Due: modTime, Path: filePath, Thumbnail: thumbnail})
		}
	}

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		panic(err)
	}
}

func MovieCategoryIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	basePath := ps.ByName("filepath")

	arg := base + basePath
	log.Printf(arg)
	fileInfos, err := ioutil.ReadDir(arg)
	if err != nil {
		fmt.Errorf("Directory cannot read %s\n", err)
		os.Exit(1)
	}

	var movies Movies
	sort.Reverse(ByName{fileInfos})
	for _, fileInfo := range fileInfos {
		fileName := (fileInfo).Name()
		if ignoreFile(fileName) {
			continue
		}

		modTime := (fileInfo).ModTime()
		isDir := (fileInfo).IsDir()
		filePath := basePath + fileName
		if isDir {
			movies = append(movies, Movie{Name: fileName, Directory: isDir, Due: modTime, Path: filePath})
		}
	}

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		panic(err)
	}
}
