package main

import (
	"namaz-time-bot/bot"
	"namaz-time-bot/config"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	config.LoadEnv() // Загружаем переменные окружения
	bot.StartBot()   // Запускаем бота
}
