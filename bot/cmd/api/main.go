package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := sql.Open("postgres", `host=db port=5432 user=postgres 
						password=12345 dbname=maxhack sslmode=disable`)
	if err != nil {
		log.Fatalf("DB connect error:%s", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	log.Printf("Token length: %d", len(os.Getenv("TOKENBOT")))
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("DB ping error:%s", err)
	}
	api, err := maxbot.New(os.Getenv("TOKENBOT"))
	if err != nil {
		log.Fatalf("Failed to create bot API: %v", err)
	}

	info, err := api.Bots.GetBot(context.Background())
	if err != nil {
		log.Fatalf("Failed to get bot info: %v", err)
	}

	fmt.Printf("Get me: %#v %#v", info, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		<-exit
		cancel()
	}()

	for upd := range api.GetUpdates(ctx) {
		switch upd := upd.(type) {
		case *schemes.BotStartedUpdate:
			go CreateUser(db, *upd)
			_, _ = api.Messages.Send(context.Background(), maxbot.NewMessage().SetChat(upd.ChatId).SetText("Привет спортсмен!\nСкорее открывай мини приложение и контролируй свой рацион!"))
		case *schemes.MessageCreatedUpdate:
			switch upd.GetCommand() {
			case "/start":
				_, _ = api.Messages.Send(context.Background(), maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText("Привет спортсмен!\nСкорее открывай мини приложение и контролируй свой рацион!"))
			}
		}
	}
}

func CreateUser(db *sql.DB, upd schemes.BotStartedUpdate) error {
	_, err := db.Exec(`
        INSERT INTO users (id, name) 
        VALUES ($1, $2)
        ON CONFLICT (id) DO NOTHING
    `, upd.GetUserID(), upd.User.FirstName)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
