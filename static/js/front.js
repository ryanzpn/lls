


var fetch = function() {
    $.getJSON('/front/friends?uname=' + $('#uname').text(), function(data) {
        if (null == data) {
            $('#contacts').empty()
            return;
        }
        $('#contacts').empty()

        $.each(data, function (i, stat) {
            $("#contacts").append("<li><a href='/chatroom?uname=" + $('#uname').text() + "&recipient=" + stat.friend + "'>" + stat.friend + " (" + stat.unread + ")</a></li>");
        });
    });
}

var addfriend = function() {
    $.post('/front/addfriend?uname=' + $('#uname').text() + '&friend_name=' + $('#friend_name').val(), 
           {}, 
           function(data) {
        if (data == null) {
            return;
        }
    });
}

var removefriend = function() {
    $.post('/front/remove_friend?uname=' + $('#uname').text() + '&friend_name=' + $('#r_friend_name').val(), 
           {}, 
           function(data) {
        if (data == null) {
            return;
        }
    });
}

$(document).ready(function () {
    fetch();

    $('#updatebtn').click(function () {
        fetch();
    });

    $('#addfriendbtn').click(function() {
        addfriend();
    });

    $('#removefriendbtn').click(function() {
        removefriend();
    });
});
