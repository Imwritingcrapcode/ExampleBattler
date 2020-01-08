function addFriend(name) {
    console.log("add: " + name);
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/friendlist', true);
    xhr.send(JSON.stringify(["Add", name]));
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                alert(xhr.responseText);
            } else {
                alert(xhr.responseText);
            }
        }
    };
}

function GetInfo() {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/afterbattle', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                let stuff = JSON.parse(xhr.responseText);
                console.log(stuff);
                ParseMatches(stuff);
            } else {
                console.log("???");
            }
        }
    };
}

function ParseMatches(info) {
    switch (info.BattleResult) {
        case 1:
            result = "Victory! ★";
            resColour = win;
            break;
        case 2:
            result = "Defeat...";
            resColour = dark;
            break;
        case 3:
            result = "Draw!";
            resColour = light;
            break;
        case 4:
            result = "Opponent gave up! ★";
            resColour = win;
            break;
        case 5:
            result = "I gave up...";
            resColour = dark;
            break;
        default:
            result = "No new rewards...";
    }
    if (info.LastOpponentsName) {
        oppName = info.LastOpponentsName;
        let size = 20;
        let t = "Add friend";
        strokeWeight(1);
        textSize(size);
        let w = textWidth(t);
        objects.push(new StandardButton(x + initial_w - w - 10, y + 125 + size, 5, t, size, info.LastOpponentsName));
        getElement(info.LastOpponentsName).clicked = function () {
            this.colour = this.clickedColour;
            this.clickTimer = this.clickLinger;
            addFriend(this.id);
        };
    }
    globalinfo = info;
    matches = info.Matches;
    to_add = info.ToAdd;
    let to_set_matches = info.TotalMatches - info.ToAdd;
    level = 1;
    passed = 0;
    while ((to_set_matches >= matches[level - 1]) && level < 20) {
        to_set_matches = to_set_matches - matches[level - 1];
        passed += matches[level - 1];
        level += 1;
    }
    curr_matches = to_set_matches;
    percentage = 100 * curr_matches / matches[level - 1];
    if ("w" in info.Dusts) {
        console.log("hasDust!");
        dust = info.Dusts["w"];
    } else {
        console.log("no Dust!");
    }
    gname = info.Name;
}

function setup() {
    can = createCanvas(600, 350);
    can.parent('rewards_sketch');
    bg_color = color(BG);
    dark = color(DARKC);
    light = color(LIGHTC);
    right = color(RIGHTC);
    win = color(WINC);
    rectangle = dark;
    matches = [];
    level = 1;
    curr_matches = 0;
    to_add = 0;
    dust = 0;
    added = 0;
    levelled_up = false;
    percentage = 0.0;
    new_percentage = 0.0;
    speed = to_add / 4;
    gname = "";
    oppName = "";
    result = "";
    resColour = dark;
    current = undefined;
    objects = [];
    x = 50;
    y = 100;
    h = 50;
    initial_w = 500;
    r = 10;
    let size = 30;
    let t = "Battle again";
    strokeWeight(1);
    textSize(size);
    let w = textWidth(t);
    objects.push(new StandardButton(x + (initial_w - w - 10) / 2, y + 195, 5, t, size, "girlList"));
    getElement("girlList").clicked = function () {
        this.colour = this.clickedColour;
        this.clickTimer = this.clickLinger;
        window.location = "/girllist";
    };
    GetInfo();
}

function draw() {
    background(bg_color);
    let stop = lerpColor(dark, right, percentage / 100);
    drawText(stop);
    w = round(initial_w * percentage / 100);
    setGradient(x, y, w, h, dark, stop);
    drawRect(x, y, initial_w, h);
    if (new_percentage > percentage) {
        percentage = percentage + speed;
        if (percentage > 100) {
            percentage = 100;
            new_percentage = 100;
        }
    } else if (percentage >= 100 && level < 20) {
        level = level + 1;
        levelled_up = true;
        if (level < 20) {
            percentage = 0.0;
            new_percentage = 0.0;
        }
    } else if (to_add > added) {
        add_match();
    }

    for (let obj of objects) {
        if (obj.clickable && obj.clickTimer > 0) {
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

function drawRect(x, y, w, h) {
    stroke(rectangle);
    strokeWeight(1);
    noFill();
    rect(x, y, w, h, r);
}

function drawText(stop) {
    strokeWeight(1);
    fill(stop);
    noStroke();
    textSize(30);
    textAlign(LEFT);
    text("Lv. " + level, x, y - 10);
    textAlign(CENTER);
    text(gname, x + 200, y - 10);
    textAlign(RIGHT);
    if (level === 20) {
        percentage = 100;
        new_percentage = 0;
        fill(255, 204, 102);
        text("★ Max ★", (initial_w - x) / 2 + 250, y - 10);
        right = color(130, 190, 255)
    } else {
        text(round(percentage) + "%", (initial_w - x) / 2 + 310, y - 10);
        if (levelled_up) {
            fill(win);
            text("Lvl up! ★", (initial_w - x) / 2 + 250, y - 10);
        }
    }
    textAlign(LEFT);
    fill(dark);
    text("+ Dust: " + dust, x, y + 100);
    if (oppName) {
        textAlign(RIGHT);
        textSize(25);
        text("Opponent: " + oppName, x + initial_w, y + 130);
    }
    textAlign(CENTER);
    textSize(40);
    fill(resColour);
    stroke(resColour);
    text(result, 300, 40);
}

function setGradient(x, y, w, h, c1, c2) {
    noFill();
    strokeWeight(2);
    //circle at the beginning
    for (let i = x; (i < x + r) && (i < x + w); i++) {
        let inter = map(i, x, x + w, 0, 1);
        let c = lerpColor(c1, c2, inter);
        stroke(c);
        let top_y = y + r - sqrt(sq(r) - sq(i - x - r));
        let bot_y = y + h - r + sqrt(sq(r) - sq(i - x - r));
        line(i, top_y, i, bot_y);
    }

    for (let i = x + r; i < x + w - r; i++) {
        let inter = map(i, x, x + w, 0, 1);
        let c = lerpColor(c1, c2, inter);
        stroke(c);
        line(i, y, i, y + h);
    }
    //circle at the end
    for (let i = x + w; (i >= x + r) && (i >= x + w - r); i--) {
        let inter = map(i, x, x + w, 0, 1);
        let c = lerpColor(c1, c2, inter);
        stroke(c);
        let top_y = y + r - sqrt(sq(r) - sq(i - x - w + r));
        let bot_y = y + h - r + sqrt(sq(r) - sq(i - x - w + r));
        line(i, top_y, i, bot_y);
    }
}

function add_match() {
    curr_matches = curr_matches + 1;
    added += 1;
    let to_level_up = matches[level - 1];
    if (to_add > to_level_up && to_add > 1) {
        speed = map(added, 0, to_add, to_add / 4, 0.5);
    } else {
        speed = 0.5
    }
    new_percentage = 100 * curr_matches / to_level_up;
    if (to_level_up === curr_matches) {
        curr_matches = 0
    }
}

function mousePressed() {
    let x = mouseX;
    let y = mouseY;
    for (obj of objects) {
        if (obj.clickable && obj.in(x, y)) {
            obj.clicked();
        }
    }
}

function getElement(id) {
    for (obj of objects) {
        if (obj.id === id) {
            return obj
        }
    }
}

function keyPressed() {
    if (key === ' ') {
        getElement("girlList").clicked()
    }
}