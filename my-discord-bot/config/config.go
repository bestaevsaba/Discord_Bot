package config
import (
    "encoding/json"
    "fmt"
    "os"
)

var (
    Token string
    BotPrefix string
    OpenAIApiKey string
    Hugobot string
    MongoDBURL string

    config *configStruct
)

type configStruct struct {
    Token string
    BotPrefix string
    OpenAIApiKey string
    Hugobot string
    MongoDBURL string
}

func ReadConfig() error {
    file, err := os.ReadFile("config/config.json")
    if err != nil {
        fmt.Println(err)
    }
    err = json.Unmarshal(file, &config)
    if err != nil {
        fmt.Println(err)
    }
    Token = config.Token
    BotPrefix = config.BotPrefix
    OpenAIApiKey = config.OpenAIApiKey
    Hugobot = config.Hugobot
    MongoDBURL = config.MongoDBURL

    return nil
}