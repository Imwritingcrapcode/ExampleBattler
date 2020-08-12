'use strict';
let response;
let currentGirls = [];
//let num;
let tags = [];
let colours = [];
let rarities = [];
let curNum = -1;
let mainNum = -1;
let secNum = -1;
let clearClickable = false;
let ws;
UpdateFreeData();
let bottomReady = false;
new p5(leftSketch, 'girls');
new p5(bottomSketch, 'bottomcanvas');
new p5(rightSketch, 'rightcanvas');

document.addEventListener("keydown", e => {
   if (e.code === 'Enter' || e.code === "Space") {
       e.preventDefault();
       battle();
   } else if (e.code === 'Digit1') {
       e.preventDefault();
       setasmain();
   } else if (e.code === 'Digit2') {
       e.preventDefault();
       setassec();
   } else if (e.code === 'Digit3' || e.code === 'KeyC') {
       e.preventDefault();
       clearsetgirls();
   } else if (e.code === 'KeyR') {
       e.preventDefault();
       random();
   }
});

function getGirlByNumber(number) {
    for (let girl of response) {
        if (girl.Number === number) {
            return girl
        }
    }
}

function refreshCurGirl(number) {
    curNum = number;
    let curGirl = getGirlByNumber(number);
    //put her info in the doc
    getElementRight("name").setText(curGirl.Name + " (" + curGirl.Number + ")");
    getElementRight("rarity").setText(curGirl.Rarity);
    switch (curGirl.Rarity) {
        case "LF":
            //star
            getElementRight("rarity").setColour(LFCOLOUR);
            break;
        case "RP":
            //green
            getElementRight("rarity").setColour(RPCOLOUR);
            break;
        case "SP":
            //yellow
            getElementRight("rarity").setColour(SPCOLOUR);
            break;
        case "AD":
            //blue
            getElementRight("rarity").setColour(ADCOLOUR);
            break;
        case "ST":
            //white
            getElementRight("rarity").setColour(STCOLOUR);
            break;
        default:
            console.log("EXCUSE ME");
            break;
    }
    getElementRight("level").setText("Level: " + curGirl.Level);
    getElementRight("matchesPlayed").setText("Matches played: " + curGirl.MatchesPlayed);
    let winrate = "";
    if (curGirl.MatchesPlayed !== 0) {
        winrate = " (" + Math.floor(curGirl.MatchesWon / curGirl.MatchesPlayed * 100) + "%)";
    }
    getElementRight("matchesWon").setText("Matches won: " + curGirl.MatchesWon + winrate);
    //let name = "Speed";
    let name = curGirl.Name;
    if (!imageBox.contains(name)) {
        if (existsPortrait(curGirl.Number)) {
            imageBox.add(new InterfaceImage(rightP, 519.5, 5, "/images/locked/" + name + "_portrait.png", "girl", name, 201, 268, curGirl.MainColour));
        } else {
            imageBox.add(new InterfaceImage(rightP, 519.5, 5, "/images/locked/Placeholder_portrait.png", "girl", name, 201, 268, curGirl.MainColour));
        }
    }
    //setloadingscreencolours;
    p_screen.setColours(curGirl.SkillColourCodes[0], curGirl.SkillColourCodes[1], curGirl.SkillColourCodes[2], curGirl.SkillColourCodes[3]);
    // update the pic
    getElementRight("girl").set(imageBox.get(name));
    //p_screen.restart()
    p_screen.restart();
    let tags = "";
    for (let tag of curGirl.Tags.sort()) {
        tags += tag + ", "
    }
    tags = tags.slice(0, tags.length - 2);
    getElementRight("tags").setText("Tags: " + tags);
    let skills = "";
    for (let s of curGirl.Skills) {
        skills += s + ", "
    }
    skills = skills.slice(0, skills.length - 2);
    getElementRight("skills").setText("Skills: " + skills);
    let colours = "";
    for (let s of curGirl.SkillColours) {
        colours += s + ", "
    }
    colours = colours.slice(0, colours.length - 2);
    getElementRight("skillcolours").setText("Skill colours: " + colours);
    let description = curGirl.Description;
    let lines = interfaceCalculateLines(rightP, description, 738, 28);
    if (lines > 1) {
        can3.resize(748, 413 + (lines - 1) * 33);
    } else {
        can3.resize(748, 413);
    }
    getElementRight("description").setText(description);
    //document.getElementById("description").textContent = curGirl.Description;*


    //check if she can be picked
    if (bottomReady) {
        getElementBottom("setmain").clickable = !(curNum === mainNum);
        getElementBottom("setsec").clickable = !(curNum === secNum);

    }
    if (curNum === mainNum && bottomReady) {
        getElementRight("set").setText("Set as Main!");
        //document.getElementById("mainset").textContent = "Set as Main!";
    } else if (curNum === secNum && bottomReady) {
        getElementRight("set").setText("Set as Secondary!");
        //document.getElementById("secset").textContent = "Set as Secondary!";
    } else if (bottomReady) {
        getElementRight("set").setText("");
        //document.getElementById("secset").textContent = "";
    }
}

function display(girllist) {
    let textS = 45;
    let height = (75) * girllist.length + 5;
    if (height > canHeight && can.height !== height) {
        console.log("YES", girllist.length, height);
        can.resize(canWidth, height);
    } else if (can.height !== height) {
        console.log("NO", girllist.length, height);
        can.resize(canWidth, canHeight);
    }
    emptyLeft();
    for (let i = 0; i < girllist.length; i++) {
        let girl = girllist[i];
        let t = girl.Name;
        leftP.textSize(textS);
        let w = leftP.textWidth(t);
        let button = new InterfaceButton(leftP, 5, (75) * i + 5, t, textS, "girl" + i, "B", 250);
        button.onClick = function () {
            this.index = girl.Number;
            refreshCurGirl(this.index);
        };
        let c = leftP.lerpColor(leftP.color(girl.MainColour), light, 0.45);
        button.setColour(c.toString());
        leftobjects.push(button);
    }
}

function okayforfilter(girl, filter, type) {
    return (type === "colour" && girl.SkillColours.indexOf(filter) !== -1 ||
        type === "rarity" && girl.Rarity === filter ||
        type === "tag" && girl.Tags.indexOf(filter) !== -1);
}

function checkallothers(girl, filter, type) {
    if (type === "colour") {
        for (let colour of girl.SkillColours) {
            if (colour !== filter && document.getElementById(colour).checked) {
                return true;
            }
        }
        return false;
    } else if (type === "tag") {
        for (let tag of girl.Tags) {
            if (tag !== filter && document.getElementById(tag).checked) {
                return true;
            }
        }
        return false;

    }
}

function oneofall(girl) {
    if (!document.getElementById(girl.Rarity).checked) {
        return false;
    } else {
        let verdict = false;
        for (let colour of girl.SkillColours) {
            if (document.getElementById(colour).checked) {
                verdict = true;
                break;
            }
        }
        if (!verdict) {
            return verdict;
        }
        verdict = false;
        for (let tag of girl.Tags) {
            if (document.getElementById(tag).checked) {
                verdict = true;
                break;
            }
        }
        return verdict;
    }
}

function updatedisplay(filter, type) {
    //1. get a list of everything needed
    let tyanlist = [];
    let checked = document.getElementById(filter).checked;

    if (checked) {
        for (let girl of response) {
            if (oneofall(girl)) {
                tyanlist.push(girl);
            }
        }
    } else {
        for (let girl of currentGirls) {
            if (!okayforfilter(girl, filter, type) || okayforfilter(girl, filter, type) &&
                type !== "rarity" && checkallothers(girl, filter, type)) {
                tyanlist.push(girl);
            }
        }
    }

    //2. display it
    currentGirls = tyanlist.slice(0);
    display(currentGirls);

}

function parse(response) {
    let presets = getCookie("GirlList");
    if (presets !== null) {
        let arr = presets.split(" ");
        mainNum = +arr[0];
        secNum = +arr[1];
        if (bottomReady) {
            getElementBottom("clear").clickable = true;

        }
        clearClickable = true;
    } else {
        if (bottomReady) {
            getElementBottom("clear").clickable = false;
        }
        clearClickable = false;
    }

    display(response);
    refreshCurGirl(response[0].Number);

    for (let girl of response) {
        //sobiraem colors
        for (let colour of girl.SkillColours) {
            if (colours.indexOf(colour) === -1) {
                colours.push(colour);
            }
        }

        //sobiraem rarities
        if (rarities.indexOf(girl.Rarity) === -1) {
            rarities.push(girl.Rarity);
        }

        //sobiraem tags
        for (let tag of girl.Tags) {
            if (tags.indexOf(tag) === -1) {
                tags.push(tag);
            }
        }
    }
    //sorting
    colours.sort((a, b) => (colourMap.get(a) > colourMap.get(b)) ? 1 : ((colourMap.get(b) > colourMap.get(a)) ? -1 : 0));
    rarities.sort((a, b) => (raritiesMap.get(a) > raritiesMap.get(b)) ? 1 : ((raritiesMap.get(b) > raritiesMap.get(a)) ? -1 : 0));
    tags.sort();

    //filters.
    let i = 4;
    for (let colour of colours) {
        document.getElementById("coloursfilter").innerHTML += "<input onchange='updatedisplay(\"" + colour
            + "\",\"colour\"" + ")' type=\"checkbox\" id=\"" + colour
            + "\" checked><label for=\"" + colour + "\"> " + colour + " </label>";
        i--;
        if (i % 4 === 0) {
            document.getElementById("coloursfilter").innerHTML += "<br>";
        }
    }

    for (let rarity of rarities) {
        document.getElementById("raritiesfilter").innerHTML += "<input onchange='updatedisplay(\"" + rarity
            + "\",\"rarity\"" + ")' type=\"checkbox\" id=\"" + rarity
            + "\" checked><label for=\"" + rarity + "\"> " + rarity + " </label><br>";
    }
    i = 3;
    for (let tag of tags) {
        document.getElementById("tagsfilter").innerHTML += "<input onchange='updatedisplay(\"" + tag
            + "\",\"tag\"" + ")' type=\"checkbox\" id=\"" + tag
            + "\" checked><label for=\"" + tag + "\"> " + tag + " </label>";
        i--;
        if (i % 3 === 0) {
            document.getElementById("tagsfilter").innerHTML += "<br>";
        }
    }
    console.log("parse is rdy");
}

function getgirllist() {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/girllist', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                let all = JSON.parse(xhr.responseText);
                response = all.girls.sort((a, b) => (a.Number > b.Number) ? 1 : ((b.Number > a.Number) ? -1 : 0));
                if (bottomReady) {
                    getElementBottom("prompts").setText(all.response);
                }
                if (response.length > 0) {
                    console.log(response);
                    currentGirls = response.slice(0);
                    parse(response);
                } else {
                    console.log("0 girls found!!!")
                }
            } else if (xhr.status === 303) {
                window.location.reload(true);
            }
        }
    };
}

function random() {
    let old_prompt = "";
    if (bottomReady) {
        getElementBottom('clear').clickable = true;
        old_prompt = getElementBottom("prompts").text;
    }
    let len = response.length;
    let rand = Math.floor(Math.random() * len);
    let rand2 = Math.floor(Math.random() * len);
    while (rand2 === rand) {
        rand2 = Math.floor(Math.random() * len);
    }
    let wasMain = false;
    let wasSec = false;
    if (mainNum === curNum) {
        wasMain = true;
    }
    if (secNum === curNum) {
        wasSec = true;
    }
    mainNum = response[rand].Number;
    secNum = response[rand2].Number;
    console.log(mainNum, secNum, rand, rand2, response.length);
    display(currentGirls);
    if ((wasMain && mainNum !== curNum) || (wasSec && secNum !== curNum) ||
        (!wasMain && mainNum === curNum) || (!wasSec && secNum === curNum)) {
        refreshCurGirl(curNum);
    } else {
        console.log("no need to refresh", curNum);
    }

    if (old_prompt !== "") {
        getElementBottom("prompts").setText(old_prompt);
    }
}

function setasmain() {
    if (curNum === secNum) { //if we try to set the sec as main, swap them
        secNum = mainNum;
        getElementBottom("setsec").clickable = true;
        //clearsetgirls();

    }
    mainNum = curNum;
    if (bottomReady) {
        getElementBottom("setmain").clickable = false;
        getElementBottom("clear").clickable = true;

    }
    getElementRight("set").setText("Set as Main!")
    //document.getElementById("mainset").textContent = "Set as Main!";
}

function setassec() {
    if (curNum === mainNum) {
        mainNum = secNum;
        getElementBottom("setmain").clickable = true;
        //clearsetgirls();
    }
    secNum = curNum;
    if (bottomReady) {
        getElementBottom("setsec").clickable = false;
        getElementBottom("clear").clickable = true;

    }
    getElementRight("set").setText("Set as Secondary!")
    //document.getElementById("secset").textContent = "Set as Secondary!";
}

function clearsetgirls() {
    mainNum = -1;
    secNum = -1;
    if (bottomReady) {
        getElementBottom("setmain").clickable = true;
        getElementBottom("setsec").clickable = true;
        getElementBottom("clear").clickable = false;
    }

    getElementRight("set").setText("");
    //document.getElementById("mainset").textContent = "";
    //document.getElementById("secset").textContent = "";

}

function battle() {
    if (mainNum !== -1 && secNum !== -1) {
        if (ws) {
            if (bottomReady) {
            getElementBottom("prompts").setText("Stopped the search.");
            getElementBottom('battle').setText("Battle");
        }
            ws.close();
            return
        }

        let loc = window.location, new_uri;
        if (loc.protocol === "https:") {
            new_uri = "wss:";
        } else {
            new_uri = "ws:";
        }
        new_uri += "//" + loc.host + "/queue";
        ws = new WebSocket(new_uri);

        ws.onopen = function (evt) {
            if (bottomReady) {
                getElementBottom("prompts").setText("Searching for an opponent...");
                getElementBottom('battle').setText("Cancel");
            }
            console.log("OPEN");
            console.log("SEND: ", JSON.stringify({MainGirl: mainNum, SecondaryGirl: secNum,}));
            ws.send(JSON.stringify({MainGirl: mainNum, SecondaryGirl: secNum,}));
        };
        ws.onclose = function (evt) {
            console.log("CLOSE");
            ws = null;
        };
        ws.onmessage = function (evt) {
            console.log("RESPONSE: " + evt.data);
            let battleresponse = JSON.parse(evt.data);
            if (bottomReady) {
                getElementBottom("prompts").setText(battleresponse.Prompt);
            }
            if (battleresponse.hasOwnProperty("OK") && battleresponse.OK) {
                console.log("RELOCATE~!!");
                setCookie("GirlList", mainNum + " " + secNum, 3);
                location = battleresponse.Location;
            } else {
                console.log("Server dced me :<");
                ws.close();
            }
        };
        ws.onerror = function (evt) {
            console.log("ERROR: " + evt.toString());
        };


    } else {
        if (bottomReady) {
            getElementBottom("prompts").setText("Please choose the girls first!");
        }
    }
}
