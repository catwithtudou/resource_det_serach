package data

import (
	"context"
	"gorm.io/gorm"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/utils"
	"testing"
)

func newUserRepoTest(t *testing.T) (*userRepo, context.Context) {
	data, _ := newData(t)
	return &userRepo{data: data}, context.Background()
}

func TestInsertUser(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	err := u.InsertUser(ctx, &biz.User{
		Name:   "zky",
		Email:  "1@qq.com",
		Pswd:   "1",
		Role:   "学生",
		Sex:    "男",
		School: "CQUPT",
		Sid:    "2018211113",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	user, err := u.GetUserByEmail(ctx, "111@qq.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(utils.JsonToString(user))
}

func TestGetUserBySid(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	user, err := u.GetUserBySid(ctx, "2018211111")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(utils.JsonToString(user))
}

func TestGetUserById(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	user, err := u.GetUserById(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(utils.JsonToString(user))
}

func TestListUser(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	users, err := u.ListUser(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range users {
		t.Logf(utils.JsonToString(v))
	}
}

func TestUpdateUser(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	user := &biz.User{
		Model: gorm.Model{ID: 7},
		//Pswd:     "111111111",
		//Avatar:   "avatar",
		//Intro:    "intro",
		IsActive: false,
	}
	t.Logf(utils.JsonToString(user))
	err := u.UpdateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	reUser, err := u.GetUserById(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(utils.JsonToString(reUser))

}

func TestDeleteUser(t *testing.T) {
	u, ctx := newUserRepoTest(t)
	err := u.DeleteUser(ctx, 6)
	if err != nil {
		t.Fatal(err)
	}
}
