function setup() {
    shopobjects = [];
    current = undefined;
    let can = createCanvas(1024, 400);
    can.parent('shop');
    bg_color = color(BG);
    dark = color(DARKC);
    right = color(RIGHTC);
    light = color(LIGHTC);
    clickc = color(CLICKABLEC);
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
    init();
}

function draw() {
    background(bg_color);
    for (let obj of shopobjects) {
        if (obj.clickTimer > 0) {
            obj.clickTimer--;
            if (obj.clickTimer === 0) {
                obj.unclick();
            }
        } else if (obj.hoverable && obj.in()) { //found an "in"
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
        }
            if (!!obj.transparency) {
                tint(255, obj.transparency);
                obj.display();
                tint(255, 255);
            } else {
                obj.display();
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
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/shopitems', true);
    xhr.send(JSON.stringify(ID));
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            console.log(xhr.responseText);
            init();
        }
    };
}

function sendSkill(ID) {
    purchase(ID)
}

function mouseClicked() {
    let x = mouseX;
    let y = mouseY;
    for (obj of shopobjects) {
        if (obj.clickable && obj.in(x, y)) {
            obj.clicked();
        }
    }
}

function keyPressed() {
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