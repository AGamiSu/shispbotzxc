package main

import (
    "fmt"
    "log"
    "net/http"
)

const donationToken = "MDmUS0LsuloNvUHonPmxLKhsKVxPSGuPCDocEDK9"

func main() {
    url := "https://www.donationalerts.com/api/v1/alerts/donations"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer "+donationToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        fmt.Println("Токен работает корректно!")
    } else {
        fmt.Printf("Ошибка: %s\n", resp.Status)
    }
}