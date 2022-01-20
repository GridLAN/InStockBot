package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID         = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken        = os.Getenv("TOKEN")
	WatchedProducts = []string{"G4 Doorbell Pro", "Camera G4 Instant"}
	ChannelID       = "933566907666292768"
)

func getJson(url string, target interface{}) error {
	// Create a new HTTP client with a timeout of 10 seconds
	var myClient = http.Client{Timeout: 10 * time.Second}

	// Build a GET request to the meal API endpoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add the required headers to the request
	req.Header.Set("User-Agent", "InStockBot")

	// Send the request and store the response
	r, getErr := myClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Close the response body when done
	defer r.Body.Close()

	// Read the response body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON response into the target interface
	err = json.Unmarshal(b, &target)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
	}

	return json.Unmarshal(b, target)
}

type UbiquitiProducts struct {
	Products []struct {
		ID          int64    `json:"id"`
		Title       string   `json:"title"`
		Handle      string   `json:"handle"`
		BodyHTML    string   `json:"body_html"`
		PublishedAt string   `json:"published_at"`
		CreatedAt   string   `json:"created_at"`
		UpdatedAt   string   `json:"updated_at"`
		Vendor      string   `json:"vendor"`
		ProductType string   `json:"product_type"`
		Tags        []string `json:"tags"`
		Variants    []struct {
			ID               int64       `json:"id"`
			Title            string      `json:"title"`
			Option1          string      `json:"option1"`
			Option2          interface{} `json:"option2"`
			Option3          interface{} `json:"option3"`
			Sku              string      `json:"sku"`
			RequiresShipping bool        `json:"requires_shipping"`
			Taxable          bool        `json:"taxable"`
			FeaturedImage    interface{} `json:"featured_image"`
			Available        bool        `json:"available"`
			Price            string      `json:"price"`
			Grams            int         `json:"grams"`
			CompareAtPrice   interface{} `json:"compare_at_price"`
			Position         int         `json:"position"`
			ProductID        int64       `json:"product_id"`
			CreatedAt        string      `json:"created_at"`
			UpdatedAt        string      `json:"updated_at"`
		} `json:"variants"`
		Images []struct {
			ID         int64         `json:"id"`
			CreatedAt  string        `json:"created_at"`
			Position   int           `json:"position"`
			UpdatedAt  string        `json:"updated_at"`
			ProductID  int64         `json:"product_id"`
			VariantIds []interface{} `json:"variant_ids"`
			Src        string        `json:"src"`
			Width      int           `json:"width"`
			Height     int           `json:"height"`
		} `json:"images"`
		Options []struct {
			Name     string   `json:"name"`
			Position int      `json:"position"`
			Values   []string `json:"values"`
		} `json:"options"`
	} `json:"products"`
}

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	if BotToken == "" {
		log.Fatal("Token cannot be empty")
	}

	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// message the channel to inform that bot is up and monitoring the WatchedProducts
	_, _ = s.ChannelMessageSend(ChannelID, "Bot is up and monitoring the following products: "+strings.Join(WatchedProducts, ", "))

	// check Ubiquiti stock every 30 seconds
	for {
		checkStock()
		time.Sleep(30 * time.Second)
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}

func checkStock() {
	var target UbiquitiProducts
	getJson("https://store.ui.com/products.json", &target)
	for _, product := range target.Products {
		// if the product is a watched product
		if contains(WatchedProducts, product.Title) {
			// and if the product is available
			if product.Variants[0].Available {
				// message the channel that the product is available
				s.ChannelMessageSend(ChannelID, fmt.Sprintf("%v is in stock!", product.Title))
				log.Printf(product.Title + " is in stock!")
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
