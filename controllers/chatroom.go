package controllers

import (
    "os"
    "strings"
    "strconv"
    "time"
)

type ChatroomController struct {
    baseController
}

type Msg struct {
    Line string `json:"line"`
    Speaker string `json:"speaker"`
    TimeStamp int `json:"timestamp"`
}

func (this *ChatroomController) Get() {
    uname := this.GetString("uname")
    recipient := this.GetString("recipient")

    this.TplName = "chatroom.html"
    this.Data["IsNewWebSocket"] = true
    this.Data["UserName"] = uname 
    this.Data["Recipient"] = recipient 
}

func (this *ChatroomController) History() {
    uname := this.GetString("uname")
    recipient := this.GetString("recipient")

    if 0 == len(uname) || 0 == len(recipient) {
        this.Redirect("/", 302)
        return
    }

    unread_dir := strings.Join([]string{"status", uname, recipient, "unread.txt"}, "/")
    unread := ReadFileAll(unread_dir)

    read_dir := strings.Join([]string{"status", uname, recipient, "read.txt"}, "/")
    WriteMsg(unread, read_dir, true)
    os.Remove(unread_dir)

    read := ReadFileAll(read_dir)

    history := PackInMsg(strings.Split(read, "\n"), uname, recipient)
    this.Data["json"] = &history
    this.ServeJSON()
}

func (this *ChatroomController) Talk() {
    uname := this.GetString("uname")
    recipient := this.GetString("recipient")
    content := this.GetString("content")

    if 0 == len(uname) || 0 == len(recipient) {
        this.Redirect("/", 302)
        return
    }

    // uname side
    read_dir := strings.Join([]string{"status", uname, recipient, "read.txt"}, "/")
    now := time.Now().Unix()
    timestamp := strconv.FormatInt(now, 10)

    read := timestamp + "\t0\t" + content + "\n"
    WriteMsg(read, read_dir, true)

    // recipient side
    // TODO
    unread_dir := strings.Join([]string{"status", recipient, uname, "unread.txt"}, "/")
    unread := timestamp + "\t1\t" + content + "\n"
    WriteMsg(unread, unread_dir, true)

    // return 
    this.Data["json"] = &Msg{Line:content, Speaker:uname, TimeStamp:int(now)}
    this.ServeJSON()
}

func (this *ChatroomController) Listen() {
    uname := this.GetString("uname")
    recipient := this.GetString("recipient")

    unread_dir := strings.Join([]string{"status", uname, recipient, "unread.txt"}, "/")
    unread := ReadFileAll(unread_dir)

    read_dir := strings.Join([]string{"status", uname, recipient, "read.txt"}, "/")
    WriteMsg(unread, read_dir, true)

    os.Remove(unread_dir)
    msgs := PackInMsg(strings.Split(unread, "\n"), uname, recipient)

    this.Data["json"] = &msgs
    this.ServeJSON()
}


