package handlers

//
//import (
//	"fmt"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/robfig/cron/v3"
//)
//
////func StartScheduler() {
////	c := cron.New()
////	c.AddFunc("0 5 * * *", func() { sendReminder("Фаджр") })
////	c.AddFunc("0 12 * * *", func() { sendReminder("Зухр") })
////	c.AddFunc("0 15 * * *", func() { sendReminder("Аср") })
////	c.AddFunc("0 18 * * *", func() { sendReminder("Магриб") })
////	c.AddFunc("0 20 * * *", func() { sendReminder("Иша") })
////	c.Start()
////}
////
////func sendReminder(prayer string) {
////	msg := fmt.Sprintf("Напоминание: Время намаза – %s", prayer)
////	botAPI.Send(tgbotapi.NewMessage(-1001234567890, msg))
////}
////
////func GetPrayerTimes() string {
////	return "Время намазов: Фаджр - 5:00, Зухр - 12:00, Аср - 15:00, Магриб - 18:00, Иша - 20:00"
////}
