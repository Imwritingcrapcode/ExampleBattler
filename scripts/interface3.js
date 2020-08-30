var touch;
var FS;

function setup() {
    //textFont('Calibri');
    //TESTING = true;
    TESTING = false;
    touch = is_touch_device4();
    //touch = true;
    PICS = true;
    //PICS = false;
    if (TESTING) {
        STATE = "{\"Instruction\":\"\",\"TurnNum\":4,\"TurnPlayer\":1,\"PlayerNum\":51,\"OppNum\":33,\"PlayerName\":\"Milana\",\"OppName\":\"Speed\",\"HP\":114,\"MaxHP\":114,\"OppHP\":113,\"OppMaxHP\":113,\"Effects\":{\"8\":\"2\"},\"OppEffects\":{\"13\":\"3\",\"6\":\"10\"},\"SkillState\":{\"E\":2,\"Q\":0,\"R\":-1,\"W\":0},\"OppSkillState\":{\"OppE\":0,\"OppQ\":0,\"OppR\":0,\"OppW\":0},\"SkillNames\":{\"E\":\"Mint Mist\",\"Q\":\"Royal Move\"" +
            ",\"R\":\"Pride\",\"W\":\"Composure\"},\"OppSkillNames\":{\"OppE\":\"Speed\",\"OppQ\":\"Run\",\"OppR\":\"Stab\",\"OppW\":\"Weaken\"},\"Defenses\":{\"1\":0,\"10\":0,\"11\":-2,\"12\":2,\"2\":-1,\"3\":1,\"4\":1,\"5\":1,\"6\":0,\"7\":-1,\"8\":0,\"9\":0},\"OppDefenses\":{\"1\":0,\"10\":0,\"11\":2,\"12\":-2,\"2\":0,\"3\":0,\"4\":4,\"5\":0,\"6\":0,\"7\":-2,\"8\":0,\"9\":0},\"SkillColours\":{\"E\":\"rgb(232,255,243)\",\"Q\":\"rgb(49,255,185)\"," +
            "\"R\":\"rgb(115,255,240)\",\"W\":\"rgb(232,255,243)\"},\"OppSkillColours\":{\"OppE\":\"rgb(14,51,20)\",\"OppQ\":\"rgb(14,51,20)\",\"OppR\":\"rgb(0,10,0)\",\"OppW\":\"rgb(0,10,0)\"},\"EndState\":0}";
        S = JSON.parse(STATE);
        STATE2 = "{\"Instruction\":\"\",\"TurnNum\":17,\"TurnPlayer\":2,\"PlayerNum\":9,\"OppNum\":1,\"PlayerName\":\"Euphoria\",\"OppName\":\"Storyteller\",\"HP\":209,\"MaxHP\":209,\"OppHP\":207,\"OppMaxHP\":207,\"Effects\":{},\"OppEffects\":{\"12\":\"\",\"18\":\"6\"},\"SkillState\":{\"E\":0,\"Q\":0,\"R\":0,\"W\":0},\"OppSkillState\":{\"OppE\":0,\"OppQ\":1,\"OppR\":-100,\"OppW\":1},\"SkillNames\":{\"E\":\"Your Dream\",\"Q\":\"Your Number\",\"R\":\"My Story\",\"W\":\"Your Colour\"},\"OppSkillNames\":{\"OppE\":\"Pink Sphere\",\"OppQ\":\"Ampleness\",\"OppR\":\"Euphoria\",\"OppW\":\"Exuberance\"},\"Defenses\":{\"1\":-1,\"2\":1,\"3\":0,\"4\":-2,\"5\":-1,\"6\":1,\"7\":1,\"8\":0,\"9\":-1,\"10\":-1,\"11\":-2,\"12\":1},\"OppDefenses\":{\"1\":0,\"2\":2,\"3\":0,\"4\":0,\"5\":-3,\"6\":0,\"7\":0,\"8\":3,\"9\":0,\"10\":-4,\"11\":0,\"12\":0},\"SkillColours\":{\"E\":\"rgb(104,022,253)\",\"Q\":\"rgb(255,69,002)\",\"R\":\"rgb(29,104,255)\",\"W\":\"rgb(237,235,243)\"},\"OppSkillColours\":{\"OppE\":\"rgb(255,135,173)\",\"OppQ\":\"rgb(255,173,135)\",\"OppR\":\"rgb(255,135,175)\",\"OppW\":\"rgb(255,173,135)\"},\"EndState\":0}";
        S2 = JSON.parse(STATE2);
        STATE3 = "{\"Instruction\":\"\",\"TurnNum\":19,\"TurnPlayer\":2,\"PlayerNum\":9,\"OppNum\":1,\"PlayerName\":\"Euphoria\",\"OppName\":\"Storyteller\",\"HP\":221,\"MaxHP\":221,\"OppHP\":213,\"OppMaxHP\":219,\"Effects\":{},\"OppEffects\":{\"12\":\"\"},\"SkillState\":{\"E\":0,\"Q\":0,\"R\":0,\"W\":0},\"OppSkillState\":{\"OppE\":0,\"OppQ\":0,\"OppR\":-100,\"OppW\":0},\"SkillNames\":{\"E\":\"Your Dream\",\"Q\":\"Your Number\",\"R\":\"My Story\",\"W\":\"Your Colour\"},\"OppSkillNames\":{\"OppE\":\"Pink Sphere\",\"OppQ\":\"Ampleness\",\"OppR\":\"Euphoria\",\"OppW\":\"Exuberance\"},\"Defenses\":{\"1\":-1,\"2\":1,\"3\":0,\"4\":-2,\"5\":-1,\"6\":1,\"7\":1,\"8\":0,\"9\":-1,\"10\":-1,\"11\":-2,\"12\":1},\"OppDefenses\":{\"1\":0,\"2\":2,\"3\":0,\"4\":0,\"5\":-3,\"6\":0,\"7\":0,\"8\":3,\"9\":0,\"10\":-4,\"11\":0,\"12\":0},\"SkillColours\":{\"E\":\"rgb(104,022,253)\",\"Q\":\"rgb(255,69,002)\",\"R\":\"rgb(29,104,255)\",\"W\":\"rgb(237,235,243)\"},\"OppSkillColours\":{\"OppE\":\"rgb(255,135,173)\",\"OppQ\":\"rgb(255,173,135)\",\"OppR\":\"rgb(255,135,175)\",\"OppW\":\"rgb(255,173,135)\"},\"EndState\":0}";
        S3 = JSON.parse(STATE3);
        //S2 = JSON.parse(STATE);
    }
    //Battler vars
    isTicking = false;
    ws = undefined;
    connected = false;
    Are_buttons_disabled = true;
    timeleft = 0;
    doredirect = false;
    redirectwhere = "/girllist";

    //Canvas setup part
    bg_color = color(BG);
    light = color(LIGHTC);
    dark = color(DARKC);
    rightc = color(RIGHTC);
    clickc = color(CLICKABLEC);
    activec = color(ACTIVEC);
    hoverc = color(HOVERC);
    clickedc = color(CLICKEDC);
    red2 = color(RED);
    orange = color(ORANGE);
    yellow = color(YELLOW);
    green2 = color(GREEN);
    cyan = color(CYAN);
    blue2 = color(BLUE);
    violet = color(VIOLET);
    pink = color(PINK);
    grey = color(GREY);
    brown = color(BROWN);
    black = color(BLACK);
    white = color(WHITE);
    gc = color(WINC);
    other = color(OTHERLIGHTC);


    IMAGEBOX = new ImageBox();
    effDescs = new Map(EFFECTDESCRIPTIONS);
    current = undefined;


    let can = createCanvas(1280, 550);
    can.parent('interface');
    /*can.touchEnded(touchEnded2);
    can.touchMoved(touchMoved2);
    can.touchStarted(touchStarted2);*/
    can.mouseClicked(clicked);
    leftPanel = new Panel(0, 0, 550, 550, 5);
    rightPanel = new Panel(730, 0, 550, 550, 5);
    topPanel = new Panel(550, 0, 180, 230, 5);
    bottomPanel = new Panel(550, 230, 180, 320, 5);

    //top panel!
    topPanel.add(new TextInfo(555, 60, dark, "", 25, "timer", "", false, false, false));
    topPanel.add(new TextInfo(555, 30, dark, "", 25, "turnNumber", "", false, false, true));
    topPanel.add(new TurnLog(550, 230, dark, 20, "turnLog", 180, 175, true));
    //left panel!
    leftPanel.add(new TextInfo(5, 30, dark, "", 25, "playerName", "", false, false, false));
    //leftPanel.add(new TextInfo(200, 30, dark, "", 25, "techinfo", "", false, false, false));
    leftPanel.add(new TextInfo(5, 55, dark, "", 20, "playerHP", "", false, false, true));
    leftPanel.add(new TextInfo(5, 235, dark, "", 25, "effects", "effects", 215, undefined, true));
    leftPanel.add(new TextInfo(330, 235, dark, "", 25, "effects2", "effects", 215, undefined, true));
    leftPanel.add(new SkillButton(26, 265, 1, "", "Q", true));
    leftPanel.add(new SkillButton(132, 315, 2, "", "W", true));
    leftPanel.add(new SkillButton(286, 315, 3, "", "E", true));
    leftPanel.add(new SkillButton(440, 265, 4, "", "R", true));
    //right panel!
    rightPanel.add(new TextInfo(735, 30, dark, "", 25, "oppName", "", false, false, false));
    rightPanel.add(new TextInfo(735, 55, dark, "", 20, "oppHP", "", false, false, true));
    rightPanel.add(new TextInfo(735, 235, dark, "", 25, "oppEffects", "effects", 215, undefined, true));
    rightPanel.add(new TextInfo(1060, 235, dark, "", 25, "oppEffects2", "effects", 215, undefined, true));
    rightPanel.add(new SkillButton(756, 265, 1, "", "OppQ", false));
    rightPanel.add(new SkillButton(862, 315, 2, "", "OppW", false));
    rightPanel.add(new SkillButton(1016, 315, 3, "", "OppE", false));
    rightPanel.add(new SkillButton(1170, 265, 4, "", "OppR", false));
    rightPanel.add(new StandardButton(1196, 5, 5, "Give up", 20, "GiveUp"));
    //bottom panel!!
    let infoEl = new TextInfo(555, 245, red2, "", 20, "info", "info", 170, 180, false);
    bottomPanel.add(infoEl);


    //DO IT
    disableButtons(1);
    disableOppButtons(1);
    document.addEventListener("keydown", e => {
        if ((e.code === 'Space' || e.code === 'Enter') && getElement("back") && getElement("back").clickable) {
            getElement("back").clicked();
            e.preventDefault();
        } else if (e.code === 'Escape' && getElement('GiveUp') && getElement('GiveUp').clickable) {
            getElement('GiveUp').clicked();
            e.preventDefault();
        } else if (!Are_buttons_disabled) {
            if ((e.code === 'KeyQ' || e.code === 'KeyA') && getElement("Q").clickable) {
                getElement("Q").clicked();
                e.preventDefault();
            } else if ((e.code === 'KeyW' || e.code === 'KeyZ') && getElement("W").clickable) {
                getElement("W").clicked();
                e.preventDefault();
            } else if (e.code === 'KeyE' && getElement("E").clickable) {
                getElement("E").clicked();
                e.preventDefault();
            } else if (e.code === 'KeyR' && getElement("R").clickable) {
                getElement("R").clicked();
                e.preventDefault();
            }
        }
    });
    loading = true;
    infoEl.setColour(rightc);
    infoEl.setText("Loading...");
    //let el = getElement("techinfo");
    //el.setText(touch + " " + FullScreen(touch));
    if (TESTING === false) {
        connectToServer();
    } else {
        //parseState(S2);
        //parseState(S3);
        parseState(S);
        disableButtons();
        parseInstruction("Animation:Q", false);
        parseInstruction("Animation:E", true);
        parseInstruction("Animation:W", false);
        parseInstruction("Animation:Q", true);
        parseInstruction("Animation:W", true);
        parseInstruction("Animation:E", false);
        parseInstruction("Animation:R", true);
        parseInstruction("Animation:R", false);
        //S2.OppHP = 32;
        //parseState(S2);
    }
}

function FullScreen(onOff) {
    let doc = window.document;
    let docEl = doc.documentElement;

    let requestFullScreen = docEl.requestFullscreen || docEl.mozRequestFullScreen || docEl.webkitRequestFullScreen || docEl.msRequestFullscreen;
    let cancelFullScreen = doc.exitFullscreen || doc.mozCancelFullScreen || doc.webkitExitFullscreen || doc.msExitFullscreen;

    if (onOff && !doc.fullscreenElement && !doc.mozFullScreenElement && !doc.webkitFullscreenElement && !doc.msFullscreenElement) {
        requestFullScreen.call(docEl);
        return false;
    }
    else {
        //cancelFullScreen.call(doc);
        return true;
    }
}

function clicked() {
    if (touch) return;
    if (document.activeElement.classList.contains("navigation")) {
        return
    }

    let x = mouseX;
    let y = mouseY;

    if (leftPanel.in(x, y)) {
        for (obj of leftPanel.objects) {
            if (obj.clickable && obj.in(x, y)) {
                obj.clicked();
            }
        }
    }
    if (rightPanel.in(x, y)) {
        for (obj of rightPanel.objects) {
            if (obj.clickable && obj.in(x, y)) {
                console.log('CLICKED !!!', obj.id);
                obj.clicked();
            }
        }
    }
    if (bottomPanel.in(x, y)) {
        for (obj of bottomPanel.objects) {
            if (obj.clickable && obj.in(x, y)) {
                obj.clicked();
            }
        }
    }
    if (topPanel.in(x, y)) {
        for (obj of topPanel.objects) {
            if (obj.clickable && obj.in(x, y)) {
                obj.clicked();
            }
        }
    }

}

function touchStarted() {
    if (!touch) return;
    if (!FS) {
        FS = FullScreen(touch);
    }

    let x = mouseX;
    let y = mouseY;

    for (let Panel of [leftPanel, rightPanel, bottomPanel, topPanel]) {
        if (Panel.in(x, y)) {
            for (obj of Panel.objects) {
                if (obj.in(x, y)) {
                    current = obj;
                    return;
                }
            }
        }
    }
}

function touchMoved() {
    if (!touch) return;
    if (document.activeElement.classList.contains("navigation")) {
        return
    }

    let x = mouseX;
    let y = mouseY;

    for (let Panel of [leftPanel, rightPanel, bottomPanel, topPanel]) {
        if (Panel.in(x, y)) {
            for (obj of Panel.objects) {
                if (obj.hoverable && obj.in()) { //found an "in"
                    if (!current) { //outside to something
                        console.log("outside to", obj.id);
                        current = obj;
                        obj.hovered();
                    } else if (current.id !== obj.id) { //switched from another 2 this
                        console.log("switched from", current.id, "to", obj.id);
                        current.unhovered();
                        current = obj;
                        obj.hovered();
                    } else if (current.id === obj.id && !current.isHovered) {
                        current.hovered();
                    }
                    return;
                } else if (obj.hoverable && current && obj.id === current.id) { //went outside
                    console.log("went outside", current.id, "to", obj.id);
                    obj.unhovered();
                    current = undefined;
                    return;
                }
            }
        }
    }

}

function touchEnded() {
    if (!touch) return;
    if (document.activeElement.classList.contains("navigation")) {
        return
    }
    let x = mouseX;
    let y = mouseY;
    if (current) {
        if (current.clickable && current.hoverTimer < current.hoverLinger && current.in(x, y)) {
            console.log(current.id, current.hoverTimer, current.hoverLinger);
            current.clicked();
        }
        current.unhovered();
        current = undefined;
    }
}

function draw() {
//    frameRate(10);
    clear();
    leftPanel.display();
    rightPanel.display();
    topPanel.display();
    bottomPanel.display();
}

function displayTimer(number) {
    let timer = getElement("timer");
    if (number <= 0) {
        timer.setText("");
    } else {
        timer.setText("" + number);
        let c = lerpColor(light, dark, number / 60);
        timer.setColour(c.toString());
    }
}

function parseState(i) {
    if (!i.hasOwnProperty("Instruction")) { //this is a time update
        let num = parseInt((i.split(":")[1]), 10);
        if (!isTicking) {
            isTicking = true;
            countdown(num, redirectwhere, doredirect);
        } else {
            settimeleft(num);
        }
        return;
    } else {
        parseInstruction(i.Instruction, i.TurnPlayer === 1);
    }
    if (i.EndState === 6) {
        disableButtons(0);
        disableOppButtons(0);
        console.log("they dced :<");
        setwhere("/girllist");
        redirect(true);
        bottomPanel.add(new StandardButton(564, 300, 5, "Back to rewards", 20, "back", other));
        if (!isTicking) {
            console.log("COUNTING DOWN DUE TO ENDING", i.EndState);
            isTicking = true;
            countdown(10, redirectwhere, doredirect);
        } else {
            settimeleft(10);
        }
        return;
    }
    //char pics
    if (PICS) {
        setMyChar(i.PlayerName, i.PlayerNum, color(i.SkillColours["Q"]), color(i.SkillColours["W"]), color(i.SkillColours["E"]), color(i.SkillColours["R"]));
        setOppChar(i.OppName, i.OppNum, color(i.OppSkillColours["OppQ"]), i.OppSkillColours["OppW"], color(i.OppSkillColours["OppE"]), color(i.OppSkillColours["OppR"]));
        if (backP === 3) {
            new p5(backSketch, "bacc");
        }
    }
    //turn number and names
    let turn = "Turn " + i.TurnNum;
    getElement("turnNumber").setText(turn + " [" + getTurnNum(i.TurnNum) + "]");
    let name = i.PlayerName + " (" + i.PlayerNum + ")";
    getElement("playerName").setText(name);
    let oppName = i.OppName + " (" + i.OppNum + ")";
    getElement("oppName").setText(oppName);
    //HP and defenses
    //setHP(i.HP, i.MaxHP);
    //setOppHP(i.OppHP, i.OppMaxHP);
    let plHP = getElement("playerHP");
    let oppHP = getElement("oppHP");
    plHP.defenses = i.Defenses;
    oppHP.defenses = i.OppDefenses;


    //Animate HP

    if (plHP.MaxHP && plHP.MaxHP !== i.MaxHP) {
        let speed = calculateHPperFrame(plHP.MaxHP, i.MaxHP);
        plHP.startAnimationMax(speed, i.MaxHP);
    }
    if (plHP.HP && plHP.HP !== i.HP) {
        let speed = calculateHPperFrame(plHP.HP, i.HP);
        plHP.startAnimation(speed, i.HP)
    } else if (!plHP.HP) {
        setHP(plHP, i.HP, i.MaxHP);
        plHP.HP = i.HP;
        plHP.MaxHP = i.MaxHP;
    }


    if (oppHP.MaxHP && oppHP.HP && oppHP.MaxHP !== i.OppMaxHP) {
        let speed = calculateHPperFrame(oppHP.MaxHP, i.OppMaxHP);
        oppHP.startAnimationMax(speed, i.OppMaxHP);
    }
    if (oppHP.HP && oppHP.HP !== i.OppHP) {
        let speed = calculateHPperFrame(oppHP.HP, i.OppHP);
        oppHP.startAnimation(speed, i.OppHP)
    } else if (!oppHP.HP) {
        setOppHP(oppHP, i.OppHP, i.OppMaxHP);
        oppHP.HP = i.OppHP;
        oppHP.MaxHP = i.OppMaxHP;
    }

    //set effects
    IMAGEBOX.clearDisplayed();
    setEffects(i.Effects, "effects");
    setEffects(i.OppEffects, "oppEffects");

    //the order is important. (but I forgot why)
    setSkillNames(i.SkillNames);
    setSkillNames(i.OppSkillNames);
    setSkillColours(i.SkillColours);
    setSkillColours(i.OppSkillColours);
    let info = getElement("info");
    switch (i.EndState) {
        case 0:
            break;
        case 1:
            disableButtons(0);
            disableOppButtons(0);
            rightPanel.discard(getElement("GiveUp"));
            setwhere("/afterbattle");
            redirect(true);
            info.setColour(gc);
            info.setText("★ Victory! ★");
            bottomPanel.add(new StandardButton(564, 300, 5, "Back to rewards", 20, "back", other));
            if (!isTicking) {
                console.log("COUNTING DOWN DUE TO ENDING", i.EndState);
                isTicking = true;
                countdown(10, redirectwhere, doredirect);
            } else {
                settimeleft(10);
            }
            return;
        case 2:
            disableButtons(0);
            disableOppButtons(0);
            rightPanel.discard(getElement("GiveUp"));
            setwhere("/afterbattle");
            redirect(true);
            info.setColour(dark);
            info.setText("Defeat...");
            bottomPanel.add(new StandardButton(564, 300, 5, "Back to rewards", 20, "back", other));
            if (!isTicking) {
                console.log("COUNTING DOWN DUE TO ENDING", i.EndState);
                isTicking = true;
                countdown(10, redirectwhere, doredirect);
            } else {
                settimeleft(10);
            }
            return;
        case 3:
            disableButtons(0);
            disableOppButtons(0);
            rightPanel.discard(getElement("GiveUp"));
            setwhere("/afterbattle");
            redirect(true);
            info.setColour(gc);
            info.setText("Draw! ★");
            bottomPanel.add(new StandardButton(564, 300, 5, "Back to rewards", 20, "back", other));
            if (!isTicking) {
                console.log("COUNTING DOWN DUE TO ENDING", i.EndState);
                isTicking = true;
                countdown(10, redirectwhere, doredirect);
            } else {
                settimeleft(10);
            }
            return;
        case 4:
            disableButtons(0);
            disableOppButtons(0);
            rightPanel.discard(getElement("GiveUp"));
            setwhere("/afterbattle");
            redirect(true);
            info.setColour(light);
            info.setText("Opponent gave up.");
            bottomPanel.add(new StandardButton(564, 300, 5, "Back to rewards", 20, "back", other));
            if (!isTicking) {
                console.log("COUNTING DOWN DUE TO ENDING", i.EndState);
                isTicking = true;
                countdown(10, redirectwhere, doredirect);
            } else {
                settimeleft(10);
            }
            return;
        case 5:
            disableButtons(0);
            disableOppButtons(0);
            rightPanel.discard(getElement("GiveUp"));
            setwhere("/afterbattle");
            redirect(true);
            info.setColour(dark);
            info.setText("Gave up...");
            bottomPanel.add(new StandardButton(564, 300, 5, "Back to rewards", 20, "back", other));
            if (!isTicking) {
                console.log("COUNTING DOWN DUE TO ENDING", i.EndState);
                isTicking = true;
                countdown(10, redirectwhere, doredirect);
            } else {
                settimeleft(10);
            }
            return;
        default:
            return;
    }
    if (i.TurnPlayer === 1) { //it's my turn
        enableButtons(i.SkillState);
    } else { //it's opp's turn
        disableButtons(1);
    }
    enableOppButtons(i.OppSkillState);


}

function parseInstruction(t, isMine) {
    //instruction can be empty; Animation: (used skill, add to turn log), Input (you have typed something wrong)
    // or Error.
    //display Input and Error somewhere else, like below the top panel or in a separate window.
    let errs = getElement("info");
    if (t !== "") {
        let parts = t.split(":");
        switch (parts[0]) {
            case "Animation":
                errs.setText("");
                let used = parts[1][parts[1].length - 1]; //get skill used
                let colour;
                let text;
                if (isMine) { //I should ask myself for the colours and for the names.
                    let skill = getElement(used);
                    colour = color(skill.baseColour.toString());
                    colour.setAlpha(255);
                    text = skill.getText();
                } else { //I am asking my opp.
                    let skill = getElement("Opp" + used);
                    colour = color(skill.baseColour.toString());
                    colour.setAlpha(255);
                    text = skill.getText();
                }
                if (text) {
                    getElement("turnLog").push("[" + used + "] " + text, colour, isMine);
                }
                break;
            case "Input":
                errs.setColour(rightc);
                errs.setText(t.split(":")[1]);
                break;
            case "System":
                errs.setColour(rightc);
                errs.setText(t.split(":")[1]);
                break;
            case "Error":
                errs.setColour(color(red2));
                errs.setText(t.split(":")[1]);
                break;
            default:
                errs.setText("");
                break;
        }
    } else {
        errs.setText("");
        settimeleft(-1);
        displayTimer(-1);
    }
}

function setHP(plHP, HP, MaxHP) {
    let c;
    let inter;
    if (HP > MaxHP / 4) {
        inter = map(HP, MaxHP / 4, MaxHP, 0, 1);
        c = lerpColor(light, dark, inter);
    } else if (HP > 0) {
        inter = map(HP, 0, MaxHP / 4, 0, 1);
        c = lerpColor(color(255, 85, 85), light, inter);
    } else {
        inter = map(HP, 0, -MaxHP / 4, 0, 1);
        c = lerpColor(color(255, 85, 85), color(255, 0, 0), inter);
    }
    plHP.setColour(c.toString());
    let tHP = HP + "/" + MaxHP + " (" + roundUp(HP / MaxHP * 100) + "%)";
    plHP.setText(tHP);
}

function setOppHP(oppHPEl, OppHP, OppMaxHP) {
    let inter;
    let c;
    if (OppHP > OppMaxHP / 4) {
        inter = map(OppHP, OppMaxHP / 4, OppMaxHP, 0, 1);
        c = lerpColor(light, dark, inter);
    } else if (OppHP > 0) {
        inter = map(OppHP, 0, OppMaxHP / 4, 0, 1);
        c = lerpColor(color(255, 85, 85), light, inter);
    } else {
        inter = map(OppHP, 0, -OppMaxHP / 4, 0, 1);
        c = lerpColor(color(255, 85, 85), color(255, 0, 0), inter);
    }
    let toppHP = OppHP + "/" + OppMaxHP + " (" + roundUp(OppHP / OppMaxHP * 100) + "%)";
    oppHPEl.setText(toppHP);
    oppHPEl.setColour(c.toString());

}

function setSkillNames(State) {
    for (let property in State) {
        if (State.hasOwnProperty(property)) {
            getElement(property).setText(State[property]);
        }
    }
}

function setSkillColours(State) {
    console.log("Set skill colours");
    for (let property in State) {
        if (State.hasOwnProperty(property)) {
            getElement(property).setColour(State[property]);
        }
    }
}

function setEffects(effects, id) {
    let effs = sortMap(effects);
    let info = [];
    let info2 = [];
    let len = 0;
    for (let key of effs) {
        if (len >= 5) {
            if (effs[key] === "") {
                info2.push(key + " ");
            } else {
                info2.push(key + " " + effects[key]);
            }
        } else {
            if (effs[key] === "") {
                info.push(key + " ");
            } else {
                info.push(key + " " + effects[key]);
            }
        }
        len += 1;
    }
    getElement(id).setText(info);
    if (len >= 4) {
        getElement(id + "2").setText(info2);
    }
}

function disableButtons(reason) {
    console.log("Disabled buttons");
    Are_buttons_disabled = true;
    getElement("Q").setState(-100);
    getElement("W").setState(-100);
    getElement("E").setState(-100);
    getElement("R").setState(-100);
}

function disableOppButtons(reason) {
    console.log("Disabled opp buttons");
    getElement("OppQ").setState(-100);
    getElement("OppW").setState(-100);
    getElement("OppE").setState(-100);
    getElement("OppR").setState(-100);
}

function enableButtons(State) {
    Are_buttons_disabled = false;
    console.log("Enabled buttons");
    for (let property in State) {
        if (State.hasOwnProperty(property)) {
            getElement(property).setState(State[property]);
        }
    }
}

function enableOppButtons(State) {
    console.log("Enabled opp buttons");
    for (let property in State) {
        if (State.hasOwnProperty(property)) {
            getElement(property).setState(State[property]);
        }
    }
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
        if (loading) {
            getElement("info").setText('');
            loading = false;
        }
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

function sendSkill(Skill) {
    if (connected) {
        if (Are_buttons_disabled && Skill !== "GiveUp") {
            return
        }
        //disableButtons(0);
        Are_buttons_disabled = true;
        ws.send(JSON.stringify(Skill));
    } else {
        let info = getElement("info");
        info.setColour(red2);
        info.setText("You are not connected to the server.");
        loading = true;
        console.log("Not connected!\n");
        connectToServer();
    }
}

function getElement(id) {
    for (let obj of leftPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
    for (let obj of rightPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
    for (let obj of topPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
    for (let obj of bottomPanel.objects) {
        if (obj.id === id) {
            return obj
        }
    }
}

function calculateHPperFrame(HP, targetHP) {
    let totalFrames = FRAMESFORANIMATIONS;
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