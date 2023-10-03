package main

import (
	"fmt"
	"simify/menu"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func main() {
	art := `
				_____ ________  ________________   __
				/  ___|_   _|  \/  |_   _|  ___\ \ / /
				\ '--.  | | | .  . | | | | |_   \ V / 
				 '--. \ | | | |\/| | | | |  _|   \ /  
				/\__/ /_| |_| |  | |_| |_| |     | |  
				\____/ \___/\_|  |_/\___/\_|     \_/  V1.0


				 [x] Created By @hidden404_bot [x]
				 
				 
	`
	menu.SlowPrintArt(art, time.Millisecond*30)
	menuOptions := `
		1. Bulk Number Generator				2. HRL Lookup Using Single Telnyx API


		3. HRL Lookup Using Bulk Telnyx APIs			4. GMAIL SMTPs To SMS


		5. OTHER SMTPs To SMS					5. Exit

		
	`
	menu.SlowPrintMenu(menuOptions, time.Millisecond*30)
	var choiceStr string
	fmt.Print("Enter your option :> ")
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		color.New(color.FgRed).Println("invalid choice. Exiting Program...")
		return
	}
	menu.MenuSelection(choice)
}
