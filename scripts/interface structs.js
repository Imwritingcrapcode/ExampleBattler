function Panel(x, y, w, h, s) {
    this.x = x;
    this.y = y;
    this.width = w;
    this.height = h;
    this.smooth = s;
    this.objects = [];

    this.display = function () {
        for (obj of this.objects) {
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
        //display the rect
        stroke(dark);
        noFill();
        strokeWeight(3);
        rect(this.x, this.y, this.width, this.height, this.smooth);
        //display hover
        if (current && current.in()) {
            current.displayHover();
        }

    };

    this.in = function () {
        let x = mouseX;
        let y = mouseY;
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y))
    };

    this.add = function (obj) {
        this.objects.push(obj)
    }
}

function CanvasImage(x, y, path, id, name, width, height) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.path = path;
    if (!path) {
        this.width = width;
        this.height = height;
        this.image = undefined;
        this.name = "";
    } else {
        this.width = width;
        this.height = height;
        this.image = loadImage(path);
        this.name = name;
    }
    this.hoverable = false;
    this.clickable = false;

    this.getImage = function () {
        return this.image;
    };

    this.getName = function () {
        return this.name;
    };

    this.copy = function () {
        let img = new CanvasImage(0, 0, "", "", "", this.width, this.height);
        img.image = this.image;
        img.path = this.path;
        return img;

    };

    this.open = function (path, name, width, height) {
        this.x = this.x + (550 - width) / 2;
        this.width = width;
        this.height = height;
        this.image = loadImage(path);
        this.name = name;
    };

    this.display = function () {
        if (this.image) {
            image(this.image, this.x, this.y, this.width, this.height);
        }
    }
}

function ImageBox() {
    this.images = [];
    this.isDisplayed = [];

    this.add = function (image) {
        this.images.push(image);
        this.isDisplayed.push(true);
    };

    this.isTaken = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                return this.isDisplayed[i];
            }
        }
        return undefined;
    };

    this.take = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                this.isDisplayed[i] = true;
                return this.images[i];
            }
        }
        return undefined;
    };

    this.clearDisplayed = function () {
        for (let i = 0; i < this.isDisplayed.length; i++) {
            this.isDisplayed[i] = false;
        }
    };

    this.contains = function (name) {
        for (let image of this.images) {
            if (image.name === name) {
                return true;
            }
        }
        return false;
    };

}

function TextInfo(x, y, colour, t, size, id, type, width, height, hoverable) {
    this.visible = true;
    this.id = id;
    this.x = x;
    this.y = y;
    this.textSize = size;
    this.type = type;

    if (this.type === "effects") {
        this.text = [];
        this.images = [];
        this.hoverTexts = [];
        this.hoverHeights = [];
        this.ids = [];
        this.width = width;
        this.height = 175;
    } else if (this.type === "info") {
        this.text = t;
        this.width = width;
        this.height = height;
    } else {
        if (this.id === "playerHP" || this.id === "oppHP") {
            this.HP = undefined;
            this.MaxHP = undefined;
            this.targetHP = undefined;
            this.speed = 0;
            this.framesLeft = 0;
            this.targetMaxHP = undefined;
            this.speedMax = 0;
            this.framesLeftMax = 0;
        }
        this.text = t;
        this.width = size * t.length;
        this.height = size;
    }

    this.textColour = colour;
    this.clickable = false;
    this.hoverable = hoverable;
    if (this.hoverable) {
        this.hoverTimer = 0;
        this.hoverLinger = 180;
    }
    this.hoverText = "";


    this.hide = function () {
        this.visible = false;
        this.wasClickable = this.clickable;
        this.wasHoverable = this.hoverable;
        this.clickable = false;
        this.hoverable = false;
    };

    this.show = function() {
        this.visible = true;
        this.clickable = this.wasClickable;
        this.hoverable = this.wasHoverable;

    };

    this.display = function () {
        if (this.visible) {
            /*let temp_c = color(rightc.toString());
            temp_c.setAlpha(150);
            stroke(temp_c);
            strokeWeight(2);
            temp_c.setAlpha(50);
            fill(temp_c);
            if (this.type !== "info") {
                rect(this.x, this.y-this.height, this.width, this.height);
            } else {
                rect(this.x, this.y, this.width, this.height);
            }*/

            //do the rest
            if (this.type === "effects") {
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                let asc = textAscent();
                let desc = textDescent();
                if (this.id[this.id.length - 1] === "2") {
                    textAlign(RIGHT);
                    for (let i = this.text.length - 1; i >= 0; i--) {
                        this.images[i].display();
                        let y_pos = asc - desc + this.images[i].height / 4;
                        text(this.text[i], this.x + this.width - this.images[i].width - 5, this.y - this.height + this.images[i].height * (i) + y_pos);
                    }
                } else {
                    textAlign(LEFT);
                    for (let i = 0; i < this.text.length; i++) {
                        this.images[i].display();
                        let y_pos = asc - desc + this.images[i].height / 4;
                        text(this.text[i], this.x + this.images[i].width + 5, this.y - this.height + this.images[i].height * (i) + y_pos);
                    }
                }
            }
            else if (this.type === "info") {
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                textAlign(CENTER);
                text(this.text, this.x, this.y, this.width);
            }
            else {
                if ((this.id === "playerHP" || this.id === "oppHP") && this.framesLeft > 0) {
                    this.framesLeft--;
                    if (this.HP + this.speed > this.targetHP && this.speed < 0 ||
                        this.HP + this.speed < this.targetHP && this.speed > 0) { //if it's worth it yet
                        this.HP = this.HP + this.speed;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP)
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                    } else { //end
                        this.HP = this.targetHP;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP);
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                        this.framesLeft = 0;
                    }
                }
                if (((this.id === "playerHP" || this.id === "oppHP") && this.framesLeftMax > 0)) {
                    this.framesLeftMax-=1;
                    if (this.MaxHP + this.speedMax > this.targetMaxHP && this.speedMax < 0 ||
                        this.MaxHP + this.speedMax < this.targetMaxHP && this.speedMax > 0) { //if it's worth it yet
                        this.MaxHP = this.MaxHP + this.speedMax;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP, 1);
                            console.log(this.MaxHP);
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                    } else { //end
                        this.MaxHP = this.targetMaxHP;
                        if (this.id === "playerHP") {
                            setHP(this, this.HP, this.MaxHP);
                        } else {
                            setOppHP(this, this.HP, this.MaxHP)
                        }
                        this.framesLeftMax = 0;
                    }
                }
                stroke(this.textColour);
                strokeWeight(1);
                fill(this.textColour);
                textSize(this.textSize);
                textAlign(LEFT);
                text(this.text, this.x, this.y - 5);
            }
        }
    };

    this.in = function () {
        let x = mouseX;
        let y = mouseY;
        return (this.x <= x && x <= (this.width + this.x) && y <= this.y && (this.y - this.height) <= y);
    };

    this.hovered = function () {
    };

    this.unhovered = function () {
        this.hoverTimer = 0;

    };

    this.setColour = function (c) {
        this.textColour = color(c);
    };

    this.setText = function (t) {
        if (this.type === "info") {
            this.text = t;
        } else if (this.type === "effects") {
            this.text = [];
            this.images = [];
            this.hoverTexts = [];
            this.hoverHeights = [];
            let add = 0;
            let effIconSize = 35;
            for (let i = 0; i < t.length; i++) {
                if (this.id[this.id.length - 1] === "2") {
                    add = this.width - effIconSize;
                }
                let index = t[i].indexOf(" ");
                let items = [];
                items.push(t[i].slice(0, index));
                items.push(t[i].slice(index + 1, t[i].length));
                this.text.push(items[1]);
                //preparing images!
                let name = items[0];
                //and btw also descriptions www
                this.hoverTexts.push(effDescs.get(parseInt(name)));
                this.hoverHeights.push(calculateLines(effDescs.get(parseInt(name))));
                //check if they are already in the box.
                if (!IMAGEBOX.contains(name)) {
                    //if they are not, add them (Adding a canvas image.).
                    let new_image = new CanvasImage(this.x + add, this.y - this.height + effIconSize * i,
                        "/images/locked/" + items[0] + ".png", this.id + "_" + items[0], items[0], effIconSize,
                        effIconSize);
                    IMAGEBOX.add(new_image);
                    this.images.push(new_image);
                } else {
                    //if they are in the box, check if they are already displayed.
                    let image = IMAGEBOX.take(name);
                    if (IMAGEBOX.isTaken(name)) {
                        //if they are already displayed, make a copy with name + '_2' and add it.
                        //for now it's guaranteed there can't be more than 2 identical effect images.
                        let new_image = image.copy();
                        //now we gotta set x, y, id, name.
                        new_image.x = this.x + add;
                        new_image.y = this.y - this.height + effIconSize * i;
                        new_image.name = name + "_2";
                        new_image.id = this.id + "_" + new_image.name;
                        IMAGEBOX.add(new_image);
                        this.images.push(new_image);
                    } else {
                        //if they are not currently displayed, 'take' them.
                        //then change the image so that it would stay where you want it to.
                        image.x = this.x + add;
                        image.y = this.y - this.height + effIconSize * i;
                        this.images.push(image);
                    }
                }
            }
        }
        else {
            this.text = t;
            this.height = this.textSize;
            textSize(this.textSize);
            this.width = textWidth(this.text);
        }
    };

    this.displayHover = function () {
        if (this.hoverTimer < this.hoverLinger) {
            this.hoverTimer += 1;
            return;
        }
        if (this.type === "effects" && this.images.length > 0) {
            this.hoverText = "";
            for (let i = 0; i < 5; i++) {
                if (i >= this.images.length) {
                    return;
                }
                if (Math.floor((mouseY - this.y + this.height) / this.images[0].height) === i) {
                    this.hoverText = this.hoverTexts[i];
                    this.hoverLines = this.hoverHeights[i];
                    break;
                }
            }
        } else if (this.type === "effects") {
            return;
        }
        if ((this.id === "playerHP" || this.id === "oppHP") && this.hasOwnProperty("defenses")) {
            let hoverSize = 15;
            strokeWeight(0.5);
            textSize(hoverSize);
            noStroke();
            fill(hoverc);
            rect(mouseX, mouseY, 355, 70, 5);
            let y_pos = mouseY + hoverSize + 5;
            let x_pos = mouseX + 10;
            let map = new Map(COLOURIDS);
            for (let i = 0; i < COLOURIDS.length; i++) {
                //noStroke();
                textAlign(LEFT);
                let name = COLOURIDS[i][0];
                fill(map.get(name));
                stroke(map.get(name));
                text(name + ":", x_pos, y_pos); //Yellow:
                let amount = this.defenses[String(i + 1)]; //number
                let c;
                if (amount > 0) {
                    c = lerpColor(dark, rightc, amount / 5);
                } else {
                    c = lerpColor(dark, red2, amount / -5);
                }
                fill(c);
                stroke(c);
                textAlign(RIGHT);
                text(amount, x_pos + 75, y_pos);
                if ((i + 1) % 4 === 0) {
                    x_pos = mouseX + 10;
                    y_pos += hoverSize + 5;
                } else {
                    x_pos += 85;
                }
            }
        }
        else if (this.hoverable) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines)
        }
    };

    this.startAnimation = function (speed, target) {
        this.speed = speed;
        this.targetHP = target;
        this.framesLeft = 60;
    };

    this.startAnimationMax = function (speed, target) {
        this.speedMax = speed;
        this.targetMaxHP = target;
        this.framesLeftMax = 30;
    }
}

function TurnLog(x, y, colour, size, id, width, height, hoverable) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.textSize = size;
    this.maxlen = 8;
    this.text = [];
    this.textColours = [];
    this.aligns = [];
    this.turns = [];
    this.len = 0;
    this.width = width;
    this.height = height;
    this.clickable = false;
    this.hoverable = hoverable;
    this.isHovered = false;
    this.hoverText = "";
    this.maxFramesDown = 45;
    this.maxFramesUp = 20;
    this.frame = 0;
    this.isTransitioning = false;
    this.transitioningUp = false;

    this.display = function () {
        strokeWeight(1);
        textSize(this.textSize);
        for (let i = this.len - 1; i >= 0; i--) {
            if (this.isTransitioning) {
                this.frame += 1;
                let c;
                let col;
                if (this.transitioningUp) {
                    col = lerpColor(bg_color, this.textColours[i], 1.05 - (this.len - i - 2) / (this.maxlen - 2));
                    c = lerpColor(col, this.textColours[i], this.frame / this.maxFramesUp);
                } else {
                    col = lerpColor(bg_color, this.textColours[i], 1.05 - (this.len - i - 2) / (this.maxlen - 2));
                    c = lerpColor(this.textColours[i], col, this.frame / this.maxFramesDown);
                }
                stroke(c);
                fill(c);
                if (this.frame > this.maxFramesUp && this.transitioningUp
                    || this.frame > this.maxFramesDown && !this.transitioningUp) {
                    this.isTransitioning = false;
                    this.frame = 0;
                }
            } else if (this.isHovered) {
                stroke(this.textColours[i]);
                fill(this.textColours[i]);
            } else {
                let col = lerpColor(bg_color, this.textColours[i], 1.05 - (this.len - i - 2) / (this.maxlen - 2));
                stroke(col);
                fill(col);
            }
            if (this.aligns[i]) {
                textAlign(LEFT);
                text(this.text[i], this.x + 5, this.y - (this.textSize) * (this.len - i - 1) - 5);
            } else {
                textAlign(RIGHT);
                text(this.text[i], this.x + this.width - 5, this.y - (this.textSize) * (this.len - i - 1) - 5);
            }
        }
    };

    this.push = function (text, colour, isMine) {
        if (this.len < this.maxlen) {
            this.text.push(text);
            //this.textColours.push(colour);
            this.textColours.push(lerpColor(color(dark.toString()), colour, 0.7));
            this.aligns.push(isMine);
            this.len++;
        } else {
            this.pop();
            this.push(text, colour, isMine);
        }
    };

    this.pop = function () {
        let len = this.len;
        if (len > 0) {
            this.text = this.text.slice(1, len);
            this.turns = this.turns.slice(1, len);
            this.textColours = this.textColours.slice(1, len);
            this.aligns = this.aligns.slice(1, len);
            this.len--;
        } else {
            console.log("popping when turn log is empty.")
        }
    };

    this.in = function () {
        let x = mouseX;
        let y = mouseY;
        return (this.x <= x && x <= (this.width + this.x) && y <= this.y && (this.y - this.height) <= y);
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
            this.transitioningUp = true;
        } else {
            this.frame = Math.ceil((this.maxFramesDown - this.frame) / this.maxFramesDown * this.maxFramesUp);
            //this.frame = this.maxframes - this.frame;
        }
    };

    this.unhovered = function () {
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
            this.transitioningUp = false;
        } else {
            this.frame = Math.ceil((this.maxFramesUp - this.frame) / this.maxFramesUp * this.maxFramesDown);
            //this.frame = this.maxframes - this.frame;
        }
    };

    this.displayHover = function () {
        if (this.isHovered && this.hoverText !== "") {
            displayStandardHoverBubble(this.hoverText, calculateLines(this.hoverText));
        }
    }
}

function SkillButton(x, y, type, t, id, mine) {
    this.isMine = mine;
    this.id = id;
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = 50;
    this.textColour = color(dark.toString());
    this.borderColour = color(dark.toString());
    this.borderWidth = 5;
    this.type = type;
    if ((type === 1) || (type === 4)) {
        this.width = 82;
        this.height = 240;
    } else {
        this.width = 132;
        this.height = 210;
    }

    this.maxframes = 7;
    this.frame = 0;
    this.destColour = this.baseColour;
    this.previousColour = this.baseColour;
    this.isTransitioning = false;
    this.colour = clickc;
    this.baseColour = clickc;
    this.hoverColour = clickc;
    this.clickedColour = clickc;
    this.isHovered = false;
    this.hoverText = "";
    this.clickable = mine;
    this.hoverable = true;
    this.clickLinger = 4;
    this.clickTimer = 0;
    if (this.hoverable) {
        this.hoverTimer = 0;
        this.hoverLinger = 180;
    }

    this.display = function () {
        if (this.isTransitioning) {
            if (this.frame <= this.maxframes) {
                this.frame++;
                this.colour = lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
            } else {
                this.frame = 0;
                this.isTransitioning = false;
            }
        }

        let border = this.borderColour;
        let c = this.colour;
        let x = this.x;
        let y = this.y;
        let w = this.width;
        let h = this.height;
        let t = this.text;
        let deg = PI * 0.20872;
        let height = this.textSize;
        let len = this.text.length;
        switch (this.type) {
            case 5:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                rect(this.x, this.y, this.width, this.height, 15, 15, 15, 15);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                /*stroke(color(255, 0, 0));
                line(x, y + h/2, x+w, y + h/2);*/
                for (let i = 0; i < len; i++) {
                    text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
                }
                break;
            case 1:
                //shape
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                beginShape();
                vertex(x + 5, y - 25);
                vertex(x + w - 25, y + 15);
                vertex(x + w - 25, y + 15);
                quadraticVertex(x + w, y + 35, x + w, y + 70);
                vertex(x + w, y + h - 10);
                vertex(x + w, y + h - 10);
                quadraticVertex(x + w, y + h, x + w - 5, y + h - 5);
                vertex(x + 25, y + h - 45);
                vertex(x + 25, y + h - 45);
                quadraticVertex(x, y + h - 65, x, y + h - 100);
                vertex(x, y + h - 100);
                vertex(x, y - 20);
                vertex(x, y - 20);
                quadraticVertex(x, y - 25, x + 5, y - 25);
                endShape(CLOSE);

                //text
                if (this.text !== "") {
                    noStroke();
                    fill(this.textColour);
                    textAlign(CENTER);
                    textSize(this.textSize);
                    for (let i = 0; i < len; i++) {
                        translate(x + w / 2, y + h / 2 - height / 2 + height * (i + 1 - len / 2) / 0.75);
                        rotate(deg);
                        text(t[i], -textWidth("1") / 2, 0);
                        rotate(-deg);
                        translate(-x - w / 2, -y - h / 2 + height / 2 - height * (i + 1 - len / 2) / 0.75);
                    }
                }
                break;
            case 2:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                rect(this.x, this.y, this.width, this.height, 15, 15, 5, 40);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                /*stroke(color(255, 0, 0));
                line(x, y + h/2, x+w, y + h/2);*/
                for (let i = 0; i < len; i++) {
                    text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
                }
                break;
            case 3:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                rect(this.x, this.y, this.width, this.height, 15, 15, 40, 5);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                for (let i = 0; i < len; i++) {
                    text(t[i], x + w / 2, y + h / 2 + height * (i + 1 - len / 2) - 5);
                }
                break;
            case 4:
                stroke(border);
                fill(c);
                strokeWeight(this.borderWidth);
                beginShape();
                vertex(x + w - 5, y - 25);
                vertex(x + 25, y + 15);
                vertex(x + 25, y + 15);
                quadraticVertex(x, y + 35, x, y + 70);
                vertex(x, y + h - 10);
                vertex(x, y + h - 10);
                quadraticVertex(x, y + h, x + 5, y + h - 5);
                vertex(x + w - 25, y + h - 45);
                vertex(x + w - 25, y + h - 45);
                quadraticVertex(x + w, y + h - 65, x + w, y + h - 100);
                vertex(x + w, y + h - 100);
                vertex(x + w, y - 20);
                vertex(x + w, y - 20);
                quadraticVertex(x + w, y - 25, x + w - 5, y - 25);
                endShape(CLOSE);
                noStroke();
                fill(this.textColour);
                textAlign(CENTER);
                textSize(this.textSize);
                deg = PI * 0.20872;
                for (let i = 0; i < len; i++) {
                    translate(x + w / 2, y + h / 2 - height / 2 + height * (i + 1 - len / 2) / 0.75);
                    rotate(-deg);
                    text(t[i], +textWidth("1") / 2, 0);
                    rotate(+deg);
                    translate(-x - w / 2, -y - h / 2 + height / 2 - height * (i + 1 - len / 2) / 0.75);
                }
                break;
            default:
                console.log("wrong skill button type: " + this.id);
                break;
        }
    };

    this.in = function () {
        let x = mouseX;
        let y = mouseY;
        if (this.type === 1) {
            let top_y = 0.77792 * x + (this.y - 25 - (this.x + 5) * 0.77792);

            return (this.x <= x && x <= (this.width + this.x) &&
                y - top_y <= this.height - 35 &&
                0 <= y - top_y);

        } else if (this.type === 4) {
            let top_y = -0.77792 * x + (this.y + 15 + (this.x + 25) * 0.77792);
            return (this.x <= x && x <= (this.width + this.x) &&
                y - top_y <= this.height - 35 &&
                0 <= y - top_y);
        } else {
            return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y));
        }
    };

    this.clicked = function () {
        //console.log("CLICKED " + this.id);
        this.colour = this.clickedColour;
        this.clickTimer = this.clickLinger;
        sendSkill(this.id);
    };

    this.unclick = function () {
        if (this.isHovered) {
            this.colour = this.hoverColour;
        } else {
            this.colour = this.baseColour;
        }
        //console.log("UNCLICKED", this.id);
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;
        }
        this.previousColour = this.colour;
        this.destColour = this.hoverColour;
    };

    this.unhovered = function () {
        this.hoverTimer = 0;
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;

        }
        this.previousColour = this.colour;
        this.destColour = this.baseColour;
    };

    this.setColour = function (stringColour) {
        this.frame = 0;
        this.isTransitioning = false;
        this.baseColour = color(stringColour);
        this.baseColour.setAlpha(0.87 * 255);
        this.colour = this.baseColour;
        this.hoverColour = color(stringColour);
        this.clickedColour = color(stringColour);
        let hoverchange = 17;
        let clickchange = 34;
        if (isLight(this.baseColour)) {
            this.textColour = color(dark.toString());
            this.textColour.setAlpha(0.87 * 255);
            this.hoverColour.setRed(red(this.colour) - hoverchange);
            this.hoverColour.setGreen(green(this.colour) - hoverchange);
            this.hoverColour.setBlue(blue(this.colour) - hoverchange);
            this.clickedColour.setRed(red(this.colour) - clickchange);
            this.clickedColour.setGreen(green(this.colour) - clickchange);
            this.clickedColour.setBlue(blue(this.colour) - clickchange);
            this.hoverColour.setAlpha(0.87 * 255);
            this.clickedColour.setAlpha(0.87 * 255);
        } else {
            this.hoverColour.setRed(red(this.colour) + hoverchange);
            this.hoverColour.setGreen(green(this.colour) + hoverchange);
            this.hoverColour.setBlue(blue(this.colour) + hoverchange);
            this.clickedColour.setRed(red(this.colour) + clickchange);
            this.clickedColour.setGreen(green(this.colour) + clickchange);
            this.clickedColour.setBlue(blue(this.colour) + clickchange);
            this.hoverColour.setAlpha(0.87 * 255);
            this.clickedColour.setAlpha(0.87 * 255);
            /*this.hoverColour.setAlpha(0.87 * 255);
            this.clickedColour.setAlpha(0.91 * 255);*/
            this.textColour = color(light.toString());
            this.textColour.setAlpha(0.87 * 255);
        }

    };

    this.setText = function (t) {
        this.hoverText = SKILLDESCRIPTIONS.get(t);
        if (!!this.hoverText) {
            this.hoverLines = calculateLines(this.hoverText);
        }
        this.text = split(t, " ");
        let w = this.width;
        this.textSize = 50;
        textSize(this.textSize);
        for (let word of this.text) {
            let width = textWidth(word);
            while (width >= w) {
                this.textSize--;
                textSize(this.textSize);
                width = textWidth(word);
            }
        }
    };

    this.getText = function () {
        if (this.text !== "") {
            return this.text.join(" ");
        } else {
            return "";
        }
    };

    this.setState = function (State) {
        this.state = State;
        switch (State) {
            //0 - active, -1 - on CD, -2 - dis by eff ??? -100 - disabled
            case 0:
                this.baseColour.setAlpha(0.87 * 255);
                this.hoverColour.setAlpha(0.87 * 255);
                this.clickedColour.setAlpha(0.87 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.87 * 255);
                this.borderColour = dark;
                this.clickable = this.isMine;
                this.borderWidth = 4.5;
                break;
            case -1:
                this.baseColour.setAlpha(0.5 * 255);
                this.hoverColour.setAlpha(0.5 * 255);
                this.clickedColour.setAlpha(0.5 * 255);
                this.textColour = color(light.toString());
                this.textColour.setAlpha(0.5 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 3;
                break;
            case -2:
                this.baseColour.setAlpha(0.5 * 255);
                this.hoverColour.setAlpha(0.5 * 255);
                this.clickedColour.setAlpha(0.5 * 255);
                this.textColour = color(light.toString());
                this.textColour.setAlpha(0.5 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 3;
                break;
            case -100:
                this.baseColour.setAlpha(0.3 * 255);
                this.hoverColour.setAlpha(0.3 * 255);
                this.clickedColour.setAlpha(0.3 * 255);
                this.textColour = color(light.toString());
                this.textColour.setAlpha(0.3 * 255);
                this.borderColour = light;
                this.clickable = false;
                this.borderWidth = 3;
                break;
            default:
                this.baseColour.setAlpha(0.87 * 255);
                this.hoverColour.setAlpha(0.87 * 255);
                this.clickedColour.setAlpha(0.87 * 255);
                if (isLight(this.baseColour)) {
                    this.textColour = color(dark.toString())
                } else {
                    this.textColour = color(light.toString())
                }
                this.textColour.setAlpha(0.87 * 255);
                this.borderColour = dark;
                this.clickable = this.isMine;
                this.borderWidth = 4.5;
                break;
        }
    };

    this.displayHover = function () {
        if (this.hoverTimer < this.hoverLinger) {
            this.hoverTimer += 1;
            return;
        }
        displayStandardHoverBubble(this.hoverText, this.hoverLines);
    };
}

function StandardButton(x, y, s, t, size, id, col) {
    this.id = id;
    if (this.id === "GiveUp") {
        this.hoverText = "Click here to give up and end the match.";
        this.hoverLines = calculateLines(this.hoverText);
    } else if (this.id === "back") {
        this.hoverText = "Click here to return to your rewards page.";
        this.hoverLines = calculateLines(this.hoverText);
    } else {
        this.hoverText = "";
    }
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = size;
    textSize(size);
    this.width = textWidth(t) + 10;
    this.smooth = s;
    if (!!col) {
        this.baseColour = color(col.toString());
    } else {
        this.baseColour = color(light.toString());
    }
    this.colour = this.baseColour;
    let hoverchange = 20;
    let clickchange = 35;
    this.maxframes = 7;
    this.frame = 0;
    this.destColour = this.baseColour;
    this.previousColour = this.baseColour;
    this.isTransitioning = false;
    this.hoverColour = color(this.colour.toString());
    this.clickedColour = color(this.colour.toString());
    this.hoverColour.setRed(red(this.colour) - hoverchange);
    this.hoverColour.setGreen(green(this.colour) - hoverchange);
    this.hoverColour.setBlue(blue(this.colour) - hoverchange);
    this.clickedColour.setRed(red(this.colour) - clickchange);
    this.clickedColour.setGreen(green(this.colour) - clickchange);
    this.clickedColour.setBlue(blue(this.colour) - clickchange);
    this.textColour = dark;
    this.height = size + 10;
    this.warned = false;
    this.clickable = true;
    this.hoverable = true;
    this.clickLinger = 10;
    this.clickTimer = 0;
    if (this.hoverable) {
        this.hoverTimer = 0;
        this.hoverLinger = 180;
    }
    this.visible = true;

    this.hide = function () {
        this.visible = false;
        this.wasClickable = this.clickable;
        this.wasHoverable = this.hoverable;
        this.clickable = false;
        this.hoverable = false;
    };

    this.show = function() {
        this.visible = true;
        this.clickable = this.wasClickable;
        this.hoverable = this.wasHoverable;

    };

    this.display = function () {
        if (this.visible) {
            if (this.isTransitioning) {
                if (this.frame <= this.maxframes) {
                    this.frame++;
                    this.colour = lerpColor(this.previousColour, this.destColour, (this.frame + 1) / this.maxframes)
                } else {
                    this.frame = 0;
                    this.isTransitioning = false;
                }
            }
            noStroke();
            fill(this.colour);
            strokeWeight(1);
            textSize(this.textSize);
            rect(this.x, this.y, this.width, this.height, this.smooth);
            stroke(this.textColour);
            fill(this.textColour);
            textAlign(LEFT, CENTER);
            text(this.text, this.x + 5, this.y + this.height / 2);
            textAlign(LEFT, BASELINE);
        }
    };

    this.in = function () {
        let x = mouseX;
        let y = mouseY;
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y))
    };

    this.clicked = function () {
        this.colour = this.clickedColour;
        this.clickTimer = this.clickLinger;
        if (this.id === "back") {
            console.log("RELOCATED!~", "/afterbattle");
            window.location = '/afterbattle';
            return
        }
        if (!this.warned && this.id === "GiveUp") {
            this.warned = true;
            let info = getElement("info");
            info.setColour(dark);
            info.setText("Are you sure you want to give up? If so, click this button again.");
        } else if (this.id === "GiveUp") {
            sendSkill(this.id);
        }
    };

    this.unclick = function () {
        if (this.isHovered) {
            this.colour = this.hoverColour;
        } else {
            this.colour = this.baseColour;
        }
    };

    this.setText = function (t) {
        this.text = t;
        textSize(this.textSize);
        this.width = textWidth(t) + 10;
    };

    this.hovered = function () {
        this.isHovered = true;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;
        }
        this.previousColour = this.colour;
        this.destColour = this.hoverColour;
    };

    this.unhovered = function () {
        this.hoverTimer = 0;
        this.isHovered = false;
        if (!this.isTransitioning) {
            this.isTransitioning = true;
        } else {
            this.frame = this.maxframes - this.frame;

        }
        this.previousColour = this.colour;
        this.destColour = this.baseColour;
    };

    this.displayHover = function () {
        if (this.hoverTimer < this.hoverLinger) {
            this.hoverTimer += 1;
        } else if (this.hoverText) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines);
        }
    }
}
