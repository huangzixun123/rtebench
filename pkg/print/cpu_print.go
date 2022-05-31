package print

import "fmt"

func NewCpuPrint(reports []string) {
	fmt.Println("Test Record")
	for k, v := range reports {
		fmt.Printf(" Record %d\n", k)
		fmt.Printf("%v\n", v)
	}
}
