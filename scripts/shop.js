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
    let b1 = new SkillButton(60, 115, 2, "ST", "ST", true);
    b1.setColour(ST.toString());
    b1.setText("ST pack");
    b1.setState(-100);
    let b2 = new SkillButton(252, 115, 2, "AD", "AD", true);
    b2.setColour(AD.toString());
    b2.setText("AD pack");
    b2.setState(-100);
    let b3 = new SkillButton(444, 115, 5, "SP", "SP", true);
    b3.setColour(lerpColor(SP, color("ffffff"), 0.3).toString());
    b3.setText("SP pack");
    b3.setState(-100);
    let b4 = new SkillButton(636, 115, 3, "RP", "RP", true);
    b4.setColour(RP.toString());
    b4.setText("RP pack");
    b4.setState(-100);
    let b5 = new SkillButton(828, 115, 3, "LF", "LF", true);
    b5.setColour(lerpColor(LF, color("ffffff"), 0.2).toString());
    b5.setText("LF pack");
    b5.setState(-100);
    shopobjects.push(b1);
    shopobjects.push(b2);
    shopobjects.push(b3);
    shopobjects.push(b4);
    shopobjects.push(b5);
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
            for (let item of response) {
                let ID = item.ID;
                let el = getElement(ID);
                el.setState(0);
            }
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