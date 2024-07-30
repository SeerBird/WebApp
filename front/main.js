window.onload = function () {
    var conn;
    /*
    {
        var msg = document.getElementById("msg");
        var log = document.getElementById("log");

        function appendLog(item) {
            var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
            log.appendChild(item);
            if (doScroll) {
                log.scrollTop = log.scrollHeight - log.clientHeight;
            }
        }

        document.getElementById("form").onsubmit = function () {
            if (!conn) {
                return false;
            }
            if (!msg.value) {
                return false;
            }
            var message = { "type": "input", "value": msg.value }
            conn.send(JSON.stringify(message));
            msg.value = "";
            return false;
        };
        if (window["WebSocket"]) {
            conn = new WebSocket("ws://" + document.location.host + "/ws/bababoi");
            conn.onclose = function (evt) {
                var item = document.createElement("div");
                item.innerHTML = "<b>Connection closed.</b>";
                appendLog(item);
            };
            conn.onmessage = function (evt) {
                var messages = evt.data.split('\n');
                for (var i = 0; i < messages.length; i++) {
                    var item = document.createElement("div");
                    item.innerText = messages[i];
                    appendLog(item);
                }
            };
        } else {
            var item = document.createElement("div");
            item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
            appendLog(item);
        }
    }
    */
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws/bababoi");
        conn.onclose = function (evt) {
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var m = 0; m < messages.length; m++) {
                var data = JSON.parse(messages[m]);
                if (data.Tag == "update") {
                    for (var i = 0; i < 3; i++) {
                        for (var j = 0; j < 3; j++) {
                            grid[i][j] = data.Msg[i][j];
                        }
                    }
                }
            }
            redraw();
        };
    }
    else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        return;
    }
    var canvas = document.querySelector(".myCanvas");
    var grid = [[-1, -1, -1], [-1, -1, -1], [-1, -1, -1]];
    if (canvas == null) {
        return;
        //handle pls?
    }
    canvas.addEventListener("click", onClick);
    var width = (canvas.width = window.innerWidth);
    var height = (canvas.height = window.innerHeight);
    var ctx = canvas.getContext("2d");
    if (ctx == null) {
        return;
        //handle pls?
    }
    var side = Math.min(width, height) / 3;
    var ox = width / 2 - side - side / 2;
    var oy = height / 2 - side - side / 2;
    function redraw() {
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, width, height);
        ctx.strokeStyle = "black";
        ctx.beginPath();
        ctx.moveTo(ox, oy);
        ctx.lineTo(ox, oy + 3 * side);
        ctx.lineTo(ox + 3 * side, oy + 3 * side);
        ctx.lineTo(ox + 3 * side, oy);
        ctx.lineTo(ox, oy);
        ctx.moveTo(ox + side, oy);
        ctx.lineTo(ox + side, oy + 3 * side);
        ctx.moveTo(ox + 2 * side, oy);
        ctx.lineTo(ox + 2 * side, oy + 3 * side);
        ctx.moveTo(ox, oy + side);
        ctx.lineTo(ox + 3 * side, oy + side);
        ctx.moveTo(ox, oy + 2 * side);
        ctx.lineTo(ox + 3 * side, oy + 2 * side);
        ctx.stroke();
        for (var i = 0; i < 3; i++) {
            for (var j = 0; j < 3; j++) {
                if (grid[i][j] == 1) {
                    ctx.beginPath();
                    ctx.moveTo(ox + side * i, oy + side * j);
                    ctx.lineTo(ox + side * (i + 1), oy + side * (j + 1));
                    ctx.moveTo(ox + side * (i + 1), oy + side * j);
                    ctx.lineTo(ox + side * i, oy + side * (j + 1));
                    ctx.stroke();
                }
                else if (grid[i][j] == 0) {
                    ctx.beginPath();
                    ctx.arc(ox + side * i + side / 2, oy + side * j + side / 2, side / 2, 0, 7);
                    ctx.stroke();
                }
            }
        }
    }
    function onClick(e) {
        if (this instanceof HTMLCanvasElement) {
            var i = Math.floor((e.x - ox) / side);
            var j = Math.floor((e.y - oy) / side);
            if (i > -1 && i < 3 && j > -1 && j < 3) {
                conn.send(JSON.stringify({ "type": "input", "value": { "i": i, "j": j } }));
            }
        }
    }
    ctx.fillStyle = "rgb(0 255 0)";
    ctx.fillRect(0, 0, width, height);
};
//# sourceMappingURL=main.js.map