package main

import (
    "archive/tar"
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
)

func Tar(src, dst string) (err error) {
    fw, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer fw.Close()

    gw := gzip.NewWriter(fw)
    defer gw.Close()

    tw := tar.NewWriter(gw)
    defer tw.Close()

    return filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        hdr, err := tar.FileInfoHeader(fi, "")
        if err != nil {
            return err
        }

        hdr.Name = strings.TrimPrefix(fileName, string(filepath.Separator))

        if err := tw.WriteHeader(hdr); err != nil {
            return err
        }

        if !fi.Mode().IsRegular() {
            return nil
        }

        fr, err := os.Open(fileName)
        if err != nil {
            return err
        }
        defer fr.Close()

        n, err := io.Copy(tw, fr)
        if err != nil {
            return nil
        }

        log.Printf("拷贝字节数为：%d", n)

        return nil
    })
}

func main() {
    src := "src/nginx-1.15"
    dst := fmt.Sprintf("%s.tar.gz", src)

    if err := Tar(src, dst); err != nil {
        log.Fatalln(err)
    }
}
