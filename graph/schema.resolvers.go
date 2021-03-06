package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph/generated"
	"github.com/minskylab/supersense/graph/model"
	"github.com/pkg/errors"
)

func (r *mutationResolver) Emit(ctx context.Context, token string, draft model.EventDraft) (string, error) {
	if r.spokesman == nil {
		return "", errors.New("you cannot emit because Spokesman is note enabled")
	}

	// if token == {}

	url := ""
	if draft.ShareURL != nil {
		url = *draft.ShareURL
	}

	entities := draftEntitiesToSSEntities(draft)

	if draft.Actor != nil {
		r.spokesman.BroadcastWithActor(
			draft.Actor.Name,
			draft.Actor.Username,
			draft.Actor.Photo,
			draft.Title, draft.Message, entities, url, nil,
		)
	} else {
		r.spokesman.Broadcast(draft.Title, draft.Message, entities, url, nil)
	}

	return draft.Message, nil
}

func (r *queryResolver) SharedBoard(ctx context.Context, buffer int) ([]*supersense.Event, error) {
	if buffer < 1 || buffer > 100 { // hard coded limit
		buffer = 100
	}

	currentState, err := r.store.CurrentSharedState(int64(buffer))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return currentState.Board, nil
}

func (r *queryResolver) Header(ctx context.Context) (*model.SuperHeader, error) {
	if r.conf == nil {
		return nil, errors.New("header settings not found")
	}

	return &model.SuperHeader{
		Title:      r.conf.ObserverTitle,
		Hashtag:    r.conf.ObserverHashtag,
		Brand:      r.conf.ObserverBrand,
		Buffer:     r.conf.ObserverBuffer,
		DarkColor:  r.conf.ObserverColorDark,
		LightColor: r.conf.ObserverColorLight,
	}, nil
}

func (r *subscriptionResolver) EventStream(ctx context.Context, filter *model.EventStreamFilter) (<-chan *supersense.Event, error) {
	pipe := make(chan *supersense.Event)

	if filter != nil {
		go r.mux.Register(pipe, ctx.Done(), filter.Sources...)
	} else {
		go r.mux.Register(pipe, ctx.Done())
	}

	return pipe, nil
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
func (r *queryResolver) Login(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented"))
}
