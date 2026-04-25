package main

import "context"

type nopSippy struct{}

func (nopSippy) Get(_ context.Context, _ string, _ map[string]string) ([]byte, error) {
	return nil, nil
}

type nopRC struct{}

func (nopRC) Get(_ context.Context, _ string, _ map[string]string) ([]byte, error) {
	return nil, nil
}

func (nopRC) GetForArch(_ context.Context, _, _ string, _ map[string]string) ([]byte, error) {
	return nil, nil
}

type nopSearch struct{}

func (nopSearch) Search(_ context.Context, _ string, _ map[string]string) ([]byte, error) {
	return nil, nil
}

func (nopSearch) Get(_ context.Context, _ string, _ map[string]string) ([]byte, error) {
	return nil, nil
}
