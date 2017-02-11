package main

import (
    "github.com/astaxie/beego"
    "github.com/beego/i18n"

    "lls/controllers"
)

func main() {

    beego.Router("/", &controllers.AppController{})
    beego.Router("/login", &controllers.AppController{}, "get:Login")
    beego.Router("/register", &controllers.AppController{}, "post:Register")
    beego.Router("/front", &controllers.FrontController{})
    beego.Router("/front/friends", &controllers.FrontController{}, "get:Friends")
    beego.Router("/front/addfriend", &controllers.FrontController{}, "post:AddFriend")
    beego.Router("/front/remove_friend", &controllers.FrontController{}, "post:RemoveFriend")
    beego.Router("/chatroom", &controllers.ChatroomController{})
    beego.Router("/chatroom/talk", &controllers.ChatroomController{}, "post:Talk")
    beego.Router("/chatroom/listen", &controllers.ChatroomController{}, "post:Listen")
    beego.Router("/chatroom/history", &controllers.ChatroomController{}, "get:History")

    beego.AddFuncMap("i18n", i18n.Tr)
    beego.Run()
}

