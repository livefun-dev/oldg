package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/lifevun-dev/oldg/graph/generated"
	"github.com/lifevun-dev/oldg/graph/model"
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
type Observer struct {
	msgChan chan model.Command
}

var observers map[string]Observer = map[string]Observer{}

func (r *subscriptionResolver) NewMessage(ctx context.Context) (<-chan string, error) {
	msgChan := make(chan string, 10)

	go func() {
		cnt := 0
		for {
			msgChan <- fmt.Sprintf("msg: %d", cnt)
			cnt += 1
			time.Sleep(2 * time.Second)
		}
	}()

	return msgChan, nil
}
