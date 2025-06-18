package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/benkiemle/gogator/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	handlers map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

func (cmds *commands) run(s *state, cmd command) error {
	handlerFunction, ok := cmds.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("command \"%s\" does not exist", cmd.name)
	}

	return handlerFunction(s, cmd)
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.handlers[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func GetCommands() commands {
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))
	return cmds
}

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("ERROR: expected argument `username`")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	if user.Name != cmd.args[0] {
		return fmt.Errorf("ERROR: user does not exist")
	}

	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("User has been set to '%s'\n", user.Name)
	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("ERROR: expected argument `name`")
	}

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})

	if err != nil {
		return err
	}

	s.config.SetUser(newUser.Name)

	fmt.Println("New user", newUser.Name, "created at", newUser.CreatedAt, "updated at", newUser.UpdatedAt, "with id", newUser.ID)
	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err == nil {
		fmt.Println("users table has been successfully reset")
	}
	return err
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		isCurrent := user.Name == s.config.CurrentUserName
		var name string
		if isCurrent {
			name = user.Name + " (current)"
		} else {
			name = user.Name
		}
		fmt.Printf("* %s\n", name)
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("ERROR: expected argument `time_between_reqs`")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s.db)
		if err != nil {
			return err
		}
	}
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("ERROR: expected arguments `name` and `url`")
	}

	usr, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    usr.ID,
	})

	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil

}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsView(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println("Feed", feed.Name, "with url", feed.Url, "created by", feed.UserName)
	}

	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("ERROR: expected argument `url`")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Feed:", feedFollow.FeedName, "User:", feedFollow.UserName)

	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println("-", feed.FeedName)
	}

	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("ERROR: expected argument `url`")
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	return err
}

func handleBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	limit = 2
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err == nil {
			limit = int32(parsedLimit)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%+v\n", post)
	}

	return nil
}
