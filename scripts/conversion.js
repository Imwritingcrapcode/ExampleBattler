function setup() {
    updatespersecond = 60;
    timeleft = -1;
    function handleVisibilityChangeConv() {
        if (!document.hidden && CONVINFO.IsConvertingRN) {
            convert("?");
        }
    }
    document.addEventListener("visibilitychange", handleVisibilityChangeConv, false);
    frameRate(updatespersecond);
    convobjects = [];
    current = undefined;
    touch = is_touch_device4();
    let can = createCanvas(1024, 70);
    can.parent('items');
    bg_color = color(BG);
    dark = color(DARKC);
    right = color(RIGHTC);
    light = color(LIGHTC);
    blueDust = color(ADCOLOUR);
    yellowDust = color(SPCOLOUR);
    purpleDust = color(RPCOLOUR);
    starDust = color(LFCOLOUR);
    maxValue = 0;
    currentType = "";
    let text = "Convert!";
    let tSize = 40;
    bar = new LoadingBar(10, 10, 512, 50, 10, "bar", color(dark.toString()), color(light.toString()));
    bar.clicked = function (x) {
        let p = (x - this.x) / this.width;
        if (1 - p < 0.01) {
            p = 1
        }
        this.setPercentage(p * 100);
        this.setNewPercentage(p * 100);
        document.getElementById("number").innerText = Math.ceil(p * maxValue) + "";
    };
    convobjects.push(bar);
    textSize(tSize);
    let B = new StandardButton(768 - (textWidth(text) + 10) / 2, 9, 10, text, tSize, "c");
    B.hide();
    B.clicked = function () {
        let inner = document.getElementById("number").innerText;
        let get = int(Math.floor(CONVINFO.ConversionRate[currentType] * int(inner)));
        let cost = get / CONVINFO.ConversionRate[currentType];
        if (cost < 1 || get < 1 || get === 1 && cost !== Math.trunc(cost)) {
            while (cost !== Math.trunc(cost) || get < 1) {
                get++;
                cost = get / CONVINFO.ConversionRate[currentType];
            }
            document.getElementById("number").innerText = "You need to convert at least " + cost + " dust!";
            return
        }
        if (inner > 0 && get >= 1.0) {
            this.hide();
            bar.makeNotDraggable();
            document.getElementById("w").disabled = true;
            document.getElementById("b").disabled = true;
            document.getElementById("y").disabled = true;
            document.getElementById("p").disabled = true;
            document.getElementById("s").disabled = true;
            convert("!", int(inner), currentType);
        }
    };
    convobjects.push(B);
    let t2 = "Claim!";
    let bw2 = textWidth(t2) + 10;
    let claim = new StandardButton(768 - bw2 / 2, 9, 10, t2, tSize, "cl");
    claim.hide();
    claim.clicked = function () {
        this.hide();
        convert("!");
    };
    convobjects.push(claim);

    let willGet = new TextInfo(768, 60, dark, "", 50, "willGet");
    convobjects.push(willGet);
    willGet.hide();
    convert("?");
}

function mousePressed() {
    let x = mouseX;
    let y = mouseY;
    for (obj of convobjects) {
        if (obj.clickable && obj.in(x, y)) {
            obj.clicked(x);
        }
    }
}

function mouseDragged() {
    let x = mouseX;
    let y = mouseY;
    for (obj of convobjects) {
        if (obj.clickable && obj.draggable && obj.in(x, y)) {
            obj.clicked(x);
        }
    }
}

function draw() {
    background(bg_color);
    for (let obj of convobjects) {
        if (obj.clickTimer > 0) {
            obj.clickTimer--;
            if (obj.clickTimer === 0) {
                obj.unclick();
            }
            obj.display();
        } else if (obj.hoverable && obj.in()) { //found an "in"
            if (!current) { //outside to something
                current = obj;
                obj.hovered();
                obj.display();
            } else if (current.id === obj.id) { //currently hovered
                obj.display();
            } else { //switched from another 2 this
                current.unhovered();
                current = obj;
                obj.hovered();
                obj.display();
            }
        } else if (obj.hoverable && current && obj.id === current.id) { //went outside
            obj.unhovered();
            current = undefined;
            obj.display();
        } else {
            obj.display();
        }
    }
    if (timeleft > 0) {
        countdownGentle(timeleft)
    }
}

function getElement(id) {
    for (let obj of convobjects) {
        if (obj.id === id) {
            return obj
        }
    }
}

function removeElement(id) {
    for (let i = 0; i < convobjects.length; i++) {
        let obj = convobjects[i];
        if (obj.id === id) {
            convobjects.splice(i, 1);
            return;
        }
    }
}

function convert(requestType, amount, dustType) {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/conversion', true);
    xhr.send(JSON.stringify({
        ReqType: requestType,
        Amount: amount,
        DustType: dustType
    }));
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                console.log(xhr.responseText);
                CONVINFO = JSON.parse(xhr.responseText);
                let after = function (data) {
                    MONIES = new Map(Object.entries(data.MoneyInfo));
                    if (!CONVINFO.ConversionRate) {
                        CONVINFO.ConversionRate = CONVERSIONRATE
                    }
                    parse(CONVINFO, MONIES);
                };
                UpdateFreeData(after);
            } else if (xhr.status === 400) {
                console.log(xhr.responseText);
                convert("?");
            }
        }
    };
}

function setDustType(type, amnt) {
    bar.setPercentage(0.0);
    bar.setNewPercentage(0.0);
    maxValue = amnt;
    switch (type) {
        case "w":
            bar.setColours(color(dark.toString()), color(light.toString()));
            document.getElementById("w").checked = true;
            document.getElementById("b").checked = false;
            document.getElementById("y").checked = false;
            document.getElementById("p").checked = false;
            document.getElementById("s").checked = false;
            break;
        case "b":
            bar.setColours(color(dark.toString()), lerpColor(color(blueDust.toString()), color(255), 0.3));
            document.getElementById("w").checked = false;
            document.getElementById("b").checked = true;
            document.getElementById("y").checked = false;
            document.getElementById("p").checked = false;
            document.getElementById("s").checked = false;
            break;
        case "y":
            bar.setColours(color(dark.toString()), lerpColor(color(yellowDust.toString()), color(255), 0.2));
            document.getElementById("w").checked = false;
            document.getElementById("b").checked = false;
            document.getElementById("y").checked = true;
            document.getElementById("p").checked = false;
            document.getElementById("s").checked = false;
            break;
        case "p":
            bar.setColours(color(dark.toString()), lerpColor(color(purpleDust.toString()), color(255), 0.2));
            document.getElementById("w").checked = false;
            document.getElementById("b").checked = false;
            document.getElementById("y").checked = false;
            document.getElementById("p").checked = true;
            document.getElementById("s").checked = false;
            break;
        case "s":
            bar.setColours(color(dark.toString()), lerpColor(color(starDust.toString()), color(255), 0.2));
            document.getElementById("w").checked = false;
            document.getElementById("b").checked = false;
            document.getElementById("y").checked = false;
            document.getElementById("p").checked = false;
            document.getElementById("s").checked = true;
            getElement("c").clickable = false;
            break;
        default:
            console.log("why is this happening to me.", type);
            return;
    }
    currentType = type;

}

function setMoney(conv, m, isConverting) {
    let needToCheck = true;
    let firstDust = "w";
    let firstAmnt = m.get("w");
    let allInnerHtml = "";
    for (let dustType of m.keys()) {
        let amnt = m.get(dustType);
        if (dustType === "s") {
            allInnerHtml += "<input onchange='setDustType(\"" + dustType + "\", " + amnt + ")' type=\"checkbox\" id =\"" + dustType +
                "\" disabled><label for = '" + dustType + "'>\t" + DUSTS.get(dustType) + "dust:\t\t\t" + amnt + "</label><br>";
        } else {
            let get;
            if (!!conv) {
                get = parseFloat(conv[dustType]) * int(amnt);
            } else {
                get = 0;
            }
            let cost = get / conv[dustType];
            let d;
            if (isConverting || cost < 1 || get < 1 || get === 1 && cost !== Math.trunc(cost)) {
                d = "disabled";
            } else {
                d = "";
                if (needToCheck) {
                    firstDust = dustType;
                    firstAmnt = amnt;
                    needToCheck = false;
                }
            }
            allInnerHtml += "<input onchange='setDustType(\"" + dustType + "\", " + amnt + ")' type=\"checkbox\" id =\"" + dustType +
                "\"" + d + "><label for = '" + dustType + "'>\t" + DUSTS.get(dustType) + " dust:\t\t" + amnt + "</label><br>";
        }
        document.getElementById("rarityselect").innerHTML = allInnerHtml;
    }
    setDustType(firstDust, firstAmnt);
}

function parse(r, m) {
    //possible states: not converting yet, converting, rdy to claim
    if (!r.IsConvertingRN) { //NOT CONVERTING YET
        timeleft = -1;
        document.getElementById("number").innerText = "";
        let c = getElement("c");
        if (!c.visible) {
            c.setText("Convert!");
            getElement("c").show();
        }
        setMoney(r.ConversionRate, m, r.IsConvertingRN);
        bar.makeDraggable();
    } else if (r.IsConvertingRN && r.Left > 0) { //CONVERTING
        console.log("currently converting...");
        let c = getElement("c");
        if (c.visible) {
            c.hide();
        }
        setMoney(r.ConversionRate, m, r.IsConvertingRN);
        bar.makeNotDraggable();
        let left = int(r.Left);
        bar.total = int(r.CurrentProgress) + left;
        bar.left = left;
        let get = r.Amount;
        let dustType = r.DustType;
        setDustType(dustType);
        textSize(50);
        let bw = textWidth(get);
        let emntEl = getElement("willGet");
        emntEl.x = 768 - bw / 2 + 65 / 2;
        emntEl.setText(get);
        if (!emntEl.visible) {
            emntEl.show();
        }
        let im = new CanvasImage(emntEl.x - 65, emntEl.y - 50, "/images/locked/" + DUSTS.get(dustType) + "_dust.png", "dustpic", DUSTS.get(dustType) + "_dust.png", 50, 50);
        convobjects.push(im);
        timeleft = left;
    } else if (r.IsConvertingRN && r.Left === 0) {//RDY TO CLAIM
        timeleft = -1;
        console.log("available for claiming!");
        bar.makeNotDraggable();
        let get = r.Amount;
        let n = document.getElementById("number");
        let dustType = r.DustType;
        n.innerHTML = " <img src=\"/images/locked/" + DUSTS.get(dustType) + "_dust.png\" style=\"width:40px;height:40px;\" alt=\"" + DUSTS.get(dustType) + " dust\"><br>" + get;
        let emntEl = getElement("willGet");
        if (emntEl.visible) {
            emntEl.hide();
        }
        let conv = getElement("c");
        if (conv.visible) {
            conv.hide();

        }
        removeElement("dustpic");
        let c = getElement("cl");
        if (!c.visible) {
            c.show();
        }
        setMoney(r.ConversionRate, m, r.IsConvertingRN);
        setDustType(dustType, m.get(dustType));
        bar.setPercentage(100);
    } else { //ERROR
        console.log("uhhhhh what?")
    }
}

function displayTimer(number) {
    let timer = document.getElementById("number");
    if (number < 0) {
        timer.innerText = "";
        convert("?");
    } else {
        timer.innerText = "" + parseSeconds(Math.ceil(number)) + " left.";
        bar.left = number;
        bar.setPercentage((bar.total - bar.left) / bar.total * 100);
    }
}

function countdownGentle(value) {
    timeleft = value - 1/updatespersecond;
    if (timeleft < 0) {
        console.log("countdown times out");
        displayTimer(-1);
    } else {
        displayTimer(timeleft);
    }
}