package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/lifevun-dev/oldg/graph/generated"
	"github.com/lifevun-dev/oldg/graph/model"
	"github.com/sirupsen/logrus"
)

func (r *mutationResolver) PinMessage(ctx context.Context, msg string, author string) (bool, error) {
	cmd := model.PinMessage{
		Msg:    msg,
		Author: author,
	}
	for _, obs := range observers {
		obs.msgChan <- cmd
	}
	return true, nil
}

func (r *mutationResolver) Unpin(ctx context.Context) (bool, error) {
	for _, obs := range observers {
		obs.msgChan <- model.Unpin{}
	}
	return true, nil
}

func (r *queryResolver) Hello(ctx context.Context, name *string) (string, error) {
	if name == nil {
		return "World!", nil
	}
	return fmt.Sprintf("Hello, %s!", *name), nil
}

func (r *subscriptionResolver) Commands(ctx context.Context) (<-chan model.Command, error) {
	obs := Observer{
		msgChan: make(chan model.Command, 10),
	}
	id := randString(10)
	observers[id] = obs

	go func() {
		<-ctx.Done() // subscription cancellata
		close(observers[id].msgChan)
		delete(observers, id)
		fmt.Printf("Subscription delete: %s\n", id)
	}()

	return obs.msgChan, nil
}

func (r *subscriptionResolver) TwitchChat(ctx context.Context, channel string) (<-chan *model.ChatMessage, error) {
	client := twitch.NewAnonymousClient()
	c := make(chan *model.ChatMessage, 1)
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		c <- &model.ChatMessage{
			Msg:    msg.Message,
			Author: msg.User.DisplayName,
		}
	})
	client.Join(channel)

	go func() {
		<-ctx.Done()
		client.Disconnect()
		close(c)
		logrus.Infof("client disconnected: %v", channel)
	}()

	go func() {
		err := client.Connect()
		if err != nil && !errors.Is(twitch.ErrClientDisconnected, err) {
			close(c)
			logrus.Errorf("error on channel %v: %v", channel, err)
		}
	}()

	return c, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
