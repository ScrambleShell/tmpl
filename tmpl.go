package main

import (
    "fmt"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "bytes"
    "io/ioutil"
    "regexp"
    "errors"
    "strings"
    "strconv"
    "text/template"
    "launchpad.net/goyaml"
)

func loadData(dataFile string) (map[interface{}]interface{}, error) {
    data, err := ioutil.ReadFile(dataFile)
    if err != nil {return nil, err}
    m := make(map[interface{}]interface{})
    err = goyaml.Unmarshal(data, m)
    if err != nil {return nil, err}
    return m, nil
}

func findAvailableTemplates() map[string]string {
    var (
        f *os.File
        list []string
    )
    m := make(map[string]string)
    me, err := user.Current()
    if err != nil { return m }
    tdir := filepath.Join(me.HomeDir, "templates")
    filepath.Walk(tdir, func(path string, info os.FileInfo, err error) error {
        if !strings.HasSuffix(path, ".tmpl") { return nil }
        name := filepath.Base(path)
        name = name[:len(name)-5]
        m[name] = path
        return nil
    })
    f, err = os.Open(".")
    if err != nil { return m }
    defer f.Close()
    list, err = f.Readdirnames(-1)
    if err != nil { return m }
    for _, name := range list {
        if !strings.HasSuffix(name, ".tmpl") { continue }
        path, _ := filepath.Abs(name)
        name = name[:len(name)-5]
        m[name] = path
    }
    return m
}

func loadTemplate(tmplName string) (*template.Template, error) {
    if strings.HasSuffix(tmplName, ".tmpl") {
        tmpl, err := template.ParseFiles(tmplName)
        if err == nil {
            return tmpl, nil
        }
        fmt.Println(err)
        tmplName = tmplName[:len(tmplName)-5]
    }
    tmplMap := findAvailableTemplates()
    path, ok := tmplMap[tmplName]
    if !ok {
        return nil, errors.New("Can't find template " + tmplName)
    }
    return template.ParseFiles(path)
}

func applyTemplate(tmplName, dataFile string, args []string) (string, error) {
    var (
        data map[interface{}]interface{}
        tmpl *template.Template
        err error
    )
    data, err = loadData(dataFile)
    if err != nil { return "", err }
    tmpl, err = loadTemplate(tmplName)
    if err != nil { return "", err }
    buffer := new(bytes.Buffer)
    err = tmpl.Execute(buffer, data)
    if err != nil { return "", err }
    str := buffer.String()
    re := regexp.MustCompile(`\$([0-9])`)
    matches := re.FindAllStringSubmatch(str, -1)
    matchMap := make(map[string]string)
    for _, match := range matches {
        matchMap[match[0]] = match[1]
    }
    if len(matchMap) != len(args) {
        return "", errors.New(
            fmt.Sprintf("Template requires %d command line arguments",
                len(matches)))
    }
    for k, v := range matchMap {
        idx, _ := strconv.Atoi(v)
        str = strings.Replace(str, k, args[idx], -1)
    }
    return str, nil
}

func runTemplate(tmplName, dataFile string, args []string) error {
    result, err := applyTemplate(tmplName, dataFile, args)
    if err != nil { return err }
    for _, line := range strings.Split(result, "\n") {
        line = strings.TrimSpace(line)
        if len(line) < 1 { continue }
        fmt.Println(line)
        err = execute(line)
        if err != nil { return err }
    }
    return nil
}

func execute(cmdString string) error {
    parts := strings.Split(cmdString, " ")
    prog := parts[0]
    args := make([]string, 0, 5)
    current := ""
    for _, part := range parts[1:] {
        part = strings.TrimSpace(part)
        if len(part) < 1 { continue }
        prefix := strings.HasPrefix(part, `"`)
        suffix := strings.HasSuffix(part, `"`)
        if prefix && suffix {
            args = append(args, part[1:len(part)-1])
            current = ""
        } else if prefix {
            current = part[1:]
        } else if suffix {
            args = append(args, current + " " + part[:len(part)-1])
            current = ""
        } else {
            if len(current) > 0 {
                current += " " + part
            } else {
                args = append(args, part)
            }
        }
    }
    cmd := exec.Command(prog, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func argumentsOk() bool {
    if len(os.Args) == 2 {
        return os.Args[1] == "list"
    }
    return len(os.Args) >= 4
}

func main() {
    var (
        result string
        err error
    )
    if !argumentsOk() {
        fmt.Println("usage: tmpl apply [template] [data]")
        fmt.Println("       tmpl run   [template] [data]")
        fmt.Println("       tmpl list")
        return
    }
    switch (os.Args[1]) {
        case "apply":
            result, err = applyTemplate(os.Args[2], os.Args[3], os.Args[4:])
            fmt.Println(result)
        case "run":
            err = runTemplate(os.Args[2], os.Args[3], os.Args[4:])
        case "list":
            tmplMap := findAvailableTemplates()
            for k, _ := range tmplMap {
                fmt.Println(k)
            }
        default:
            err = errors.New("Unknown action: " + os.Args[1])
    }
    if err != nil { fmt.Println(err); return }
}
