const size: number = 5
window.onload = function () {
    //region init game model
    var grid: string[][] = [];
    for (var i: number = 0; i < size; i++) {
        grid[i] = []
        for (var j: number = 0; j < size; j++) {
            grid[i][j] = ""
        }
    }
    var myTurn: boolean = false
    var word: coordinate[] = []
    var mousepos: coordinate = { i: -1, j: -1 }
    var playerList = []
    //endregion
    //region init canvas and size constants
    const canvas = <HTMLCanvasElement>document.querySelector(".myCanvas");
    if (canvas == null) {
        return
        //handle pls?
    }
    canvas.addEventListener("click", onClick)
    canvas.addEventListener("mousemove", onMove)
    const width = (canvas.width = window.innerWidth);
    const height = (canvas.height = window.innerHeight);
    const side = Math.min(width, height) / size
    const ox = width / 2 - side * (size / 2)
    const oy = height / 2 - side * (size / 2)
    const ctx = canvas.getContext("2d");
    const font = side * 2 / 3
    if (ctx == null) {
        return
        //handle pls?
    }
    //endregion
    //region init scoreboard
    const scoreboard = document.getElementById("scoreboard");
    function addSBItem(text:string,active:boolean){
        var entry = document.createElement('li');
        entry.appendChild(document.createTextNode(text));
        if(active){
            entry.style.color="green"
        }
        scoreboard.appendChild(entry);
    }
    function clearSB(){
        scoreboard.innerHTML=""
    }
    //endregion
    //region input
    function onClick(this, e: MouseEvent) {
        //region do nothing if it's not my turn
        if (!myTurn) {
            return
        }
        //endregion
        //region clear word selection if click is outside the canvas
        if (!(this instanceof HTMLCanvasElement)) {
            word = []
            return
        }
        const i = Math.floor((e.x - ox) / side)
        const j = Math.floor((e.y - oy) / side)
        if (i < 0 || i > size - 1 || j < 0 || j > size - 1) {
            word = []
            return //can this even happen? we're outside the grid. whatever.
        }
        //endregion
        //region start word
        if (word.length == 0) {
            word[0] = { i: i, j: j }
        }
        //endregion
        //region clear word selection and return if selection is one letter
        else if ((word.length == 1)){
            word = []
        }
        //endregion
        //region end word
        else {
            conn.send(JSON.stringify({ "type": "input", "value": word})) //validate this on server
            word = []
        }
        //endregion
        redraw()
    }
    function onMove(this, e: MouseEvent) {
        //region redraw
        redraw()
        ctx.strokeStyle="green"
        drawCircle(e.x, e.y, 5)
        //endregion
        //region validate input
        mousepos = { i: e.x, j: e.y }
        if (!(this instanceof HTMLCanvasElement)) {
            return
        }
        if (word.length == 0) {
            return
        }
        const i = Math.floor((e.x - ox) / side)
        const j = Math.floor((e.y - oy) / side)
        console.log(" ")
        console.log(i+", "+j)
        if (i < 0 || i > size - 1 || j < 0 || j > size - 1) {
            console.log("outside")
            return //can this even happen? we're outside the grid. whatever.
        }
        const lastLetter = word[word.length - 1]
        if (i == lastLetter.i && j == lastLetter.j) { // we're in the last letter
            console.log("in the last letter")
            return
        }
        if (Math.abs(i - lastLetter.i) > 1 || Math.abs(j - lastLetter.j) > 1) {
            console.log("not in a square neighbouring the previous letter")
            return // we're not in a square neighbouring the previous letter
        }
        const centerx = ox + i * side + side / 2
        const centery = oy + j * side + side / 2
        ctx.strokeStyle="red"
        drawCircle(centerx, centery, 10)
        if (Math.abs(e.x - centerx) + Math.abs(e.y - centery) > side / 2) {
            console.log("Not in the diamond")
            return // only trigger in the diamond
        }
        //endregion
        //region erase last letter or append hovered letter
        if (word.length > 1) { // we can erase tha last letter
            if (i == word[word.length - 2].i && j == word[word.length - 2].j) {
                // we're in the diamond of the letter before last
                console.log("erasing")
                word.pop()
                return
            }
        }
        for(let m=0;m<word.length;m++){
            if (i == word[m].i && j == word[m].j){
                return //this coordletter is already in the word
            }
        }
        word[word.length] = { i: i, j: j } // append hovered letter
        //endregion        
    }
    //endregion
    //region drawing
    function redraw() {
        //region grid
        ctx.fillStyle = "white"
        ctx.fillRect(0, 0, width, height);
        ctx.strokeStyle = "black"
        for (var i = 0; i < size + 1; i++) {
            drawLine(ox + i * side, oy, ox + i * side, oy + size * side);
        }
        for (var i = 0; i < size + 1; i++) {
            drawLine(ox, oy + side * i, ox + size * side, oy + side * i);
        }
        //endregion
        //region letters
        ctx.fillStyle = "red"
        ctx.font = "bold " + font + "px serif"
        for (var i = 0; i < size; i++) {
            for (var j = 0; j < size; j++) {
                var letter = grid[i][j];
                var dims = ctx.measureText(letter);
                ctx.fillText(grid[i][j], ox + side * i + side / 2 - dims.width / 2, oy + side * j + side / 2 + font / 3)
            }
        }
        //endregion
        //region selection
        if (word.length > 0) {
            //region start letter circle
            ctx.beginPath()
            var lastx = ox + side * word[0].i + side / 2
            var lasty = oy + side * word[0].j + side / 2
            ctx.arc(lastx, lasty, side / 2, 0, 6.3)
            ctx.strokeStyle="green"
            ctx.stroke()
            //endregion
            //region fixed word path
            var nextx
            var nexty
            for(var i = 1;i<word.length;i++){
                nextx = ox + side * word[i].i + side / 2
                nexty = oy + side * word[i].j + side / 2
                drawLine(lastx,lasty,nextx,nexty)
                lastx=nextx
                lasty=nexty
            }
            //endregion
            //region hanging word path link
            drawLine(lastx, lasty, mousepos.i, mousepos.j)
            //endregion
        }
        //endregion
        //region debug
        //ctx.fillText(String(myTurn), 100,100)
        //endregion
    }
    function drawLine(x0: number, y0: number, x1: number, y1: number) {
        ctx.beginPath()
        ctx.moveTo(x0, y0)
        ctx.lineTo(x1, y1)
        ctx.stroke()
    }
    function drawCircle(x: number, y: number, r: number) {
        ctx.beginPath()
        ctx.arc(x, y, r, 0, 6.3)
        ctx.stroke()
    }
    //endregion
    //region connection
    var conn: WebSocket;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws/Wordle");
        conn.onclose = function (evt) {

        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (let m = 0; m < messages.length; m++) {
                var data = JSON.parse(messages[m])
                switch (data.Tag) {
                    case "update":// grid, turn, playerList(ordered (name,score))
                        var msg = data.Msg;
                        grid = msg.Grid;
                        myTurn = msg.ClientOrder == msg.Turn;
                        playerList = msg.PlayerList;
                        clearSB()
                        for(let i=0;i<playerList.length;i++){
                            var name = playerList[i].Name
                            if (playerList[i].Name == String(msg.ClientOrder)){
                                name="You"
                            }
                            addSBItem(name + ": " + playerList[i].Score, i == msg.Turn)
                        }
                }
            }
            redraw()
        };
    } else {
        //scream
        return
    }
    //endregion
    ctx.fillStyle = "rgb(0 255 0)";
    ctx.fillRect(0, 0, width, height);
};

interface coordinate {
    i: number;
    j: number;
}
