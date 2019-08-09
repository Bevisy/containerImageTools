package main

import (
    "archive/tar"
    "io"
    "log"
    "os"
)

type file struct {
    name    string
    ext     string
    absPath string
}

func main() {
    //获取程序执行路径
    //dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    //dir := filepath.Dir(os.Args[0])
    //if err != nil {
    //  log.Printf("%v", err)
    //}

    // 定义源文件属性
    srcFile := file{
        name:    "sample",
        ext:     ".txt",
        absPath: "E:\\go\\src\\github.com\\bevisy\\imageTool\\pkg\\src\\",
    }

    // 实例化目标文件
    var desFile file
    desFile.name = srcFile.name
    desFile.ext = ".tar"
    desFile.absPath = srcFile.absPath
    // 目标文件存在则删除
    if err := os.Remove(desFile.absPath + desFile.name + desFile.ext); err != nil {
        log.Println(err)
    }
    // 创建目标文件
    fw, err := os.Create(desFile.absPath + desFile.name + desFile.ext)
    if err != nil {
        log.Println(err)
    }
    defer fw.Close()

    // 创建 tar.writer
    tw := tar.NewWriter(fw)
    defer func() {
        if err := tw.Close(); err != nil {
            log.Println(err)
        }
    }()

    // tar包包含两部分信息：文件信息 和 文件数据
    // 通过 Stat 读取 Fileinfo，然后通过 FileInfoHeader 得到 hdr tar.*Header
    fi, err := os.Stat(srcFile.absPath + srcFile.name + srcFile.ext)
    if err != nil {
        log.Println(err)
    }
    hdr, err := tar.FileInfoHeader(fi, "")
    // 将 tar 的文件信息 har 写入到 tw
    err = tw.WriteHeader(hdr)
    if err != nil {
        log.Println(err)
    }

    // 将文件数据写入
    // 打开准备写入的文件
    fr, err := os.Open(srcFile.absPath + srcFile.name + srcFile.ext)
    if err != nil {
        log.Println(err)
    }
    defer fr.Close()

    written, err := io.Copy(tw, fr)
    if err != nil {
        log.Println(err)
    }
    log.Printf("写入字符个数: %d\n", written)

}
