package main

import (
    "./media"
    "flag"
    "fmt"
    "github.com/kr/fs"
    "runtime"
)

func main() {
    var iPhotoPath string
    flag.StringVar(&iPhotoPath, "path", "path", "full path to iphoto library")
    flag.Parse()

    rootdir := iPhotoPath

    nCPU := runtime.NumCPU()
    // Removed this because it broke using it :-(
    //runtime.GOMAXPROCS(nCPU)

    var medialist []*media.Media

    fswalker := fs.Walk(rootdir)
    for fswalker.Step() {
        if err := fswalker.Err(); err != nil {
            fmt.Println("ERROR: ", err)
            continue
        }

        if stat := fswalker.Stat(); !stat.IsDir() {
            path := fswalker.Path()

            m, _ := media.NewMedia(path)
            if m.IsValid() {
                medialist = append(medialist, m)
            }
        }
    }

    completed := make(chan *media.Media, len(medialist)/nCPU)
    for _, m := range medialist {
        go func(m media.Media) {
            m.CalcMd5()
            completed <- &m
        }(*m)
    }

    duplicates := make(map[string][]string)

    for i := 0; i < len(medialist); i++ {
        res := <-completed

        md5string := fmt.Sprintf("%x", res.Md5Sum)

        duplicates[md5string] = append(duplicates[md5string], res.FullPath)
    }

    for md5sum, pathlist := range duplicates {
        if len(pathlist) >= 2 {
            //delete(duplicates, md5sum)
            fmt.Printf("Duplicates found, md5: %s\n", md5sum)
            for _, path := range pathlist {
                fmt.Printf("\tPath: %s\n", path)
            }
        }
    }
}
