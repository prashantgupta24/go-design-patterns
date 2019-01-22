package main

import (
	"errors"
	"fmt"
	"strconv"
)

type UserFinder interface {
	FindUser(id int32) (User, error)
}

type User struct {
	ID int32
}

type UserListProxy struct {
	SomeDatabase              *UserList
	StackCache                *UserList
	StackCapacity             int
	DidDidLastSearchUsedCache bool
}

func (u *UserListProxy) FindUser(id int32) (User, error) {

	var targetUser User
	proxyCache := u.StackCache

	for _, user := range *proxyCache {
		if user.ID == id {
			targetUser = *user
		}
	}

	if targetUser.ID != 0 {
		return targetUser, nil
	}

	if len(*proxyCache) >= u.StackCapacity {
		pc := *proxyCache
		for _, user := range pc {
			fmt.Printf(strconv.Itoa(int(user.ID)) + ", ")
		}
		fmt.Println()
		pc = pc[1:]
		for _, user := range pc {
			fmt.Printf(strconv.Itoa(int(user.ID)) + ", ")
		}
		proxyCache = &pc
	}

	*proxyCache = append(*proxyCache, &User{id})

	return User{}, errors.New("Could not find user in Cache")
}

type UserList []*User

func main() {

	db := UserList{}
	for i := 0; i < 1000000; i++ {
		db = append(db, &User{ID: int32(i)})
	}
	fmt.Println("The size of the db is :", len(db))

	proxy := &UserListProxy{
		SomeDatabase:  &db,
		StackCache:    &UserList{},
		StackCapacity: 10,
	}

	targetUsers := &UserList{}

	for i := 0; i < 12; i++ {
		//id := rand.Intn(10)
		*targetUsers = append(*targetUsers, &User{ID: int32(i)})
	}

	for _, user := range *targetUsers {
		fmt.Printf("Trying to find user %v in proxy?\n", user.ID)

		if _, err := proxy.FindUser(user.ID); err == nil {
			fmt.Println("User found in cache!")
		} else {
			fmt.Println("User not in cache")
		}
	}

	for _, user := range *proxy.StackCache {
		fmt.Printf(strconv.Itoa(int(user.ID)) + ", ")
	}
	fmt.Println()

}
