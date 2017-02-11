package controllers

import (
    "strings"
    "github.com/astaxie/beego/logs"
)

type FrontController struct {
    baseController
}

type UnreadStat struct {
    Friend string `json:"friend"`
    Unread int `json:"unread"`
}

func (this *FrontController) Get() {
    uname := this.GetString("uname")
    if len(uname) == 0 {
        this.Redirect("/", 302)
        return
    }

    this.TplName = "front.html"
    this.Data["IsWebSocket"] = true
    this.Data["UserName"] = uname
}

func (this *FrontController) Friends() {
    uname := this.GetString("uname")
    friends := GetFriends(uname)

    var stat []UnreadStat

    for _, friend := range friends {
        line_cnt := CountUnread(uname, friend)
        stat = append(stat, UnreadStat{Friend:friend, Unread:line_cnt})
    }
    this.Data["json"] = &stat
    this.ServeJSON()
}

func (this *FrontController) AddFriend() {
    this.TplName = "front.html"
    friend_name := this.GetString("friend_name")
    uname := this.GetString("uname")

    if 0 == len(friend_name) || 0 == len(uname) {
        logs.Error("invalid parameter:" + strings.Join([]string{friend_name, uname}, ","))
        return
    }

    if !IsRegistered(friend_name) {
        logs.Error("Unregistered user")
        return
    }

    AddNewFriend(friend_name, uname)
    AddNewFriend(uname, friend_name)
} 

func (this *FrontController) RemoveFriend() {
    this.TplName = "front.html"
    friend_name := this.GetString("friend_name")
    uname := this.GetString("uname")

    if 0 == len(friend_name) || 0 == len(uname) {
        logs.Error("invalid parameter:" + strings.Join([]string{friend_name, uname}, ","))
        return
    }

    RemoveOldFriend(friend_name, uname)
    RemoveOldFriend(uname, friend_name)
}


