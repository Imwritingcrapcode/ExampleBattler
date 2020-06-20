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

function is_touch_device4() {

    let prefixes = ' -webkit- -moz- -o- -ms- '.split(' ');

    let mq = function (query) {
        return window.matchMedia(query).matches;
    };

    if (('ontouchstart' in window) || window.DocumentTouch && document instanceof DocumentTouch) {
        return true;
    }

    // include the 'heartz' as a way to have a non matching MQ to help terminate the join
    // https://git.io/vznFH
    let query = ['(', prefixes.join('touch-enabled),('), 'heartz', ')'].join('');
    return mq(query);
}

function countdown(value, where, yesno) {
    timeleft = value - 1;
    if (timeleft < 0) {
        if (yesno) {
            console.log("RELOCATED!~", where);
            //window.location = where;
            return;
        }
        isTicking = false;
        console.log("countdown times out");
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

function parseSeconds(n, strip) {
    let mins = Math.floor(n / 60);
    let rem_secs = n - mins * 60;
    let hrs = Math.floor(mins / 60);
    let rem_mins = mins - hrs * 60;
    let days = Math.floor(hrs / 24);
    let rem_hrs = hrs - days * 24;
    let full = "";
    if (days > 0) {
        if (days === 1) {
            full += days + " day";
        } else {
            full += days + " days";
        }
        if (strip) {
            if (days > 364) {
                return "-1"
            } else {
                return full
            }
        } else {
            full += ", "
        }
    }
    if (rem_hrs > 0) {
        if (rem_hrs === 1) {
            full += rem_hrs + " hour";
        } else {
            full += rem_hrs + " hours";
        }
        if (strip) {
            return full
        } else {
            full += ", "
        }
    }
    if (rem_mins > 0) {
        if (rem_mins === 1) {
            full += rem_mins + " minute";
        } else {
            full += rem_mins + " minutes";
        }
        if (strip) {
            return full
        } else {
            full += ", "
        }
    }
    if (rem_secs >= 0) {
        if (rem_secs === 1) {
            full += rem_secs + " second";
        } else {
            full += rem_secs + " seconds";
        }
        if (strip) {
            return "less than a minute"
        }
    }

    return full;
}

function UpdateFreeData(after) {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/freeinfo', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                let response = JSON.parse(xhr.responseText);
                console.log(response);
                let welcome = "Welcome, " + response.Username;
                if ( welcome.length > 19) {
                    document.getElementById("username").innerText = "Hi, " + response.Username;
                } else {
                    document.getElementById("username").innerText = welcome;
                }
                if (!!document.getElementById("moneytable")) {
                    document.getElementById("wDust").innerText = response.MoneyInfo["w"];
                    document.getElementById("bDust").innerText = response.MoneyInfo["b"];
                    document.getElementById("yDust").innerText = response.MoneyInfo["y"];
                    document.getElementById("pDust").innerText = response.MoneyInfo["p"];
                    document.getElementById("sDust").innerText = response.MoneyInfo["s"];
                }
                if (!!after) {
                    after(response);
                }
            } else {
                console.log(xhr.responseText);
            }
        }
    };
}

function UpdateProfileData() {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/profileinfo', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                let response = JSON.parse(xhr.responseText);
                console.log(response);
                let welcome = "Welcome, " + response.Username;
                if ( welcome.length > 19) {
                    document.getElementById("username").innerText = "Hi, " + response.Username;
                } else {
                    document.getElementById("username").innerText = welcome;
                }
                document.getElementById("wDust").innerText = response.MoneyInfo["w"];
                document.getElementById("bDust").innerText = response.MoneyInfo["b"];
                document.getElementById("yDust").innerText = response.MoneyInfo["y"];
                document.getElementById("pDust").innerText = response.MoneyInfo["p"];
                document.getElementById("sDust").innerText = response.MoneyInfo["s"];
                document.getElementById("username2").innerText = response.Username;
                if (response.BattlesTotal > 0) {
                    document.getElementById("battles2").innerText = "Battle stats: " + response.BattlesWon + "/" + response.BattlesTotal + " (" + roundUp(response.BattlesWon / response.BattlesTotal * 100) + "% winrate)";
                } else {
                    document.getElementById("battles2").innerText = "Battle stats: " + 0 + "/" + 0 + " (" + 0 + "% winrate)";
                }
            } else {
                console.log(xhr.responseText);
            }
        }
    };
}

