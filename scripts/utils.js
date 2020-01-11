function isLight(colour) {
    return lightness(colour) > 50;
}

function sortMap(map) {
    let tupleArray = [];
    for (let key in map)
        if (map.hasOwnProperty(key)) {
            tupleArray.push(key);
        }
    tupleArray.sort((a, b) => (int(a) > int(b)) ? 1 : ((int(b) > int(a)) ? -1 : 0));
    return tupleArray;
}

function connectToServer() {
    if (ws) {
        return
    }
    let loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    } else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.host + "/battler";
    ws = new WebSocket(new_uri);

    ws.onopen = function (evt) {
        console.log("OPEN");
        connected = true;
    };
    ws.onclose = function (evt) {
        console.log("CLOSE");
        ws = null;
    };
    ws.onmessage = function (evt) {
        console.log("RESPONSE:");
        //JSONED = evt.data;
        let battleresponse = JSON.parse(evt.data);
        console.log(battleresponse);
        parseState(battleresponse);
    };
    ws.onerror = function (evt) {
        let errs = getElement("info");
        errs.setColour(color(red2));
        errs.setText("Failed to connect to the server!");
        console.log("ERROR: " + evt);
    };
}

function countdown(value, where, yesno) {
    timeleft = value - 1;
    if (timeleft < 0) {
        if (yesno) {
            console.log("RELOCATED!~", where);
            window.location = where;
            return;
        }
        isTicking = false;
        console.log("turn times out");
        displayTimer(-1);
    } else {
        displayTimer(timeleft);
        console.log(timeleft);
        window.setTimeout("countdown(timeleft, redirectwhere, doredirect)", 1000);
    }
}

function settimeleft(value) {
    timeleft = value;
}

function setwhere(where) {
    redirectwhere = where;
}

function redirect(yesno) {
    doredirect = yesno;
}

function keyPressed() {
    if (key === ' ' && getElement("back") && getElement("back").clickable) {
        getElement("back").clicked()
    }
    if (!Are_buttons_disabled) {
        if ((key === 'q' || key === 'a') && getElement("Q").clickable) {
            getElement("Q").clicked()
        } else if ((key === 'w' || key === 'z') && getElement("W").clickable) {
            getElement("W").clicked()
        } else if (key === 'e' && getElement("E").clickable) {
            getElement("E").clicked()
        } else if (key === 'r' && getElement("R").clickable) {
            getElement("R").clicked()
        }
    }

}

function getTurnNum(i) {
    return Math.floor((i - 1) / 2 + 1);
}

function getChar(event) {
    if (event.which == null) { // IE
        if (event.keyCode < 32) return "";
        return String.fromCharCode(event.keyCode)
    }
    if (event.which !== 0 && event.charCode !== 0) { // not IE

        if (event.which < 32) return ""; // symb
        return String.fromCharCode(event.which); // the rest
    }
    return ""; // symb

}

function sendSkill(Skill) {
    if (connected) {
        if (Are_buttons_disabled) {
            return
        }
        disableButtons(0);
        ws.send(JSON.stringify(Skill));
    } else {
        let info = getElement("info");
        info.setColour(red2);
        info.setText("You are not connected to the server.");
        console.log("Not connected yet!\n");
        connectToServer();
    }
}

function getElement(id) {
    for (obj of leftPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
    for (obj of rightPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
    for (obj of topPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
    for (obj of bottomPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
}

function roundUp(num) {
    return Math.round(num * 10) / 10
}

function calculateLines(hoverText, width, size) {
    if (!width) {
        width = 290;
    }
    if (!size) {
        size = 15;
    }
    let height = 0;
    textSize(size);
    for (let line of hoverText.split("\n")) {
        height += 1;
        let x_pos = 0;
        for (let word of line.split(" ")) {
            /*strokeWeight(5);
            stroke(red2);
            point(mouseX+10+x_pos + change, mouseY+height*size);*/
            if (x_pos + textWidth(word + " ") < width) { //we are still on that line
                x_pos += textWidth(word + " ")
            } else { //start a new line
                height += 1;
                x_pos = textWidth(word + " ");
            }
        }
    }
    return height;
}

function displayStandardHoverBubble(hoverText, lines) {
    let hoverSize = 15;
    textAlign(LEFT);
    textSize(hoverSize);
    noStroke();
    fill(hoverc);
    let changerForFlipping;
    let w = textWidth(hoverText);
    if (mouseX >= 965 && w < 290) {
        changerForFlipping = -w - 20;
    } else if (mouseX >= 965) {
        changerForFlipping = -310;
    } else {
        changerForFlipping = 0;
    }
    let changerForFlippingY;
    let h = lines * hoverSize + (lines - 1) * 5;
    if (mouseY + h + 10 >= 550) {
        changerForFlippingY = -h;
    } else {
        changerForFlippingY = 0;
    }
    if (h === hoverSize && w < 290) {
        rect(mouseX + changerForFlipping, mouseY + changerForFlippingY, w + 20, h + 10, 5);
    } else {
        rect(mouseX + changerForFlipping, mouseY + changerForFlippingY, 310, h + 10, 5);
    }
    strokeWeight(0.5);
    stroke(dark);
    fill(dark);
    textLeading(20);
    text(hoverText, mouseX + 10 + changerForFlipping, mouseY + changerForFlippingY + 5, 290);
}

function calculateHPperFrame(HP, targetHP) {
    let totalFrames = 60;
    let num = targetHP - HP;
    if (num > 0) {
        return Math.ceil(num / totalFrames)
    } else if (num < 0) {
        return Math.floor(num / totalFrames)
    } else {
        console.log("ERROR WITH SPEED", HP, targetHP);
        return 0;
    }
}
