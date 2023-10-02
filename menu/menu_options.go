package menu

import (
	"simify/bulkgen"
	"simify/lookup"
)

func MenuSelection(choice int) {
	switch choice {
	case 1:
		bulkgen.GenNumbers()
	case 2:
		lookup.HRLLOOKUP()
	}
}
