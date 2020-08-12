function setup() {
    shopobjects = [];
    STATE = "NO_WINDOW";
    current = undefined;
    touch = is_touch_device4();
    let can = createCanvas(1024, 490);
    can.parent('shop');
    bg_color = color(BG);
    dark = color(DARKC);
    right = color(RIGHTC);
    light = color(LIGHTC);
    clickc = color(CLICKABLEC);
    hoverc = color(HOVERC);
    ST = color(light);
    AD = color(ADCOLOUR);
    SP = color(SPCOLOUR);
    RP = color(RPCOLOUR);
    LF = color(LFCOLOUR);
    let b1 = new SkillButton(60, 6, 2, "ST", "ST", true);
    b1.setColour(ST.toString());
    b1.setText("ST pack");
    b1.setState(-100);
    let pic1 = new CanvasImage(65, 231, "", "STimage", "STimage", 40, 40);
    let text1 = new TextInfo(115, 271, light, "???", 40, "STtext");
    let b2 = new SkillButton(252, 5, 2, "AD", "AD", true);
    b2.setColour(AD.toString());
    b2.setText("AD pack");
    b2.setState(-100);
    let pic2 = new CanvasImage(257, 231, "", "ADimage", "ADimage", 40, 40);
    let text2 = new TextInfo(307, 271, light, "???", 40, "ADtext");
    let b3 = new SkillButton(444, 6, 5, "SP", "SP", true);
    b3.setColour(lerpColor(SP, color("ffffff"), 0.3).toString());
    b3.setText("SP pack");
    b3.setState(-100);
    let pic3 = new CanvasImage(449, 231, "", "SPimage", "SPimage", 40, 40);
    let text3 = new TextInfo(499, 271, light, "???", 40, "SPtext");
    let b4 = new SkillButton(636, 6, 3, "RP", "RP", true);
    b4.setColour(RP.toString());
    b4.setText("RP pack");
    b4.setState(-100);
    let pic4 = new CanvasImage(641, 231, "", "RPimage", "RPimage", 40, 40);
    let text4 = new TextInfo(691, 271, light, "???", 40, "RPtext");
    let b5 = new SkillButton(828, 6, 3, "LF", "LF", true);
    b5.setColour(lerpColor(LF, color("ffffff"), 0.2).toString());
    b5.setText("LF pack");
    b5.setState(-100);
    let pic5 = new CanvasImage(833, 226, "", "LFimage", "LFimage", 40, 40);
    let text5 = new TextInfo(883, 266, light, "???", 40, "LFtext");
    shopobjects.push(b1);
    shopobjects.push(text1);
    shopobjects.push(pic1);
    shopobjects.push(b2);
    shopobjects.push(text2);
    shopobjects.push(pic2);
    shopobjects.push(b3);
    shopobjects.push(text3);
    shopobjects.push(pic3);
    shopobjects.push(b4);
    shopobjects.push(text4);
    shopobjects.push(pic4);
    shopobjects.push(b5);
    shopobjects.push(text5);
    shopobjects.push(pic5);
    imageBox = new ImageBox();
    init();
}

function open() {
    left_pos = 136.5;
    right_pos = 686.5;
    white_window = {
        hoverable: false,
        clickable: false,
        id: 'white_window',
        x: 121.5,
        y: 0,
        width: 781,
        height: 490,
        display: function () {
            noStroke();
            let c = color(hoverc.toString());
            c.setAlpha(255 * 0.75);
            fill(c);
            rect(this.x, this.y, this.width, this.height, 5);
        },
        in: function () {
            let x = mouseX;
            let y = mouseY;
            return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.y + this.height));
        }
    };
    let t1 = new TextInfo(left_pos, 45, dark, "Loading...", 35, "name", "A");
    let t2 = new TextInfo(left_pos, 71, dark, "Rarity", 21, "rarity", "A");
    let r1 = {
        hoverable: false,
        clickable: false,
        id: 'r1',
        x: right_pos,
        y: 15,
        width: 201,
        height: 268,
        display: function () {
            fill(light);
            noStroke();
            rect(this.x, this.y, this.width, this.height, 5);
        }
    };
    let i1 = new CanvasImage(right_pos, 15, "", "girl", "", 201, 268, dark);
    p_screen = new LoadingScreenNoP(0.4 * 201 + right_pos, 15 + 0.4 * 268, 0.2 * 201, 0.2 * 201);
    p_screen.setColours(dark, bg_color, light, right);
    let r2 = {
        hoverable: false,
        clickable: false,
        id: 'r2',
        x: right_pos,
        y: 15,
        width: 201,
        height: 268,
        display: function () {
            stroke(getElement("girl").colour);
            strokeWeight(4);
            noFill();
            rect(this.x, this.y, this.width, this.height, 5);
        }
    };
    let t6 = new TextInfo(left_pos, 127, dark, "Tags", 21, "tags", "A");
    let t8 = new TextInfo(left_pos, 160, dark, "Skill colours", 21, "skillcolours", "A");
    let sk1 = new SkillButtonMini(left_pos, 170, "", "SkillQ");
    sk1.setText('Q');
    sk1.setColour(light.toString());
    let sk2 = new SkillButtonMini(left_pos + 110, 170, "", "SkillW");
    sk2.setText('W');
    sk2.setColour(light.toString());
    let sk3 = new SkillButtonMini(left_pos + 220, 170, "", "SkillE");
    sk3.setText('E');
    sk3.setColour(light.toString());
    let sk4 = new SkillButtonMini(left_pos + 330, 170, "", "SkillR");
    sk4.setText('R');
    sk4.setColour(light.toString());
    /*let sk5 = new SkillButtonMini(left_pos + 440, 170, "", "SkillD");
    sk5.setText('Your Number');
    sk5.setColour(light.toString());*/

    let t9 = new TextInfo(left_pos, 276, dark, "Description", 21, "description", "C", 560);
    let c = new StandardButton(475, 442, 5, "Close", 28, 'close');
    shopobjects.push(white_window);
    shopobjects.push(t1);
    shopobjects.push(t2);
    shopobjects.push(r1);
    shopobjects.push(i1);
    shopobjects.push(r2);
    shopobjects.push(t6);
    shopobjects.push(t8);
    shopobjects.push(t9);
    shopobjects.push(c);
    shopobjects.push(sk1);
    shopobjects.push(sk2);
    shopobjects.push(sk3);
    shopobjects.push(sk4);
    //shopobjects.push(sk5);
}

function removeElement(id) {
    for (let i = 0; i < shopobjects.length; i++) {
        let obj = shopobjects[i];
        if (obj.id === id) {
            shopobjects.splice(i, 1);
            return;
        }
    }
}

function close() {
    STATE = 'NO_WINDOW';
    removeElement('white_window');
    removeElement('name');
    removeElement('rarity');
    removeElement('r1');
    removeElement('girl');
    removeElement('r2');
    removeElement('tags');
    removeElement('skills');
    removeElement('skillcolours');
    removeElement('description');
    removeElement('close');
    removeElement('SkillQ');
    removeElement('SkillW');
    removeElement('SkillE');
    removeElement('SkillR');
    //removeElement('SkillD');
}

function draw() {
    background(bg_color);
    for (let obj of shopobjects) {
        if (obj.clickTimer > 0) {
            obj.clickTimer--;
            if (obj.clickTimer === 0) {
                obj.unclick();
            }
        } else if (obj.hoverable && obj.in() && ((STATE === 'NO_WINDOW') || (obj.id === 'close')|| (obj.id === 'SkillD')|| (obj.id === 'SkillR') || (obj.id === 'SkillE')|| (obj.id === 'SkillW') || (obj.id === 'SkillQ') || (!white_window.in()))) { //found an "in"
            if (!current) { //outside to something
                current = obj;
                obj.hovered();
            } else if (current.id !== obj.id) { //switched from another 2 this
                current.unhovered();
                current = obj;
                obj.hovered();
            }
        } else if (obj.hoverable && current && obj.id === current.id) { //went outside
            obj.unhovered();
            current = undefined;
        } else if (obj.id === "girl" && obj.loaded() && p_screen.stopped < 1) {
            p_screen.stop();
        } else if (obj.id === "girl" && obj.loaded()) {
            obj.display();
        } else if (obj.id === "girl" && !obj.loaded() && p_screen.stopped < 1) {
            p_screen.display();
        }
        if (!!obj.transparency) {
            tint(255, obj.transparency);
            obj.display();
            tint(255, 255);
        } else {
            obj.display();
        }
        if (!touch && current && current.in() || touch && current && current.isHovered && current.in()) {
            current.displayHover();
        }
    }
}

function getElement(id) {
    for (obj of shopobjects) {
        if (obj.id === id) {
            return obj
        }
    }
}

function init() {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/shopitems', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            let response = JSON.parse(xhr.responseText);
            console.log(response);

            function after(data) {
                for (let item of response) {
                    let ID = item.ID;
                    let cost = item.Cost;
                    let letter = item.Dust;
                    let el = getElement(ID);
                    let im = getElement(ID + "image");
                    let txt = getElement(ID + "text");
                    if (cost <= int(data.MoneyInfo[letter])) {
                        el.setState(0);
                        txt.setColour(dark);
                        im.transparency = 255;
                    } else {
                        el.setState(-1);
                        txt.setColour(light);
                        im.transparency = 126;
                    }
                    txt.setText(cost);
                    im.image = loadImage("/images/locked/" + DUSTS.get(item.Dust) + "_dust_small.png");
                }
            }

            UpdateFreeData(after);
        }
    };
}

function purchase(ID) {
    if (STATE !== "WINDOW") {
        STATE = "WINDOW";
        let xhr = new XMLHttpRequest();
        xhr.open('POST', '/shopitems', true);
        xhr.send(JSON.stringify(ID));
        xhr.onreadystatechange = (e) => {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    response = JSON.parse(xhr.responseText);
                    console.log(response);
                    parse(response);
                    init();
                } else {
                    addPopup(xhr.responseText);
                }
            }
        };
    }
}

function parse(curGirl) {
    open();
    getElement("name").setText(curGirl.Name + " (" + curGirl.Number + ")");
    getElement("rarity").setText(curGirl.Rarity);
    switch (curGirl.Rarity) {
        case "LF":
            //star
            getElement("rarity").setColour(LFCOLOUR);
            break;
        case "RP":
            //green
            getElement("rarity").setColour(RPCOLOUR);
            break;
        case "SP":
            //yellow
            getElement("rarity").setColour(SPCOLOUR);
            break;
        case "AD":
            //blue
            getElement("rarity").setColour(ADCOLOUR);
            break;
        case "ST":
            //white
            getElement("rarity").setColour(STCOLOUR);
            break;
        default:
            console.log("EXCUSE ME");
            break;
    }
    let name = curGirl.Name;
    if (!imageBox.contains(name)) {
        if (existsPortrait(curGirl.Number)) {
            imageBox.add(new CanvasImage(right_pos, 15, "/images/locked/" + name + "_portrait.png", "girl", name, 201, 268, curGirl.MainColour));
        } else {
            imageBox.add(new CanvasImage(right_pos, 15, "/images/locked/Placeholder_portrait.png", "girl", name, 201, 268, curGirl.MainColour));
        }
    }
    //setloadingscreencolours;
    p_screen.setColours(curGirl.SkillColourCodes[0], curGirl.SkillColourCodes[1], curGirl.SkillColourCodes[2], curGirl.SkillColourCodes[3]);
    // update the pic
    getElement("girl").set(imageBox.get(name));
    //p_screen.restart()
    p_screen.restart();
    let tags = "";
    for (let tag of curGirl.Tags.sort()) {
        tags += tag + ", "
    }
    tags = tags.slice(0, tags.length - 2);
    getElement("tags").setText(tags);
    const letters = ['Q', 'W', 'E', 'R', 'D'];
    for (let i = 0; i < curGirl.Skills.length; i++) {
        let name = curGirl.Skills[i];
        let letter = letters[i];
        let el = getElement('Skill'+letter);
        el.setText(name);
        el.setColour(curGirl.SkillColourCodes[i]);
    }
    let colours = "";
    for (let s of curGirl.SkillColours) {
        colours += s + ", "
    }
    colours = colours.slice(0, colours.length - 2);
    getElement("skillcolours").setText(colours);
    getElement("description").setText(curGirl.Description);
}

function sendSkill(ID) {
    purchase(ID)
}

function mouseClicked() {
    let x = mouseX;
    let y = mouseY;
    for (obj of shopobjects) {
        if (obj.clickable && obj.in(x, y) && STATE === "NO_WINDOW") {
            obj.clicked();
        } else if (obj.clickable && obj.in(x, y) && STATE === "WINDOW" && obj.id === 'close') {
            obj.clicked();
            close();
        }
    }
}

function keyPressed() {
    if (STATE === 'NO_WINDOW') {
        if ((key === 'q' || key === 'a') && getElement("ST").clickable) {
            getElement("ST").clicked()
        } else if ((key === 'w' || key === 'z') && getElement("AD").clickable) {
            getElement("AD").clicked()
        } else if (key === 'e' && getElement("SP").clickable) {
            getElement("SP").clicked()
        } else if (key === 'r' && getElement("RP").clickable) {
            getElement("RP").clicked()
        } else if (key === 't' && getElement("LF").clickable) {
            getElement("LF").clicked()
        }
    }

}