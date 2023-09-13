package bulkgen

import (
	"bufio"
	"fmt"
	"os"
	"simify/fileutil"
	"strconv"
	"time"
)

func GenNumbers() {
	reader := bufio.NewReader(os.Stdin)
	var mainCodeStr string
	fmt.Print("\n\nEnter Country Code Without '+' :> ")
	fmt.Scanln(&mainCodeStr)
	mainCode, err := strconv.Atoi(mainCodeStr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Print("\nEnter your area code(s) (315,406) :> ")
	areaCode, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if areaCode == "\n" {
		fmt.Printf("invalid input. Exiting Program...\n")
	}

	parsedCodes, err := ParseAreaCode(areaCode)
	if err != nil {
		fmt.Printf("invalid input. Exiting Program...\n")
		return
	}

	var totalGenStr string
	var totalGen int
	if len(parsedCodes) > 1 {
		fmt.Print("\nHow many do you want to generate for each area code? :> ")
		fmt.Scanln(&totalGenStr)
		totalGen, err = strconv.Atoi(totalGenStr)
		if err != nil {
			fmt.Printf("invalid input. Exiting Program...\n")
			return
		}
	} else {
		fmt.Print("\nHow many do you want to generate? :> ")
		fmt.Scanln(&totalGenStr)
		totalGen, err = strconv.Atoi(totalGenStr)
		if err != nil {
			fmt.Printf("invalid input. Exiting Program...\n")
			return
		}
	}

	currentTime := time.Now().Unix()
	fileName := fmt.Sprintf("results_%d.txt", currentTime)
	file, err := fileutil.WriteToFile("Generated_Numbers", fileName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	filteredNums, elapsedTime := ProcessGen(mainCode, parsedCodes, totalGen, file)
	fmt.Println()
	fmt.Printf("\U0001F604 %d numbers generated after removing duplicates. The process took %0.2f\n", filteredNums, float32(elapsedTime.Seconds()))

}
