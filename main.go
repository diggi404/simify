package main

import (
	"fmt"
	"simify/menu"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func main() {
	fmt.Print("\033[H\033[2J")
	art := `
	
					░██████╗██╗███╗░░░███╗██╗███████╗██╗██╗░░░██╗
					██╔════╝██║████╗░████║██║██╔════╝██║╚██╗░██╔╝
					╚█████╗░██║██╔████╔██║██║█████╗░░██║░╚████╔╝░
					░╚═══██╗██║██║╚██╔╝██║██║██╔══╝░░██║░░╚██╔╝░░
					██████╔╝██║██║░╚═╝░██║██║██║░░░░░██║░░░██║░░░
					╚═════╝░╚═╝╚═╝░░░░░╚═╝╚═╝╚═╝░░░░░╚═╝░░░╚═╝░░░


					     [x] Created By @hidden404_bot [x]
				 
				 
	`
	menu.SlowPrintArt(art, time.Millisecond*50)
	menuOptions := `
		[1] Bulk Number Generator				[2] HRL Lookup Using Single Telnyx API


		[3] HRL Lookup Using Bulk Telnyx APIs			[4] GMAIL SMTPs To SMS (USA ONLY)


		[5] OTHER SMTPs To SMS (USA ONLY)			[6] Exit

		
	`
	menu.SlowPrintMenu(menuOptions, time.Millisecond*50)
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
