package main

import (
        "fmt"
//        "flag"
        "io/ioutil"
        flags "github.com/jessevdk/go-flags"
        "os/exec"
        "os"
        "path/filepath"
)


type Options struct {
        Name string `short:"n" long:"name" description:"A name" required:"true"`
}

var opts Options

func delete_strings(slice []string, s string) (string, []string, error) {
        ret := make([]string, len(slice))
        i := 0
        for _, x := range slice {
                if s != x {
                        ret[i] = x
                        i++
                }
        }
        if len(ret[:i]) == len(slice) {
                return "", slice, fmt.Errorf("Couldn't find")
        }
        return s, ret[:i], nil
}

func dirwalk(dir string) []string {
        files, err := ioutil.ReadDir(dir)
        if err != nil {
                panic(err)
        }
        
        var paths []string
        for _, file := range files {
                if file.IsDir() {
                        paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
                        continue
                }
                paths = append(paths, filepath.Join(dir, file.Name()))
        }
        
        return paths
}

func main(){

        _, err := flags.Parse(&opts)
        
        if err != nil {
                os.Exit(1)
        }

        fmt.Printf("Name: %s\n", opts.Name)


        files:=dirwalk(".")
        
        files = append(files, opts.Name + ".epub")
        
        // html変換
        out, err := exec.Command("make4ht", "-l" , opts.Name + ".tex").Output()
        
        if err != nil {
                fmt.Println("make4ht Exec Error.")
                os.Exit(1)
        }
        
        // 実行したコマンドの結果を出力
        fmt.Printf("make4ht: \n%s", string(out))

        // html→epub3に変換
        out2, err2 := exec.Command("pandoc" ,"-t" ,"epub3" ,"-o" , opts.Name + ".epub" ,"--css=" + opts.Name + ".css" ,  opts.Name + ".html" ).Output()
        
        if err2 != nil {
                fmt.Println("pandoc Exec Error.")
        os.Exit(1)
        }
                
        // 実行したコマンドの結果を出力
        fmt.Printf("pandoc: \n%s", string(out2))
        
        files2:=dirwalk(".")
        
        for _, value := range files {
                _, files2, err = delete_strings(files2, value)
        }

        fmt.Println(files2)
                
        for _, value := range files2 {
                if err := os.Remove( value ); err != nil {
                        fmt.Println(err)
                }
        }
        
}