package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

// Complete the caesarCipher function below.
func caesarCipher(s string, k int32) string {
    res := ""
    for _, ch := range s {
        if ch >= 'a' && ch <= 'z' {
            newCh := (((ch - 'a') + k) % 26) + 'a'
            res += string(newCh)
        } else if ch >= 'A' && ch <= 'Z' {
            newCh := (((ch - 'A') + k) % 26) + 'A'
            res += string(newCh)
        } else {
            res += string(ch)
        }
    }

    return res
}

func main() {
    reader := bufio.NewReaderSize(os.Stdin, 1024 * 1024)

    stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
    checkError(err)

    defer stdout.Close()

    writer := bufio.NewWriterSize(stdout, 1024 * 1024)

    _, err = strconv.ParseInt(readLine(reader), 10, 64)
    checkError(err)

    s := readLine(reader)

    kTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
    checkError(err)
    k := int32(kTemp)

    result := caesarCipher(s, k)

    fmt.Fprintf(writer, "%s\n", result)

    writer.Flush()
}

func readLine(reader *bufio.Reader) string {
    str, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }

    return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}
