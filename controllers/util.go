package controllers

import (
    "os"
    "io"
    "bytes"
    "strings"
    "strconv"
    "io/ioutil"
    "path/filepath"
    "github.com/astaxie/beego/logs"
)

func RemoveOldFriend(candidate, uname string) {
    friends := GetFriends(uname)
    index, iof := IsOldFriend(candidate, friends) 
    if !iof {
        logs.Error("He/She is not a friend")
        return
    }
    friends = RemoveElement(friends, index)
    ReplaceFile(strings.Join(friends, ","), 
                strings.Join([]string{"status", uname, "friends.txt"}, "/"))
}

func RemoveElement(arr []string, index int) []string {
    arr[len(arr) - 1], arr[index] = arr[index], arr[len(arr) - 1]
    return arr[:len(arr) - 1]
}

func AddNewFriend(candidate, uname string) {
    friends := GetFriends(uname)
    if _, iof := IsOldFriend(candidate, friends); iof {
        logs.Error("He/She is an old friend")
        return
    }

    filedir := strings.Join([]string{"status", uname, "friends.txt"}, "/")
    content := candidate
    if 0 < len(friends) {
        content = candidate + "," + strings.Join(friends, ",")
    }
    WriteMsg(content, filedir, false)
}

func IsOldFriend(candidate string, friends []string) (int, bool) {
    for i, friend := range friends {
        if (candidate == friend) {
            return i, true
        }
    }
    return -1, false
}

func CountUnread(uname, friend string) int {
    filedir := strings.Join([]string{"status", uname, friend, "unread.txt"}, "/")
    unread := ReadFileAll(filedir)

    if 0 == len(unread) {
        return 0
    }
    return strings.Count(unread, "\n")
}

func GetFriends(uname string) []string {
    filedir := strings.Join([]string{"status", uname, "friends.txt"}, "/")
    friends := ReadFileAll(filedir)

    if 0 == len(friends) {
        return make([]string, 0, 1) 
    }
    return strings.Split(strings.TrimSpace(friends), ",")
}

func ReadFileAll(fname string) string {
    finfo, err := os.Stat(fname)
    if os.IsNotExist(err) || finfo.IsDir() {
        return ""
    }

    infile, err := os.Open(fname)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := infile.Close(); err != nil {
            panic(err)
        }
    }()

    var allfile bytes.Buffer
    buf := make([]byte, 1024)
    for {
        n, err := infile.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if 0 == n {
            break
        }
        allfile.Write(buf[:n])
    }
    return allfile.String()
}

func ParseMsgLine(line, uname, recipient string) Msg {
    seg := strings.Split(strings.TrimSpace(line), "\t")
    ts, err := strconv.Atoi(seg[0])
    if nil != err {
        panic(err)
    }
    speaker := uname
    if "1" == seg[1] {
        speaker = recipient
    }
    return Msg{Line:seg[2], Speaker:speaker, TimeStamp:ts}
}

func PackInMsg(lines []string, uname, recipient string) []Msg {
    var msgs []Msg
    for _, line := range lines {
        if 0 != len(line) {
            msgs = append(msgs, ParseMsgLine(line, uname, recipient))
        }
    }
    return msgs
}

func WriteMsg(content, fname string, is_append bool) {
    fdir := filepath.Dir(fname)
    finfo, err := os.Stat(fdir)
    if os.IsNotExist(err) || !finfo.IsDir() {
        os.MkdirAll(fdir, 0755)
    }
    
    mask := os.O_WRONLY|os.O_CREATE
    if is_append {
        mask |= os.O_APPEND
    }

    outfile, err := os.OpenFile(fname, mask, 0644)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
   
    _, err = outfile.WriteString(content)
    if err != nil {
        panic(err)
    }
}

func ReplaceFile(content, fname string) {
    if 0 == len(content) {
        os.Remove(fname)
        return
    }

    err := ioutil.WriteFile(fname, []byte(content), 0644)
    if err != nil {
        panic(err)
    }
}

func CreateIfNotExist(uname, recipient string) {
    friend_dir := strings.Join([]string{"status", uname, recipient}, "/")

    finfo, err := os.Stat(friend_dir)
    if err != nil || !finfo.IsDir() {
        os.Mkdir(friend_dir, 0755)
    }
}

func Max(a, b int) int {
    if a >= b {
        return a
    }
    return b
}

func Min(a, b int) int {
    if a <= b {
        return a
    }
    return b
}


