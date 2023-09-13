package bulkgen

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func ProcessGen(mainCode int, areaCodes []int, totalGen int, file *os.File) (int, time.Duration) {
	generatedNumbers := make(map[string]bool)
	startTime := time.Now()
	for _, code := range areaCodes {
		generator := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < totalGen; i++ {
			num := 1000000 + generator.Intn(9999999)
			numStr := strconv.Itoa(num)
			if len(numStr) > 7 {
				num = num / 10
			}

			genNum := fmt.Sprintf("+%d%d%d\n", mainCode, code, num)

			if _, exists := generatedNumbers[genNum]; !exists {
				file.WriteString(genNum)
				generatedNumbers[genNum] = true
			}
			color.New(color.FgHiBlue).Printf("\rGenerating... %d/%d", i, totalGen)
		}
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	return len(generatedNumbers), elapsedTime
}
