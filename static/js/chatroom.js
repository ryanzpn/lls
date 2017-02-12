
var chat_history = function () {
    $.getJSON("/chatroom/history?uname=" + $('#uname').text() + "&recipient=" + $('#recipient').text(), function (data) {
        if (data == null) {
            return;
        }
        $("#chatbox").empty();

        $.each(data, function (i, stat) {
            $("#chatbox").append("<li>" + (new Date(parseInt(stat.timestamp))).toString().substring(4, 24) + " <b>" + stat.speaker + "</b>: " + stat.line + "</li>");
        });
    });
}

var talk = function() {
    var uname = $('#uname').text();
    var recipient = $('#recipient').text();
    var content = $('#sendbox').val();

    if (0 == content.length) {
        return;
    }

    $("#sendbox").val("");

    $.post("/chatroom/talk?uname=" + $('#uname').text() + "&recipient=" + $('#recipient').text(), 
           {uname:uname, recipient:recipient, content:content}, 
           function (data) {
               if (data == null) {
                   return;
               }
               var libody = (new Date(parseInt(data['timestamp']))).toString().substring(4, 24) + " <b>" + uname + "</b>: " + content;
               $("#chatbox").append("<li>" + libody + "</li>");

    }, "json");
}

var listen = function() {
    var uname = $('#uname').text();
    var recipient = $('#recipient').text();

    $.post("/chatroom/listen?uname=" + $('#uname').text() + "&recipient=" + $('#recipient').text(), 
           {}, 
           function (data) {
               if (data == null) {
                   return;
               }

               $.each(data, function (i, stat) {
                   $("#chatbox").append("<li>" + (new Date(parseInt(stat.timestamp))).toString().substring(4, 24) + " <b>" + stat.speaker + "</b>: " + stat.line + "</li>");
               });
    }, "json");
}

$(document).ready(function () {
    chat_history();

    $('#sendbtn').click(function () {
        talk();
    });

    $('#updatebtn').click(function () {
        listen();
    });

});

