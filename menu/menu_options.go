package menu

import "simify/bulkgen"

func MenuSelection(choice int) {
	switch choice {
	case 1:
		bulkgen.GenNumbers()
	}
}
