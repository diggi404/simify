package menu

import (
	"simify/bulkgen"
	"simify/lookup/single_api"
)

func MenuSelection(choice int) {
	switch choice {
	case 1:
		bulkgen.GenNumbers()
	case 2:
		single_api.SingleAPILookup()
	}
}
