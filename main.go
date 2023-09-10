package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/manif3station/shared_lib"
)

type schema struct {
	Projects []details `json:"projects"`
}

type details struct {
	Name     string `json:"name"`
	BaseDir  string `json:"base_dir"`
	BuildDir string `json:"build_dir"`
}

func main() {
	wd, _ := os.Getwd()
	config := get_config(wd)
	for _, project := range config.Projects {
		build(project)
		os.Chdir(wd)
	}
}

func build(config details) {
	prog := config.Name
	base := template(config.BaseDir)
	dir := template(config.BuildDir)

	if prog == "" {
		log.Fatal("Missing Program Name.")
	}

	if shared_lib.Dir_exists(base) {
		os.Chdir(base)
	} else {
		log.Fatal("Base directory '" + base + "' is not a valid directory.")
	}

	if !shared_lib.Dir_exists(dir) {
		log.Fatal("Build directory '" + dir + "' is not a valid directory.")
	}

	ext := ""

	if runtime.GOOS == "windows" {
		ext = "exe"
	} else if runtime.GOOS == "darwin" {

		ext = "mac"
	} else if runtime.GOOS == "linux" {

		ext = "linux"
	}

	bin := fmt.Sprintf("%s/%s.%s", dir, prog, ext)

	cmd := exec.Command("go", "build", "-o", bin, ".")

	fmt.Println(">> cd", base)
	fmt.Println(">>", cmd)

	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	if len(out) > 0 {
		fmt.Println(string(out))
	}

	fmt.Println("--------------------")
}

func get_config(wd string) schema {
	var config schema

	file := wd + "/Buildfile.json"

	if info, err := os.Stat(file); err != nil || info == nil || !info.Mode().IsRegular() {
		return config
	}

	json_str, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(json_str, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func template(path string) string {
	path = shared_lib.Replace(`{{home}}`, path, shared_lib.MyHomeFolder())
	path = shared_lib.Replace(`{{bin}}`, path, shared_lib.MyHomeItem("bin"))
	path = shared_lib.Replace(`{{web}}`, path, shared_lib.MyHomeItem("web"))
	path = shared_lib.Replace(`{{lab}}`, path, shared_lib.MyHomeItem("web")+"/lab")
	return path
}
