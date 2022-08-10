package console

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/config"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/db"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/external/himatro"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/handler"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/repository"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var botCmd = &cobra.Command{
	Use:     "bot",
	Short:   "start bot",
	Long:    "start a new bot server",
	Example: "go run main.go bot",
	Run:     runBot,
}

func init() {
	RootCmd.AddCommand(botCmd)
}

func runBot(cmd *cobra.Command, args []string) {
	db.InitializePostgresConn()

	himatroClient := himatro.NewClient(config.HimatroAPIHost())
	sessionRepo := repository.NewSessionRepo(db.PostgresDB)

	userRepo := repository.NewUserRepository(db.PostgresDB)
	userUsecase := usecase.NewUserUsecase(userRepo, himatroClient, sessionRepo)

	handler := handler.New(userUsecase)

	b, err := gotgbot.NewBot(config.Token(), &gotgbot.BotOpts{
		Client:            http.Client{},
		DisableTokenCheck: false,
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: time.Second * 20,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})

	if err != nil {
		logrus.Panic("error creating bot instance: %v", err)
	}

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: log.Default(),
		DispatcherOpts: ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				fmt.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		},
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("register", handler.RegisterHandler()))
	dispatcher.AddHandler(handlers.NewCommand("login", handler.LoginHandler()))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
