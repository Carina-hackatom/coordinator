package logic

import "log"

func botTickLog(botType string, elapsedSecs int) {
	log.Printf("----> %s Bot is ongoing for %d secs <----\n\n", botType, elapsedSecs)
}
