package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

var availableVersion string = ""
var targetVersion string = ""

func getVersionFromString(target string, start int) string {
	var version string = ""
	for i, c := start, 0; i < len(target); i++ {
		if target[i] == '.' {
			c++
		}
		if c < 2 {
			version += string(target[i])
		}
	}
	return version
}

func execute(cmd string, arg ...string) {

	out, err := exec.Command(cmd, arg[0]).Output()

	if err != nil {
		panic(err)
	}
	output := string(out[:])

	availableVersion = getVersionFromString(output, 4)
}

func getActualVersion() {
	resp, err := http.Get("https://api.github.com/repos/zed-industries/zed/releases/latest")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// i := strings.Index(string(body), "tag_name")
	rm := make(map[string]json.RawMessage)

	err = json.Unmarshal(body, &rm)

	if err != nil {
		panic(err)
	}

	targetVersion = getVersionFromString(string(rm["tag_name"]), 2)

}

func main() {
	execute("zed", "--version")
	getActualVersion()

	if availableVersion != targetVersion {
		execute("curl https://zed.dev/install.sh | sh")
		fmt.Println("\n\n\nSuccessfully update")
	} else {
		fmt.Println("The versions match")
	}
}
