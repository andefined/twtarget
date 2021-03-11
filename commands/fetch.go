package commands

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/andefined/twtarget/utils"
	"github.com/urfave/cli"
)

var (
	target *utils.Target
)

// FollowersPage data
type FollowersPage struct {
	Followers  []anaconda.User
	NextCursor string
	Error      error
}

// FriendsPage data
type FriendsPage struct {
	Friends    []anaconda.User
	NextCursor string
	Error      error
}

// Fetch : Fetch User Data
func Fetch(c *cli.Context) error {
	if c.Args().Get(0) == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	target = utils.GetTarget(c.Args().Get(0))
	api := utils.API(target.Credentials)

	user, err := api.GetUsersShow(target.User, url.Values{})
	utils.ExitOnError(err)

	// fetch user
	if c.Bool("user") {
		SaveUserInfo(user.ScreenName, []string{
			user.CreatedAt,
			user.IdStr,
			user.ScreenName,
			user.Name,
			utils.CleanText(user.Description),
			"https://twitter.com/" + user.ScreenName,
			strconv.FormatBool(user.DefaultProfileImage),
			strconv.Itoa(user.FavouritesCount),
			strconv.Itoa(user.FollowersCount),
			strconv.Itoa(user.FriendsCount),
			strconv.FormatInt(user.StatusesCount, 10),
			strconv.FormatInt(user.ListedCount, 10),
			strconv.FormatBool(user.Verified),
			"",
			user.ProfileImageUrlHttps,
		})
	}

	// fetch followers
	if c.Bool("followers") {
		fmt.Printf("WARNING: You need %v to complete this request. Continue? (y/n): ",
			time.Duration(float64(user.FollowersCount/200))*time.Minute,
		)
		if utils.PromptConfirm() {
			GetFollowers(api, user.FollowersCount)
		}
	}

	// fetch friends
	if c.Bool("friends") {
		fmt.Printf("WARNING: You need %v to complete this request. Continue? (y/n): ",
			time.Duration(float64(user.FriendsCount/200))*time.Minute,
		)
		if utils.PromptConfirm() {
			GetFriends(api, user.FriendsCount)
		}
	}

	return nil
}

// SaveUserInfo : Create a single csv holding user information
func SaveUserInfo(user string, record []string) {
	o := target.Path + "/target.csv"
	if _, err := os.Stat(o); err != nil {
		utils.CreateFile(o)
	}

	f, err := os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	utils.ExitOnError(err)
	defer f.Close()

	writer := csv.NewWriter(f)

	writer.Write([]string{
		"CreatedAt",
		"Id",
		"ScreenName",
		"Name",
		"Description",
		"Link",
		"DefaultProfileImage",
		"FavouritesCount",
		"FollowersCount",
		"FriendsCount",
		"StatusesCount",
		"ListedCount",
		"Verified",
		"Status",
		"ProfileImage",
	})

	writer.Write(record)
	writer.Flush()
}

// GetFollowers : Retrieve Followers
func GetFollowers(api *anaconda.TwitterApi, c int) {
	o := target.Path + "/followers.csv"

	if target.FollowersNextCursor == "-1" {
		os.Remove(o)
	}

	if _, err := os.Stat(o); err != nil {
		utils.CreateFile(o)
	}

	f, err := os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	utils.ExitOnError(err)
	defer f.Close()

	writer := csv.NewWriter(f)

	v := url.Values{}
	v.Set("count", "200")
	v.Set("screen_name", target.User)
	v.Set("skip_status", "false")

	nextCursor := "-1"
	if target.FollowersNextCursor != "-1" {
		nextCursor = target.FollowersNextCursor
	} else {
		writer.Write([]string{
			"CreatedAt",
			"Id",
			"ScreenName",
			"Name",
			"Description",
			"Link",
			"DefaultProfileImage",
			"FavouritesCount",
			"FollowersCount",
			"FriendsCount",
			"StatusesCount",
			"ListedCount",
			"Verified",
			"Status",
			"ProfileImage",
		})
		writer.Flush()
	}

	chFollowers := make(chan FollowersPage)
	go func(a anaconda.TwitterApi, v url.Values, next_cursor string, result chan FollowersPage) {
		// Cursor defaults to the first page ("-1")
		next_cursor = "-1"
		for {
			v.Set("cursor", next_cursor)
			c, err := a.GetFollowersList(v)
			result <- FollowersPage{c.Users, c.Next_cursor_str, err}

			next_cursor = c.Next_cursor_str
			if err != nil || next_cursor == "0" {
				close(result)
				break
			}
		}
	}(*api, v, nextCursor, chFollowers)

	for {
		data, more := <-chFollowers
		c -= len(data.Followers)
		fmt.Printf("Remaining Followers: %v, Next Cursor: %s\n", c, data.NextCursor)
		// save next cursor
		target.FollowersNextCursor = data.NextCursor
		utils.SaveTarget(target)
		for _, user := range data.Followers {
			status := ""

			if user.StatusesCount > 0 && !user.Protected && user.Status != nil {
				status = utils.CleanText(user.Status.FullText)
			}

			writer.Write([]string{
				user.CreatedAt,
				user.IdStr,
				user.ScreenName,
				user.Name,
				utils.CleanText(user.Description),
				"https://twitter.com/" + user.ScreenName,
				strconv.FormatBool(user.DefaultProfileImage),
				strconv.Itoa(user.FavouritesCount),
				strconv.Itoa(user.FollowersCount),
				strconv.Itoa(user.FriendsCount),
				strconv.FormatInt(user.StatusesCount, 10),
				strconv.FormatInt(user.ListedCount, 10),
				strconv.FormatBool(user.Verified),
				status,
				user.ProfileImageUrlHttps,
			})
			writer.Flush()
		}
		if !more {
			break
		}
	}
}

// GetFriends : Retrieve Friends
func GetFriends(api *anaconda.TwitterApi, c int) {
	o := target.Path + "/friends.csv"

	if target.FriendsNextCursor == "-1" {
		os.Remove(o)
	}

	if _, err := os.Stat(o); err != nil {
		utils.CreateFile(o)
	}

	f, err := os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	utils.ExitOnError(err)
	defer f.Close()

	writer := csv.NewWriter(f)

	v := url.Values{}
	v.Set("count", "200")
	v.Set("screen_name", target.User)
	v.Set("skip_status", "false")

	nextCursor := "-1"
	if target.FriendsNextCursor != "-1" {
		nextCursor = target.FriendsNextCursor
	} else {
		writer.Write([]string{
			"CreatedAt",
			"Id",
			"ScreenName",
			"Name",
			"Description",
			"Link",
			"DefaultProfileImage",
			"FavouritesCount",
			"FollowersCount",
			"FriendsCount",
			"StatusesCount",
			"ListedCount",
			"Verified",
			"Status",
			"ProfileImage",
		})
		writer.Flush()
	}

	chFriends := make(chan FriendsPage)
	go func(a anaconda.TwitterApi, v url.Values, next_cursor string, result chan FriendsPage) {
		// Cursor defaults to the first page ("-1")
		next_cursor = "-1"
		for {
			v.Set("cursor", next_cursor)
			c, err := a.GetFriendsList(v)
			result <- FriendsPage{c.Users, c.Next_cursor_str, err}

			next_cursor = c.Next_cursor_str
			if err != nil || next_cursor == "0" {
				close(result)
				break
			}
		}
	}(*api, v, nextCursor, chFriends)
	for {
		data, more := <-chFriends
		c -= len(data.Friends)
		fmt.Printf("Remaining Friends: %v, Next Cursor: %s\n", c, data.NextCursor)
		// save next cursor
		target.FriendsNextCursor = data.NextCursor
		utils.SaveTarget(target)
		for _, user := range data.Friends {
			status := ""

			if user.StatusesCount > 0 && !user.Protected && user.Status != nil {
				status = utils.CleanText(user.Status.FullText)
			}
			writer.Write([]string{
				user.CreatedAt,
				user.IdStr,
				user.ScreenName,
				user.Name,
				utils.CleanText(user.Description),
				"https://twitter.com/" + user.ScreenName,
				strconv.FormatBool(user.DefaultProfileImage),
				strconv.Itoa(user.FavouritesCount),
				strconv.Itoa(user.FollowersCount),
				strconv.Itoa(user.FriendsCount),
				strconv.FormatInt(user.StatusesCount, 10),
				strconv.FormatInt(user.ListedCount, 10),
				strconv.FormatBool(user.Verified),
				status,
				user.ProfileImageUrlHttps,
			})
			writer.Flush()
		}
		if !more {
			break
		}
	}
}
