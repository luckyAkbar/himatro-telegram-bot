package main

import (
	"github.com/joho/godotenv"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/console"
)

func main() {
	godotenv.Load()
	console.Execute()
}
