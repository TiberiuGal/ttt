/**
 * Created by tgal on 1/10/2017.
 */
var game = {};
var offset = {};
var cnf = {};
$(document).ready(function () {
    listGames();
    show('#gamelist');
});
function show(id) {
    $('.tab1').hide();
    $(id).show();
}


function listGames() {
    $.get("/listGames", function (data) {
        console.log(data);
        var ul = $('#gamelist ul').empty();
        for (el in data.G) {
            console.log(el);
            el = data.G[el];
            ul.append('<li data-id="' + el.Id + '">' + el.Id + ',' + el.Status + '</li>');
        }
        $('#gamelist li').click(function (event) {
            var gameId = $(this).attr('data-id');
            goPlay(gameId);
            show('#gamearea');
        });
    });
}


function goPlay(id) {
    var socket = null;
    var msgBox = $("#chatbox textarea");
    var messages = $('#messages');
    var myFill = 0;
    $('#chatbox').submit(function () {
        if (!msgBox.val()) return false;
        if (!socket) {
            console.log("err: there's no socket open");
            return false;
        }
        socket.send(JSON.stringify({"Message": msgBox.val()}));
        msgBox.val("");
        return false;
    });
    if (!window["WebSocket"]) {
        console.log("your browsers sucks");
        return false;
    }
    socket = new WebSocket("ws://"+jsg_Host+"/wg/" + id);
    socket.onclose = function () {
        console.log("youre out");
    };
    socket.onmessage = function (e) {
        var msg = JSON.parse(e.data);
        console.log(msg);
        if (msg.Txt.length > 0) {
            try {
                act = JSON.parse(msg.Txt);
                console.log(act);
                if (act.Result === "ok") {
                    switch (act.Action) {
                        case "join":
                            console.log("joined successfully as " + act.Fill);
                            myFill = act.Fill;
                            if(myFill != 0) {
                                playGame(socket, id, myFill);
                            }
                            break;
                        case "start":
                            console.log("game started, it's " + act.Fill +"'s turn");
                            game.cp = act.Fill;
                            break;
                        case "move":
                            game.cp = act.Fill;
                            console.log("move made, it's " + act.Fill +"'s turn");
                        default :
                            break;
                    }

                    load(id);
                    return;
                }

            } catch (ex) {
                console.log(ex);
            }

        }

        messages.append($('<li>').append(
            $('<img>').css({
                width: 50,
                verticalAlign: "middle"
            }).attr("src", msg.AvatarURL),
            $("<strong>").text(msg.Name + ": "),
            $("<span>").text(msg.Message)
            )
        );
    };



}
function playGame(socket, gameId, myFill) {
    cnf = {
        w: 20,
        h: 20,
        o: 5
    };
    offset = {
        x: 0,
        y: 0,
        w: 0,
        h: 0
    };
    game = {
        cp: 0,
        canvas: null,
        my : myFill
    };


    game.canvas = document.getElementById('canvas');
    load(gameId);

    $(canvas).click(function (e) {
        if (game.cp != game.my) {
            alert('it is not your turn yet');
            return;
        }
        var pos = getMousePos(this, e);

        socket.send(JSON.stringify({"Message": pos, "Txt": "df"}));

    });


}

function load(id) {
    $.ajax("//localhost:8021/game/" + id, {
        "dataType": "json", "success": function (data) {
            if (data.Status == 1) {
                drawPlayers(canvas, data.Players);
                draw(game.canvas, data.GameMap);

            } else if (data.Status == 2) {
                if (game.my == 0) {
                    drawVictorForOthers(game.canvas, data.Winner);
                } else {
                    drawVictory(game.canvas, data.Winner == game.my);
                }
            }
        }
    });
}

function drawPlayers(canvas, data) {

    ctx = canvas.getContext('2d');
    ctx.clearRect(0, 30, canvas.offsetWidth, 30);
    ctx.fillStyle = "#333";
    var p1 = data[0];
    var p2 = data[1];
    var img1 = new Image();

    img1.onload = function() {
        ctx.drawImage(img1, 20, 0, 30, 30);
    };
    img1.src = p1.Avatar;
    var img2 = new Image();
    img2.onload = function() {
        ctx.drawImage(img2, 300, 0, 30, 30);
    };
    img2.src = p2.Avatar;

    ctx.fillText("Player 1 "+ p1.Name + " as red ", 60, 10);

    ctx.fillText("vs ", 280, 10);
    ctx.fillText("Player 2 " + p2.Name + " as green", 340, 10);

}

function drawVictorForOthers(canvas, data) {
    ctx = canvas.getContext('2d');
    //ctx.clearRect(0, 0, canvas.offsetWidth, 30);
    ctx.font = "24px Georgia";

    ctx.fillStyle = "green";
    ctx.fillText("Player "+data+" has won!!!", 20, 20);
}

function drawVictory(canvas, data) {
    ctx = canvas.getContext('2d');
    //ctx.clearRect(0, 0, canvas.offsetWidth, canvas.offsetHeight);
    ctx.font = "24px Georgia";

    if (data) {
        ctx.fillStyle = "blue";
        ctx.fillText("You won!!!", 20, 20);
    } else  {
        ctx.fillStyle = "red";
        ctx.fillText("You lost!!!", 20, 20);
    }
}

function draw(canvas, data) {

    ctx = canvas.getContext('2d');
    ctx.clearRect(0, 30, canvas.offsetWidth, canvas.offsetHeight);
    var mX = 0, MX = 0, mY = 0, MY = 0;
    for (i in data) {
        x = data[i].X;
        y = data[i].Y;
        if (mX > x) {
            mX = x
        }
        if (MX < x) {
            MX = x
        }
        if (mY > y) {
            mY = y
        }
        if (MY < y) {
            MY = y
        }
    }
    ctx.strokeStyle = "gray";
    ctx.lineWidth = "1";

    offset.x = mX - cnf.o;
    offset.y = mY - cnf.o;
    offset.w = (Math.abs(mX) + cnf.o) * cnf.w;
    offset.h = (Math.abs(mY) + cnf.o) * cnf.h;
    ctx.translate(offset.w, offset.h);
    for (var j = mY - 3; j <= MY + 3; j++) {
        for (var i = mX - 3; i <= MX + 3; i++) {
            ctx.strokeRect(i * cnf.w, j * cnf.h, cnf.w, cnf.h);
            idx = "" + i + "." + j;
            if (data[idx] != undefined) {
                switch (data[idx].Content) {
                    case 1:
                        ctx.fillStyle = "red";
                        break;
                    case 2:
                        ctx.fillStyle = "green";
                        break;
                }

                ctx.fillRect(i * cnf.w, j * cnf.h, cnf.w, cnf.h);
                ctx.fillStyle = "blue";
                ctx.fillText(i + " " + j, i * cnf.w, (j * cnf.h) + 11);

            }
        }
    }
    ctx.translate(-offset.w, -offset.h);


}


function getMousePos(canvas, evt) {
    var rect = canvas.getBoundingClientRect();
    return {
        x: Math.floor((evt.clientX - rect.left) / cnf.w) + offset.x,
        y: Math.floor((evt.clientY - rect.top) / cnf.h) + offset.y
    };
}