// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Command interface {
	IsCommand()
}

type PinMessage struct {
	Msg    string `json:"msg"`
	Author string `json:"author"`
}

func (PinMessage) IsCommand() {}

type Unpin struct {
	B *bool `json:"b"`
}

func (Unpin) IsCommand() {}
