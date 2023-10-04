package menu

import (
	"fmt"
	"simify/bulkgen"
	bulkapi "simify/lookup/bulk_api"
	"simify/lookup/single_api"
	"simify/sms/gmailsmtp"
	"simify/sms/othersmtp"
)

func MenuSelection(choice int) {
	switch choice {
	case 1:
		bulkgen.GenNumbers()
	case 2:
		single_api.SingleAPILookup()
	case 3:
		bulkapi.BulkAPILookup()
	case 4:
		gmailsmtp.GmailSMTPToSMS()
	case 5:
		othersmtp.OtherSMTPToSMS()
	default:
		fmt.Println("Exiting Program...")
	}
}
