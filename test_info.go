package zardoz

import (
    "bufio"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readCode(fpath string, startLine int, endLine int) string {
    var code string = ""
    file, err := os.Open(fpath)
    defer file.Close()

    check(err)

    lineNum := 0
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        lineNum++
        if lineNum < startLine { continue }
        if lineNum > endLine { continue }

        code += scanner.Text() + "\n"
    }

    if err := scanner.Err(); err != nil {
        check(err)
    }

    return code
}
