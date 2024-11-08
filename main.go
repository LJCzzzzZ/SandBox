package main

import "fmt"

func main() {

	// 采用seccomp对系统调用进行限制
	s := fmt.Errorf("ljc: %d", 12)
	fmt.Println(s)
}
