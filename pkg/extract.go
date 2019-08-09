package main

import (
    "archive/tar"
    "io"
    "log"
    "os"
)

type file2 struct {
    name    string
    ext     string
    absPath string
}

func main() {
    // 定义源文件属性
    srcFile := file2{
        name:    "sample",
        ext:     ".tar",
        absPath: "E:\\go\\src\\github.com\\bevisy\\imageTool\\pkg\\src\\",
    }

    // 实例化目标文件
    var desFile file2
    desFile.name = srcFile.name
    desFile.ext = ".txt"
    desFile.absPath = srcFile.absPath
    // 目标文件存在则删除，不存在则报错，但程序继续执行
    if err := os.Remove(desFile.absPath + desFile.name + desFile.ext); err != nil {
        log.Println(err)
    }

    // 打开 tar 包
    fr, err := os.Open(srcFile.absPath + srcFile.name + srcFile.ext)
    if err != nil {
        log.Println(err)
    }
    defer fr.Close()
    // 创建 tar.*Reader 结构 tr；遍历 tr,将数据保存在磁盘中
    tr := tar.NewReader(fr)

    for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
        if err != nil {
            log.Println(err)
        }
        // 获取文件信息
        fi := hdr.FileInfo()
        // 创建空文件，用来写入解包后的数据
        fw, err := os.Create(desFile.absPath + fi.Name())
        if err != nil {
            log.Println(err)
        }

        w, err := io.Copy(fw, tr)
        if err != nil {
            log.Println(err)
        }
        log.Println("写入字符个数：%d", w)
        // 设置文件权限与原文件相同
        os.Chmod(desFile.absPath+fi.Name(), fi.Mode().Perm())
        // 由于处于循环中，所以未所以 defer 关闭文件
        // 如果需要使用 defer， 则可将文件写入步骤单独封装在一个函数中
        fw.Close()
    }
}
