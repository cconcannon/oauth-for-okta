package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("error loading .env file")
	}

	ctx, client, err := okta.NewClient(
		context.TODO(),
		okta.WithOrgUrl(os.Getenv("OKTA_ORG_URL")),
		okta.WithAuthorizationMode("PrivateKey"),
		okta.WithClientId(os.Getenv("CLIENT_ID")),
		okta.WithScopes(([]string{"okta.groups.manage"})),
		okta.WithPrivateKey("./secrets/key_private.pem"),
	)

	if err != nil {
		handleError(err)
	}

	actionSelector("What would you like to do?", ctx, client)
}

func actionSelector(label string, ctx context.Context, client *okta.Client) {
	var s string
	r := bufio.NewReader(os.Stdin)
	validChoice := false
	options := `
	1) create an Okta Group
	2) list Okta Groups
	3) delete an Okta Group
	`
	for !validChoice {
		fmt.Fprint(os.Stderr, label+"\n\n"+options)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		switch s {
		case "1":
			createRandomGroup(ctx, client)
			validChoice = true
		case "2":
			listGroups(ctx, client)
			validChoice = true
		case "3":
			deleteGroupSelector(label, ctx, client)
			validChoice = true
		default:
			fmt.Fprint(os.Stderr, "Sorry, but you didn't enter a valid choice! Your options are:\n"+options)
			continue
		}
	}
}

func deleteGroupSelector(label string, ctx context.Context, client *okta.Client) {
	fmt.Printf("To be continued...")
}

func handleError(err error) {
	fmt.Printf("Error: %v\n", err)
}

func createRandomGroup(ctx context.Context, client *okta.Client) {
	rand.Seed(time.Now().UnixNano())

	group := okta.Group{
		Profile: &okta.GroupProfile{
			Name:        fmt.Sprintf("Random Group %d", rand.Intn(1000000)),
			Description: "Random Group Generation for Demo",
		},
	}

	oktaGroup, oktaResponse, err := client.Group.CreateGroup(ctx, group)

	if err != nil {
		handleError(err)
	}

	fmt.Printf("Okta Response Status: %s\nOkta Group:\n%+s", oktaResponse.Status, prettyPrint(oktaGroup.Profile))
}

func listGroups(ctx context.Context, client *okta.Client) error {
	oktaGroups, oktaResponse, err := client.Group.ListGroups(ctx, nil)

	if err != nil {
		return err
	}

	fmt.Printf("Okta List Groups Call Status: %s\nGroups:\n", oktaResponse.Status)

	for _, g := range oktaGroups {
		fmt.Printf("%s", prettyPrint(g))
	}

	return nil
}

func prettyPrint(v interface{}) string {
	res, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		fmt.Printf("error doing marshalIndent on struct: %+v", v)
	}

	return string(res)
}
