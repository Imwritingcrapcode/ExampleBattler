var bcanvas;
var backP = 3;
backSketch = function (p) {
    bg_color = p.color(BG);
    left = p.color(LEFTC);
    backP = p;

    p.preload = function () {
        myChar = new InterfaceImage(p, 0, 0, "", "myChar", "", 0, 0);
        oppChar = new InterfaceImage(p, 730, 0, "", "oppChar", "", 0, 0);
        if (myChar.name !== PlN && getResolution(PlNum)[0] !== 0) {
            myChar.open("/images/locked/" + PlN + "_left.png", PlN, getResolution(PlNum)[0], getResolution(PlNum)[1]);
        } else if (myChar.name !== PlNum) {
            myChar.open("/images/locked/Placeholder_left.png", PlNum, 350, 550);
        }
        if (oppChar.name !== ON && getResolution(ONum)[0] !== 0) {
            oppChar.open("/images/locked/" + ON + "_right.png", ON, getResolution(ONum)[0], getResolution(ONum)[1]);
        } else if (oppChar.name !== ON) {
            oppChar.open("/images/locked/Placeholder_right.png", ONum, 350, 550);
        }
    };
    p.setup = function () {
        bcanvas = p.createCanvas(1280, 550);
        bcanvas.position(0, 0);
        bcanvas.style("z-index", "-1");
        p.background(bg_color);
        p.noLoop();
    };
    p.draw = function () {
        let p = backP;
        let k;
        if (isLight(PlC)) {
            k = 0.07
        } else {
            k = 0.28
        }
        let pColour = p.lerpColor(color(left.toString()), PlC, k);
        if (isLight(OC)) {
            k = 0.07
        } else {
            k = 0.28
        }
        let oppColour = p.lerpColor(color(left.toString()), OC, k);
        setGradient(p, 0, 230, 550, 320, 5, bg_color, pColour);
        setGradient(p, 730, 230, 550, 320, 5, bg_color, oppColour);
        setGradient(p, 550, 230, 180, 320, 5, bg_color, left);
        myChar.display();
        oppChar.display();
        console.log("drew girls");
    };

};

function setMyChar(PlayerName, PlayerNum, PlayerColour) {
    PlN = PlayerName;
    PlNum = PlayerNum;
    PlC = PlayerColour;
}

function setOppChar(OppName, OppNum, OpponentColour) {
    ON = OppName;
    ONum = OppNum;
    OC = OpponentColour;
}


function setGradient(p, x, y, w, h, r, c1, c2) {
    p.noFill();
    p.strokeWeight(2);
    //circle at the beginning
    for (let i = y; (i < y + r) && (i < y + h); i++) {
        let inter = p.map(i, y, y + h, 0, 1);
        let c = p.lerpColor(c1, c2, inter);
        p.stroke(c);
        let top_x = x + r - p.sqrt(p.sq(r) - p.sq(i - y - r));
        let bot_x = x + w - r + p.sqrt(p.sq(r) - p.sq(i - y - r));
        p.line(top_x, i, bot_x, i);
    }

    for (let i = y + r; i < y + h - r; i++) {
        let inter = p.map(i, y, y + h, 0, 1);
        let c = p.lerpColor(c1, c2, inter);
        p.stroke(c);
        p.line(x, i, x + w, i);
    }
    //circle at the end
    for (let i = y + h; (i >= y + r) && (i >= y + h - r); i--) {
        let inter = p.map(i, y, y + h, 0, 1);
        let c = p.lerpColor(c1, c2, inter);
        p.stroke(c);
        let top_x = x + r - p.sqrt(p.sq(r) - p.sq(i - y - h + r));
        let bot_x = x + w - r + p.sqrt(p.sq(r) - p.sq(i - y - h + r));
        p.line(top_x, i, bot_x, i);
    }
}