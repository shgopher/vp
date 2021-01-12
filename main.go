package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	getFilePath()
}

var (
	// 输入路径
	inputPath  string
	// 输出路径
	outputPath string
	//码率
	code string
	// 分辨率
	mass string
	// 视频帧率
	fps string
)

func init() {
	flag.StringVar(&inputPath, "input", ".", "输入的路径")
	flag.StringVar(&outputPath, "output", "../", "输出的路径，【注意不要和输入的路径在同一个路径下】")
	flag.StringVar(&code,"code","40k","码率 例如 400k")
	flag.StringVar(&mass,"mass","960x540","分辨率，例如 960x540")
	flag.StringVar(&fps,"fps","20","视频帧率，例如 25")
	flag.Parse()
}

func getFilePath() {
	pathChan := make(chan string, 0)
	go func() {
		defer close(pathChan)
		filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			pathChan <- path
			return nil
		})
	}()
	wg := new(sync.WaitGroup)
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()
		L:
			for {
				select {
				case p, ok := <-pathChan:
					if ok {
						_, output := filepath.Split(p)
						output = filepath.Join(outputPath,output)
						fmt.Println("正在执行：",output)
						deal(p, output)
					} else {
						fmt.Println("此goroutine 全部任务执行完毕，请示退出。")
						break L
					}
				default:
					time.Sleep(time.Second >> 2)
				}
			}
		}()
	}
	wg.Wait()
}
func deal(input, output string) {
	cmd := exec.Command("ffmpeg", "-i", input, "-b:v", code, "-s", mass,"-r",fps, output)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("ffmpeg输出: " + out.String())
}
