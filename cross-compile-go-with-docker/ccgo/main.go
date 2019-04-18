// Copyright 2019 Evan Han
// 该工具用于借助docker交叉编译GOPATH中的包，需要安装docker和设置GOPATH环境变量。

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// 参数错误
var errArgs = errors.New("args error: run 'ccgo -h' for details")

// 工具运行参数
var (
	goos        = flag.String("os", "linux", "GOOS for compiler.")
	goarch      = flag.String("arch", "amd64", "GOARCH for compiler.")
	pkg         = flag.String("pkg", "github.com/go-learning/dev-kit/cross-compile-go-with-docker/ccgo", "The pkg that to be built")
	dockerImage = flag.String("image", "golang:1.11", "The image of docker that be used to compile go code.")
)

func main() {
	flag.Parse()
	// 需要设置GOPATH环境变量
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("You have to set GOPATH enviroment variable for ccgo.")
	}

	// 检查参数
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	// 运行docker命令
	if err := runCmd(gopath); err != nil {
		log.Fatal("run cmd error:", err)
	}
	log.Println("compile successfully!")
}

func checkArgs() error {
	if *goos == "" || *goarch == "" || *pkg == "" || *dockerImage == "" {
		return errArgs
	}
	return nil
}

func runCmd(gopath string) error {
	// 为了输入方便(cd pkgDir && ccgo -pkg=$PWD), 这里取GOPATH中包的相对路径
	var srcDir string
	if gopath[len(gopath)-1] == '/' {
		srcDir = gopath + "src/"
	} else {
		srcDir = gopath + "/src/"
	}
	pkgPath := strings.Replace(*pkg, srcDir, "", 1)

	// docker命令
	cmdStr := fmt.Sprintf("docker run --rm -v %s:/go -w /go/src/%s -e GOOS=\"%s\" -e GOARCH=\"%s\" %s go build -v",
		gopath, pkgPath, *goos, *goarch, *dockerImage)

	log.Println(cmdStr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	return cmd.Run()
}
