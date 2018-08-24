package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	UPLOAD_DIR = "./uploads"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.html") //加载html
		log.Println(t.Execute(w, nil))
		return
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("image") //此值与html的取值有关
		fmt.Println(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		fmt.Println("filename is ", filename)
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		/*
		 */

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("hello")
		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}
	/* if r.Method==delete(map[typeA]typeB, typeA){


	   }
	*/
}
func viewHandle(w http.ResponseWriter, r *http.Request) {

	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)

}
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
func listHandler(w http.ResponseWriter, r *http.Request) {
	_path := strings.TrimRight(r.FormValue("id"), "/") + "/"
	fileInfoArr, err := ioutil.ReadDir("./uploads" + _path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	locals := make(map[string]interface{})
	images := []map[string]string{}
	files := []map[string]string{}
	files = append(files, map[string]string{"url": "/list?id=" + _path + "../", "name": "上一级"})
	for _, fileInfo := range fileInfoArr {
		_name := fileInfo.Name()
		_nameSlice := strings.Split(_name, ".")
		switch _nameSlice[len(_nameSlice)-1] {
		case "jpg":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "JPG":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "png":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "PNG":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "gif":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "GIF":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "bmp":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		case "BMP":
			images = append(images, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
		default:
			if fileInfo.IsDir() {
				files = append(files, map[string]string{"url": "/list?id=" + _path + _name, "name": _name})
			} else {
				files = append(files, map[string]string{"url": "/view?id=" + _path + _name, "name": _name})
			}
		}

	}
	locals["images"] = images
	locals["files"] = files
	t, err := template.ParseFiles("list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, locals)
}
func main() {

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view", viewHandle)
	http.HandleFunc("/", listHandler)
	err := http.ListenAndServe(":9010", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err.Error())
	}
}
