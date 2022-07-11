package main

// @title NFT Key Me API
// @version 1.0
// @description This is the API to query user's NFT data

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/reliablestaking/nftkeyme-discord/db"
	"github.com/reliablestaking/nftkeyme-discord/discord"
	"github.com/reliablestaking/nftkeyme-discord/nftkeyme"
	"github.com/reliablestaking/nftkeyme-discord/server"
	"golang.org/x/oauth2"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	// init database
	portInt := 5432
	port := os.Getenv("DB_PORT")
	if port != "" {
		portInt, _ = strconv.Atoi(port)
	}
	sslmode := "disable"
	if os.Getenv("DB_SSL") == "true" {
		sslmode = "require"
	}
	pgCon := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_ADDR"),
		portInt,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		sslmode)
	database, err := sqlx.Connect("postgres", pgCon)
	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to db...")
	}
	defer database.Close()
	store := db.Store{
		Db: database,
	}

	// init discord server
	//TODO: make configurable
	discordOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		Scopes:       []string{"identity"},
		Endpoint:     oauth2.Endpoint{TokenURL: os.Getenv("DISCORD_TOKEN_URL")},
	}

	nftkeymeOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("NFTKEYME_REDIRECT_URL"),
		ClientID:     os.Getenv("NFTKEYME_CLIENT_ID"),
		ClientSecret: os.Getenv("NFTKEYME_CLIENT_SECRET"),
		Scopes:       []string{"offline assets"},
		Endpoint: oauth2.Endpoint{
			TokenURL: os.Getenv("NFTKEYME_TOKEN_URL"),
			AuthURL:  os.Getenv("NFTKEYME_AUTH_URL"),
		},
	}

	discordBotToken := os.Getenv("DISCORD_BOT_TOKEN")
	if discordBotToken == "" {
		logrus.Fatalf("Discord bot token not found")
	}

	discordBot, err := discordgo.New("Bot " + discordBotToken)
	if err != nil {
		logrus.WithError(err).Fatal("Error setting up discord")
	}

	discordAuthURL := os.Getenv("DISCORD_AUTH_URL")
	if discordAuthURL == "" {
		logrus.Fatalf("Discord auth url not found")
	}

	policyIDCheck := os.Getenv("POLICY_ID_CHECK")
	if policyIDCheck == "" {
		logrus.Fatalf("Policy id check")
	}

	policyIDCheckHunter := os.Getenv("POLICY_ID_CHECK_HUNTERS")
	if policyIDCheckHunter == "" {
		logrus.Fatalf("Policy id check")
	}

	serverID := os.Getenv("DISCORD_SERVER_ID")
	if serverID == "" {
		logrus.Fatalf("Discrod server id check")
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		logrus.Fatalf("Channel id check")
	}

	roleMapString := os.Getenv("DISCORD_ROLE_MAP")
	if roleMapString == "" {
		logrus.Fatalf("Role map check")
	}
	roleMap, err := buildRoleMap(roleMapString)
	if err != nil {
		logrus.WithError(err).Fatalf("Error building role map %s", roleMapString)
	}

	// init server
	server := server.Server{
		Store:                store,
		Sha1ver:              sha1ver,
		BuildTime:            buildTime,
		DiscordOauthConfig:   discordOauthConfig,
		NftkeymeOauthConfig:  nftkeymeOauthConfig,
		DiscordClient:        discord.NewClientFromEnvironment(),
		NftkeymeClient:       nftkeyme.NewClientFromEnvironment(),
		DiscordSession:       discordBot,
		DiscordAuthCodeURL:   discordAuthURL,
		PolicyIDCheck:        policyIDCheck,
		PolicyIDCheckHunters: policyIDCheckHunter,
		DiscordServerID:      serverID,
		DiscordChannelID:     channelID,
		RoleMap:              roleMap,
	}

	// start monitor
	go server.VerifyAccess()

	// start server
	server.Start()
}

func buildRoleMap(roleMapString string) (map[int]string, error) {
	roleMap := make(map[int]string)

	roleMapSplits := strings.Split(roleMapString, ",")
	for _, roleMapSplit := range roleMapSplits {
		roleMapArray := strings.Split(roleMapSplit, ":")
		numberValue, err := strconv.Atoi(roleMapArray[0])
		if err != nil {
			return nil, err
		}
		roleMap[numberValue] = roleMapArray[1]
	}

	return roleMap, nil
}
