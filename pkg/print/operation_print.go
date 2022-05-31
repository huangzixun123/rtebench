package print

import (
	"fmt"

	"github.com/huangzixun123/rtebench/pkg/maths"
	"github.com/huangzixun123/rtebench/pkg/types"
)

func NewOperationPrint(report types.Report) {
	fmt.Println("Test Record:")

	for k := range report["Create"] {
		fmt.Printf(" Record %d\n", k)
		fmt.Printf(" Create: %v | Run: %v | CreateAndRun: %v | Destroy: %v\n",
			report["Create"][k], report["Run"][k], report["CreateAndRun"][k], report["Destroy"][k])
	}

	fmt.Println("\nGeneral Statistics:")
	min, max, avg := helper(report["Run"])
	run := fmt.Sprintf(" Run: Min:%v, Max:%v, Avg: %v\n", min, max, avg)

	min, max, avg = helper(report["Create"])
	create := fmt.Sprintf(" Create: Min:%v, Max:%v, Avg: %v\n", min, max, avg)

	min, max, avg = helper(report["CreateAndRun"])
	createAndRun := fmt.Sprintf(" CreateAndRun: Min:%v, Max:%v, Avg: %v\n", min, max, avg)

	min, max, avg = helper(report["Destroy"])
	destroy := fmt.Sprintf(" Destroy: Min:%v, Max:%v, Avg: %v\n", min, max, avg)
	fmt.Println(create + run + createAndRun + destroy)
}

func helper(sli []float64) (min, max, avg float64) {
	min, _ = maths.Min(sli)
	max, _ = maths.Max(sli)
	avg, _ = maths.Avg(sli)
	return min, max, avg
}
