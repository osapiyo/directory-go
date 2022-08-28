package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	fmt.Println("パスを入力してください")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var str = scanner.Text()
	path := strings.Trim(str, "\"")

	fmt.Println(path)

	data := [][]string{{"dirName", "size"}}

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, f := range dirs {
		if f.IsDir() {
			fmt.Printf("- %s\n", f.Name())
			var p = filepath.Join(path, f.Name())
			size, err := DirSize(p)
			if err != nil {
				fmt.Println(err.Error())
			}
			// fmt.Printf("%20s : %v\n", f.Name(), size)
			sl := []string{f.Name(), strconv.FormatInt(size, 10)}
			data = append(data, sl)
		}
	}

	hd, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}
	fn := "dir-" + filepath.Base(path) + "_" + time.Now().Format("2006-01-02-150405") + ".csv"
	fp := filepath.Join(hd, "Downloads", fn)

	if err != nil {
		log.Fatal(err.Error())
	}

	f, err := os.Create(fp)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	w := csv.NewWriter(transform.NewWriter(f, japanese.ShiftJIS.NewEncoder()))
	w.WriteAll(data)
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("完了しました")
	scanner2 := bufio.NewScanner(os.Stdin)
	scanner2.Scan()
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
