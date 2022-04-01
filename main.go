package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
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
		okta.WithPrivateKey("./secrets/private.pem"),
		okta.WithCache(false),
	)

	if err != nil {
		handleError(err)
	}

	actionSelector("\n\n\n\n\n\n\n\n\n\n\n******** MAIN MENU ********\n\nWhat would you like to do?", ctx, client)
}

func ListGroups(ctx context.Context, client *okta.Client) error {
	oktaGroups, oktaResponse, err := client.Group.ListGroups(ctx, nil)

	if err != nil {
		return err
	}

	if oktaResponse.StatusCode != 200 {
		fmt.Printf("error in status: %s", oktaResponse.Status)
	}

	fmt.Printf("\n\n********** OKTA GROUPS:LIST RESPONSE **************\n\nStatus: %s\nGroups:\n", oktaResponse.Status)

	for _, g := range oktaGroups {
		fmt.Printf("\nID: %s\nName: %s\nDescription: %s\n", g.Id, g.Profile.Name, g.Profile.Description)
	}

	return nil
}

func actionSelector(label string, ctx context.Context, client *okta.Client) {
	var s string
	r := bufio.NewReader(os.Stdin)
	exit := false
	options := `

	1) create a random Okta Group with name 'Random Group {randomNumber}'
	2) list all Okta Groups
	3) delete an Okta Group (choose from list)

	type <exit> to exit the program

	`
	makeSelectionStatement := "\n\nPlease enter a valid choice and press <Enter>"
	wrongChoiceStatement := "\n\nSorry, but you didn't enter a valid choice!"

	for !exit {
		fmt.Fprint(os.Stderr, label+options+makeSelectionStatement)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		switch s {
		case "1":
			createRandomGroup(ctx, client)
		case "2":
			ListGroups(ctx, client)
		case "3":
			deleteGroupSelector(ctx, client)
		case "exit":
			exit = true
		default:
			fmt.Fprint(os.Stderr, wrongChoiceStatement)
			continue
		}
	}
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

	fmt.Printf("\n\n********** OKTA GROUP:CREATE RESPONSE **************\n\nOkta Response Status: %s\nOkta Group:\n%+s", oktaResponse.Status, prettyPrint(oktaGroup.Profile))
}

func deleteGroupSelector(ctx context.Context, client *okta.Client) {
	var s string
	r := bufio.NewReader(os.Stdin)
	exit := false
	oktaGroups, oktaResponse, err := client.Group.ListGroups(ctx, nil)

	if err != nil {
		handleError(err)
	}

	if oktaResponse.StatusCode != 200 {
		fmt.Fprint(os.Stderr, "there was an error retrieving groups")
	}

	groupSelectorMap := map[string]string{}
	groupsInfo := []string{}
	for i, g := range oktaGroups {
		groupSelectorMap[strconv.Itoa(i)] = g.Id
		groupInfo := fmt.Sprintf(fmt.Sprintf("%d) Name: %s\nID:%s\nDescription:%s", i, g.Profile.Name, g.Id, g.Profile.Description))
		groupsInfo = append(groupsInfo, groupInfo)
	}

	makeSelectionStatement := "\n\nWhich group would you like to delete? Make a valid choice and press <Enter>\n\n"
	groupsInfoDisplay := strings.Join(groupsInfo, "\n\n")
	exitStatement := "\n\nType <exit> to abandon\n\n"
	wrongChoiceStatement := "\n\nSorry, but you didn't enter a valid choice! The prompt will be displayed again.\n\n"

	exit = false

	for !exit {
		fmt.Fprint(os.Stderr, groupsInfoDisplay+makeSelectionStatement+exitStatement)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if groupId, validChoice := groupSelectorMap[s]; validChoice {
			confirmationStatement := fmt.Sprintf("\n\nYou've chosen to delete the group with id %s - are you sure? y/n\n", groupId)
			fmt.Fprint(os.Stderr, confirmationStatement)
			s, _ = r.ReadString('\n')
			s = strings.TrimSpace(s)
			if s == "y" {
				deleteGroup(groupId, ctx, client)
			}
			exit = true
		} else if s == "exit" {
			exit = true
		} else {
			fmt.Fprint(os.Stderr, wrongChoiceStatement)
			time.Sleep(2 * time.Second)
		}
	}
}

func deleteGroup(id string, ctx context.Context, client *okta.Client) {
	oktaResponse, err := client.Group.DeleteGroup(ctx, id)

	if err != nil {
		handleError(err)
	}

	if oktaResponse.StatusCode != 204 {
		fmt.Printf("\n\n********** OKTA GROUP:DELETE RESPONSE **************\n\nerror deleting group with id %s", id)
	} else {
		fmt.Printf("\n\n********** OKTA GROUP:DELETE RESPONSE **************\n\ndeleted group with id %s", id)
	}
}

func handleError(err error) {
	fmt.Printf("Error: %v\n", err)
}

func prettyPrint(v interface{}) string {
	res, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		fmt.Printf("error doing marshalIndent on struct: %+v", v)
	}

	return string(res)
}
