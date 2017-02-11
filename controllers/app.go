package controllers

import (
    "strings"
    "fmt"
    "time"

    "database/sql"
    "github.com/astaxie/beego"
    "github.com/beego/i18n"
    "github.com/astaxie/beego/logs"
    _ "github.com/go-sql-driver/mysql"
)

var langTypes []string // Languages that are supported.

func init() {
    log := logs.NewLogger(10000)
    log.SetLogger("console")
    // Initialize language type list.
    langTypes = strings.Split(beego.AppConfig.String("lang_types"), "|")

    // Load locale files according to language types.
    for _, lang := range langTypes {
        beego.Trace("Loading language: " + lang)
        if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
            beego.Error("Fail to set message file:", err)
            return
        }
    }
}

type baseController struct {
    beego.Controller 
    i18n.Locale 
}

func (this *baseController) Prepare() {
    // Set template level language option.
    this.Data["Lang"] = "en-US"
}

type AppController struct {
    baseController 
}

func (this *AppController) Get() {
    this.TplName = "welcome.html"
}

func (this *AppController) Login() {
    uname := this.GetString("uname")
    passwd := this.GetString("passwd")

    if 0 == len(uname) || !IsValidUserPasswd(uname, passwd) {
        this.Redirect("/", 302)
        return
    }
    this.Redirect("/front?uname="+uname, 302)
}

var (
    host = "127.0.0.1"
    user = "lls"
    password = "lls"
    port = "3306"
    db = "lls"
)

func (this *AppController) Register() {
    this.TplName = "welcome.html"  
    uname := this.GetString("uname")
    passwd := this.GetString("passwd")
    confirm := this.GetString("confirm")

    if 0 == len(uname) || 0 == len(passwd) || passwd != confirm {
        logs.Error("invalid parameter:" + strings.Join([]string{uname, passwd, confirm}, ","))
        this.Redirect("/", 302)
        return
    }

    if IsValidUserPasswd(uname, passwd) {
        this.Redirect("/", 302)
        return
    }

    url := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8"
    db, err := sql.Open("mysql", url)
    if err != nil {
        logs.Error(err.Error())
        panic(err.Error())
    }
    defer db.Close()

    sqlStmt, err := db.Prepare("INSERT INTO user_info VALUES (?, ?, ?);")
    if nil != err {
        logs.Error(err.Error())
        this.Redirect("/", 302)
        panic(err.Error())
    }
    defer sqlStmt.Close()

    _, err = sqlStmt.Exec(uname, passwd, time.Now().Unix())
    if err != nil {
        logs.Error(err.Error())
        this.Redirect("/", 302)
        panic(err.Error())
    }
}

func IsValidUserPasswd(uname, passwd string) bool {
    url := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8"
    db, err := sql.Open("mysql", url)
    if err != nil {
        fmt.Println(err)
        return false
    }
    defer db.Close()

    sqlStmt := "SELECT passwd FROM user_info WHERE uname = '" + uname + "';"
    rows, err := db.Query(sqlStmt)
    if err != nil {
        fmt.Println(err)
        return false
    }

    rowCnt := 0
    var pswd string

    for rows.Next() {
        err := rows.Scan(&pswd)
        if err != nil {
            fmt.Println(err)
            return false
        }
        rowCnt += 1
    }
    if 1 != rowCnt {
        fmt.Println("rowCnt error")
        return false
    }
    if passwd != pswd {
        fmt.Println("passwd error")
        return false
    }

    return true
}

func IsRegistered(uname string) bool {
    url := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?charset=utf8"
    db, err := sql.Open("mysql", url)
    if err != nil {
        fmt.Println(err)
        return false
    }
    defer db.Close()

    sqlStmt := "SELECT passwd FROM user_info WHERE uname = '" + uname + "';"
    rows, err := db.Query(sqlStmt)
    if err != nil {
        fmt.Println(err)
        return false
    }

    rowCnt := 0
    for rows.Next() {
        rowCnt += 1
    }
    if 0 == rowCnt {
        return false
    }
    return true;
} 
