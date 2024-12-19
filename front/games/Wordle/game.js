var size = 5;
window.onload = function () {
    var conn;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws/TicTacToe");
        conn.onclose = function (evt) {
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var m = 0; m < messages.length; m++) {
                var data = JSON.parse(messages[m]);
                switch (data.Tag) {
                    case "update": // grid, turn, playerList(ordered (name,score))
                        grid = null;
                }
            }
            redraw();
        };
    }
    else {
        //scream
        return;
    }
    var canvas = document.querySelector(".myCanvas");
    if (canvas == null) {
        return;
        //handle pls?
    }
    canvas.addEventListener("click", onClick);
    canvas.addEventListener("mousemove", onMove);
    var width = (canvas.width = window.innerWidth);
    var height = (canvas.height = window.innerHeight);
    var side = Math.min(width, height) / size;
    var ox = width / 2 - side * (size / 2) - side / 2;
    var oy = height / 2 - side * (size / 2) - side / 2;
    var ctx = canvas.getContext("2d");
    if (ctx == null) {
        return;
        //handle pls?
    }
    var grid = [];
    for (var i = 0; i < size; i++) {
        grid[i] = [];
        for (var j = 0; j < size; j++) {
            grid[i][j] = "";
        }
    }
    var myTurn = false;
    var word = [];
    function onClick(e) {
        //region do nothing if it's not my turn
        if (!myTurn) {
            return;
        }
        //endregion
        //region clear word selection if click is outside the canvas or RMB
        if (e.button == 2) {
            word = [];
            return;
        }
        if (!(this instanceof HTMLCanvasElement)) {
            word = [];
            return;
        }
        //endregion
        //region start word
        if (word.length == 0) {
            var i_1 = Math.floor((e.x - ox) / side);
            var j_1 = Math.floor((e.y - oy) / side);
            if (i_1 > -1 && i_1 < size && j_1 > -1 && j_1 < size) {
                word[0] = { i: i_1, j: j_1 };
            }
        }
        //endregion
        //region end word
        else {
            conn.send(JSON.stringify(word)); //validate this on server
            word = [];
            myTurn = false; //maybe wait for update from server anyways?
        }
        //endregion
    }
    function onMove(e) {
        //region validate input
        if (!(this instanceof HTMLCanvasElement)) {
            return;
        }
        if (word.length == 0) {
            return;
        }
        var i = Math.floor((e.x - ox) / side);
        var j = Math.floor((e.y - oy) / side);
        if (i < 0 || i > size - 1 || j < 0 || j > size - 1) {
            return; //can this even happen? we're outside the grid. whatever.
        }
        var lastLetter = word[word.length - 1];
        if (i == lastLetter.i && j == lastLetter.j) { // we're in the last letter
            return;
        }
        if (Math.abs(i - lastLetter.i) > 1 || Math.abs(j - lastLetter.j) > 1) {
            return; // we're not in a square neighbouring the previous letter
        }
        var centerx = ox + i * side + side / 2;
        var centery = oy + i * side + side / 2;
        if (Math.abs(e.x - centerx) + Math.abs(e.y - centery) > side / 2) {
            return; // only trigger in the diamond
        }
        //endregion
        //region erase last letter or append hovered letter
        if (word.length > 1) { // we can erase tha last letter
            if (i == word[word.length - 2].i && j == word[word.length - 2].j) {
                // we're in the diamond of the letter before last
                word.pop();
                return;
            }
        }
        word[word.length] = { i: i, j: j }; // append hovered letter
        //endregion
    }
    function redraw() {
        //region grid
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, width, height);
        ctx.strokeStyle = "black";
        for (var i = 0; i < size + 1; i++) {
            drawLine(ctx, ox + i * side, 0, ox + i * side, height);
        }
        for (var i = 0; i < size + 1; i++) {
            drawLine(ctx, 0, oy + side * i, width, oy + side * i);
        }
        //endregion
        //region letters
        for (var i = 0; i < size + 1; i++) {
            for (var j = 0; j < size; j++) {
                ctx.fillText(grid[i][j], ox + side * i, oy + side * j);
            }
        }
        //endregion
    }
    ctx.fillStyle = "rgb(0 255 0)";
    ctx.fillRect(0, 0, width, height);
};
function drawLine(ctx, x0, y0, x1, y1) {
    ctx.beginPath();
    ctx.moveTo(x0, y0);
    ctx.lineTo(x1, y1);
    ctx.stroke();
}
//# sourceMappingURL=game.js.map