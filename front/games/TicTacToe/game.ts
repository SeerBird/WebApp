window.onload = function () {

    var conn: WebSocket;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws/TicTacToe");
        conn.onclose = function (evt) {

        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (let m = 0; m < messages.length; m++) {
                var data=JSON.parse(messages[m])
                if (data.Tag == "update") {
                    for(let i=0;i<3;i++){
                        for(let j=0;j<3;j++){
                            grid[i][j]=data.Msg[i][j]
                        }
                    }
                }
            }
            redraw()
        };
    } else {
        //scream
        return
    }
    const canvas = <HTMLCanvasElement>document.querySelector(".myCanvas");
    const grid = [[-1, -1, -1], [-1, -1, -1], [-1, -1, -1]]
    if (canvas == null) {
        return
        //handle pls?
    }

    canvas.addEventListener("click", onClick)
    const width = (canvas.width = window.innerWidth);
    const height = (canvas.height = window.innerHeight);
    const ctx = canvas.getContext("2d");
    if (ctx == null) {
        return
        //handle pls?
    }
    const side = Math.min(width, height) / 3
    const ox = width / 2 - side - side / 2
    const oy = height / 2 - side - side / 2
    function redraw() {
        //#region grid
        ctx.fillStyle = "white"
        ctx.fillRect(0, 0, width, height);
        ctx.strokeStyle = "black"
        ctx.beginPath()
        ctx.moveTo(ox, oy)
        ctx.lineTo(ox, oy + 3 * side)
        ctx.lineTo(ox + 3 * side, oy + 3 * side)
        ctx.lineTo(ox + 3 * side, oy)
        ctx.lineTo(ox, oy)

        ctx.moveTo(ox + side, oy)
        ctx.lineTo(ox + side, oy + 3 * side)
        ctx.moveTo(ox + 2 * side, oy)
        ctx.lineTo(ox + 2 * side, oy + 3 * side)
        ctx.moveTo(ox, oy + side)
        ctx.lineTo(ox + 3 * side, oy + side)
        ctx.moveTo(ox, oy + 2 * side)
        ctx.lineTo(ox + 3 * side, oy + 2 * side)
        ctx.stroke()
        //#endregion
        for (let i = 0; i < 3; i++) {
            for (let j = 0; j < 3; j++) {
                if (grid[i][j] == 1) {
                    ctx.beginPath()
                    ctx.arc(ox + side * i + side / 2, oy + side * j + side / 2, side / 2, 0, 7)
                    ctx.stroke()
                } else if (grid[i][j] == 0) {
                    ctx.beginPath()
                    ctx.moveTo(ox + side * i, oy + side * j)
                    ctx.lineTo(ox + side * (i + 1), oy + side * (j + 1))
                    ctx.moveTo(ox + side * (i + 1), oy + side * j)
                    ctx.lineTo(ox + side * i, oy + side * (j + 1))
                    ctx.stroke() //this api really is the worst huh 
                }
            }
        }
    }
    function onClick(this, e: MouseEvent) {
        if (this instanceof HTMLCanvasElement) {
            const i = Math.floor((e.x - ox) / side)
            const j = Math.floor((e.y - oy) / side)
            if (i > -1 && i < 3 && j > -1 && j < 3) {
                conn.send(JSON.stringify({ "type": "input", "value": { "i": i, "j": j } }))
            }
        }
    }
    ctx.fillStyle = "rgb(0 255 0)";
    ctx.fillRect(0, 0, width, height);
};
