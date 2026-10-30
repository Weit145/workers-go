package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Result[T any] struct {
	Data T
	Err  error
}

type User struct {
	name  string
	email string
}

func user(ctx context.Context, _ int, ch chan<- Result[User], wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Millisecond * 500):
	}
	// err := nil
	// err := errors.ErrUnsupported
	// if err != nil {
	// 	cancel()
	// }
	res := Result[User]{
		Data: User{
			name:  "Dima",
			email: "email",
		},
		Err: nil,
	}
	select {
	case <-ctx.Done():
	case ch <- res:
	}
}

type Order struct {
	name  string
	price float64
}

func order(ctx context.Context, _ int, ch chan<- Result[[]Order], wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Millisecond * 500):
	}
	// err := errors.ErrUnsupported
	// if err != nil {
	// 	cancel()
	// }
	orders := []Order{Order{
		name:  "Pen",
		price: 12,
	}, Order{
		name:  "Apple",
		price: 12000,
	}}

	res := Result[[]Order]{
		Data: orders,
		Err:  nil,
	}

	select {
	case <-ctx.Done():
	case ch <- res:
	}
}

func payment(ctx context.Context, _ int, ch chan<- Result[float64], wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Millisecond * 500):
	}
	err := errors.ErrUnsupported

	res := Result[float64]{
		Data: 1000,
		Err:  err,
	}
	select {
	case <-ctx.Done():
	case ch <- res:
	}
}

type Response struct {
	User    User
	Order   []Order
	Balance float64
}

func merge(ctx context.Context, cancel context.CancelFunc, chaUser chan Result[User], chanOrder chan Result[[]Order], chanBalance chan Result[float64]) (Response, error) {
	var response Response
	var hasUser, hasOrder, hasBalance bool
	for !(hasUser && hasOrder && hasBalance) {
		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		case user := <-chaUser:
			if user.Err != nil {
				cancel()
				return Response{}, user.Err
			}
			response.User = user.Data
			hasUser = true
		case order := <-chanOrder:
			if order.Err != nil {
				cancel()
				return Response{}, order.Err
			}
			response.Order = order.Data
			hasOrder = true
		case balance := <-chanBalance:
			if balance.Err != nil {
				cancel()
				return Response{}, balance.Err
			}
			response.Balance = balance.Data
			hasBalance = true
		}
	}
	return response, nil
}

func main() {

	chUser := make(chan Result[User])
	chOrder := make(chan Result[[]Order])
	chBalance := make(chan Result[float64])

	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	wg.Add(3)
	go user(ctx, 1, chUser, &wg)
	go order(ctx, 1, chOrder, &wg)
	go payment(ctx, 1, chBalance, &wg)

	response, err := merge(ctx, cancel, chUser, chOrder, chBalance)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		wg.Wait()
		return
	}
	fmt.Println("=== Данные профиля ===")
	fmt.Printf("Пользователь: %s\n", response.User.name)
	fmt.Printf("Баланс:        %.2f руб.\n", response.Balance)
	fmt.Printf("Всего заказов: %d\n", len(response.Order))
	wg.Wait()
}
