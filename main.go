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
				
					░██████╗██╗███╗░░░███╗██╗███████╗██╗░░░██╗
					██╔════╝██║████╗░████║██║██╔════╝╚██╗░██╔╝
					╚█████╗░██║██╔████╔██║██║█████╗░░░╚████╔╝░
					░╚═══██╗██║██║╚██╔╝██║██║██╔══╝░░░░╚██╔╝░░
					██████╔╝██║██║░╚═╝░██║██║██║░░░░░░░░██║░░░
					╚═════╝░╚═╝╚═╝░░░░░╚═╝╚═╝╚═╝░░░░░░░░╚═╝░░░


					     [x] Created By @yeptg [x]
				 
				 
	`
	menu.SlowPrintArt(art, time.Millisecond*50)
	menuOptions := `
		[1] Bulk Number Generator				[2] HRL Lookup Using Telnyx API Key


		[3] GMAIL SMTPs To SMS (USA ONLY)			[4] OTHER SMTPs To SMS (USA ONLY)

		
	`
	menu.SlowPrintMenu(menuOptions, time.Millisecond*50)
	var choiceStr string
	color.New(color.FgHiRed).Print("(")
	color.New(color.FgHiGreen).Print("~")
	color.New(color.FgHiRed).Print(")")
	fmt.Print(" Enter your option :> ")
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		color.New(color.FgRed).Println("invalid choice. Exiting Program...")
		return
	}
	menu.MenuSelection(choice)
}
