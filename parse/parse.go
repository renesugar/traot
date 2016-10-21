package parse

import "os"
import "path/filepath"
import "fmt"
import "strings"
import "io/ioutil"

var files []string

func visit(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, ".java") {
		files = append(files, path)
	}
	return nil
}

func readJava(path, orig string) {
	fmt.Println(path)
	f, _ := os.Open(path)
	data, _ := ioutil.ReadAll(f)
	str := string(data)
	lines := strings.Split(str, "\n")

	tokens := strings.Split(path, "/")
	jpackage := make([]string, 0)
	for _, toke := range tokens {
		jpackage = append([]string{toke}, jpackage...)
	}
	jpackage = jpackage[1 : len(jpackage)-1]
	gopackage := make([]string, 0)
	for _, name := range jpackage {
		if name == "java" || name == "main" || name == "src" {
			break
		}
		gopackage = append(gopackage, name)
	}
	fmt.Println(gopackage)
	gopackname := "root"
	fullgopackname := ""
	if len(gopackage) > 0 {
		gopackname = gopackage[0] //strings.Join(gopackage, "_")
		fullgopackname = strings.Join(gopackage[1:len(gopackage)-1], "_")
	}
	dirname := fmt.Sprintf("%s_go/%s", orig, gopackname)
	os.Mkdir(dirname, 0777)
	dirname = fmt.Sprintf("%s_go/%s/%s", orig, gopackname, fullgopackname)
	fmt.Println("d", dirname)
	os.Mkdir(dirname, 0777)
	endparts := strings.Split(strings.ToLower(tokens[len(tokens)-1]), ".")
	endpart := endparts[0]
	filename := fmt.Sprintf("%s_go/%s/%s/%s", orig, gopackname, fullgopackname, endpart+".go")
	//fmt.Println(filename)
	wfile, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "public ") {
			wfile.WriteString("//" + line + "\n")
		}
		if strings.HasPrefix(line, "private ") {
			wfile.WriteString("//" + line + "\n")
		}
		if strings.HasPrefix(line, "protected ") {
			wfile.WriteString("//" + line + "\n")
		}
	}
}

func Parse(path string) {
	filepath.Walk(path, visit)
	for _, file := range files {
		readJava(file, path)
	}
}
