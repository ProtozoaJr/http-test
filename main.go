package main

import (
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/robfig/cron/v3"
    "github.com/joho/godotenv"
)


func main() {
    loadEnv()

    cronSchedule := getEnv("CRON_SCHEDULE", "*/10 * * * * *")
    apiUrl := getEnv("API_URL", "https://github.com")

    fmt.Println("START HTTP TEST")

    c := cron.New()
    _, err := c.AddFunc(cronSchedule, func() {
        fmt.Println("Run HTTP TEST schedule job on:", time.Now().Format(time.RFC3339))
        hitApi(apiUrl)
    })
    if err != nil {
        log.Fatal("Error adding cron job:", err)
    }

    c.Start()

    select {} // Block forever
}

func hitApi(url string) {
    apiData, err := getAPI(url)
    if err != nil {
        fmt.Printf("Failed to get API data: %v - %s\n", time.Now().Format(time.RFC3339), err)
        return
    }
    fmt.Printf("Success on get API data: (%s) - %s\n", time.Now().Format(time.RFC3339), apiData)
}

func getAPI(url string) (string, error) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }
}

func getEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return fallback
    }
    return value
}
