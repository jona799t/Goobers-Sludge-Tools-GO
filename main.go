package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func Colorate(str string, from pterm.RGB, to pterm.RGB) string {
	strs := strings.Split(str, "")
	var fadeInfo string
	for i := 0; i < len(str); i++ {
		fadeInfo += from.Fade(0, float32(len(str)), float32(i), to).Sprint(strs[i])
	}

	return fadeInfo + "\u001B[38;5;255m"
}

func title(from pterm.RGB, to pterm.RGB) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()

	} else {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()

	}

	fmt.Println(Colorate("\n      ______            __                  ", from, to) + "  _       ____ ")
	fmt.Println(Colorate("     / ____/___  ____  / /_  ___  __________", from, to) + " | |     / / / ")
	fmt.Println(Colorate("    / / __/ __ \\/ __ \\/ __ \\/ _ \\/ ___/ ___/", from, to) + " | | /| / / /  ")
	fmt.Println(Colorate("   / /_/ / /_/ / /_/ / /_/ /  __/ /  (__  ) ", from, to) + " | |/ |/ / /___")
	fmt.Println(Colorate("   \\____/\\____/\\____/_.___/\\___/_/  /____/  ", from, to) + " |__/|__/_____/\n")
}

func accountDetails(client *fasthttp.Client, accounts []*fasthttp.Request, green pterm.RGB, cyan pterm.RGB) {
	title(green, cyan)

	i := 0
	for _, account := range accounts {
		i += 1
		account.SetRequestURIBytes([]byte("https://app.goobers.net/api/users/@me"))

		resp := fasthttp.AcquireResponse()
		client.Do(account, resp)

		fmt.Printf("%s %s\n", "   "+Colorate(strconv.Itoa(i)+".", green, cyan), string(resp.Body()))
	}

	fmt.Printf("\n   Click " + Colorate("ENTER", green, cyan) + " to go back")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	main()
}

func purchaser(client *fasthttp.Client, accounts []*fasthttp.Request, green pterm.RGB, cyan pterm.RGB) {
	title(green, cyan)

	req := fasthttp.AcquireRequest()

	req.Header.Set("accept", "application/json, text/plain, */*")
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-DK,en;q=0.9,da-DK;q=0.8,da;q=0.7,en-US;q=0.6")
	req.Header.Set("referer", "gzip, deflate, br")
	req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"100\", \"Google Chrome\";v=\"100\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")

	req.SetRequestURIBytes([]byte("https://app.goobers.net/api/shop/items"))

	resp := fasthttp.AcquireResponse()
	client.Do(req, resp)
	var jsonResp []map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &jsonResp); err != nil {
		log.Fatal(err)
	}
	i := 0
	for _, element := range jsonResp {
		i += 1
		fmt.Printf("%s %s (%.0f sludge)\n", "   "+Colorate(strconv.Itoa(i)+".", green, cyan), element["title"], element["price"])
	}
	var choice int
	fmt.Printf("\n   Select an option: ")
	fmt.Scan(&choice)

	item := jsonResp[choice-1]

	fmt.Printf("\n   Item selected: %s\n", Colorate(fmt.Sprintf("%s", item["title"]), green, cyan))
	fmt.Printf("   Price for %s: %s\n", item["title"], Colorate(fmt.Sprintf("%.0f", item["price"]), green, cyan))

	//fmt.Printf("\n   Max %s per user is %s", item["title"], Colorate(fmt.Sprintf("%.0f", item["max_per_user"]), green, cyan))

	var timeDelay float64
	fmt.Printf("\n   How much %s would you like to between each request? (Recommended: 1-2): ", Colorate("time delay", green, cyan))
	fmt.Scan(&timeDelay)

	fmt.Printf("\n   Click " + Colorate("ENTER", green, cyan) + " to go back")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	main()
}

func main() {
	client := &fasthttp.Client{
		//Dial: fasthttpproxy.FasthttpHTTPDialer("username:pass@host:port"),
	}

	green := pterm.NewRGB(0, 255, 0)
	cyan := pterm.NewRGB(0, 255, 255)

	title(green, cyan)

	var accounts []*fasthttp.Request

	accountsByte, _ := ioutil.ReadFile("data/accounts.txt")
	for _, account := range strings.Split(string(accountsByte), "\n") {
		req := fasthttp.AcquireRequest()

		req.Header.Set("accept", "application/json, text/plain, */*")
		//req.Header.Set("accept-encoding", "gzip, deflate, br")
		req.Header.Set("accept-language", "en-DK,en;q=0.9,da-DK;q=0.8,da;q=0.7,en-US;q=0.6")
		req.Header.Set("authorization", account)
		req.Header.Set("cookie", "authorization="+account+"; goobers-lasting-token="+account)
		req.Header.Set("referer", "gzip, deflate, br")
		req.Header.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"100\", \"Google Chrome\";v=\"100\"")
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")

		accounts = append(accounts, req)
	}

	fmt.Printf("%s Account Details\n", "   "+Colorate("1.", green, cyan))
	fmt.Printf("%s Purchaser\n", "   "+Colorate("2.", green, cyan))

	var choice int
	fmt.Printf("\n   Select an option: ") // Question to the user
	fmt.Scan(&choice)

	if choice == 1 {
		accountDetails(client, accounts, green, cyan)
	} else if choice == 2 {
		purchaser(client, accounts, green, cyan)
	} else {
		fmt.Printf("\n   %d is an invalid choice\n", choice)
	}
	fmt.Println()
}
