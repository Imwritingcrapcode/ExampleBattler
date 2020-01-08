function InterfaceButton(p, x, y, t, size, id, type, width, height) {
    this.p = p;
    this.id = id;
    this.x = x;
    this.y = y;
    this.text = t;
    this.textSize = size;
    this.smooth = 5;
    this.type = type;
    if (this.type === "A") {
        this.height = size + 10;
        this.p.textSize(size);
        this.width = this.p.textWidth(t) + 10;
    } else if (this.type === "B") {
        this.width = width;
        while (this.p.textWidth(t) > (this.width - 10)) {
            this.textSize--;
            this.p.textSize(this.textSize);
        }

        if (height) {
            this.height = height;
        } else {
            this.height = 70;
        }

    }
    this.textColour = dark;
    this.clickable = true;
    this.hoverable = true;
    this.clickLinger = 4;
    this.clickTimer = 0;
    this.baseColour = this.p.color(light.toString());
    this.colour = this.baseColour;
    let hoverchange = 17;
    let clickchange = 34;
    this.hoverColour = this.p.color(this.colour.toString());
    this.clickedColour = this.p.color(this.colour.toString());
    this.hoverColour.setRed(this.p.red(this.colour) - hoverchange);
    this.hoverColour.setGreen(this.p.green(this.colour) - hoverchange);
    this.hoverColour.setBlue(this.p.blue(this.colour) - hoverchange);
    this.clickedColour.setRed(this.p.red(this.colour) - clickchange);
    this.clickedColour.setGreen(this.p.green(this.colour) - clickchange);
    this.clickedColour.setBlue(this.p.blue(this.colour) - clickchange);
    /*this.hoverColour.setAlpha(0.87 * 255);
    this.clickedColour.setAlpha(0.87 * 255);
    this.baseColour.setAlpha(0.87 * 255);*/

    this.setColour = function (colour) {
        this.baseColour = this.p.color(colour);
        let hoverchange = 17;
        let clickchange = 34;
        this.hoverColour = this.p.color(this.baseColour.toString());
        this.clickedColour = this.p.color(this.baseColour.toString());
        this.hoverColour.setRed(this.p.red(this.baseColour) - hoverchange);
        this.hoverColour.setGreen(this.p.green(this.baseColour) - hoverchange);
        this.hoverColour.setBlue(this.p.blue(this.baseColour) - hoverchange);
        this.clickedColour.setRed(this.p.red(this.baseColour) - clickchange);
        this.clickedColour.setGreen(this.p.green(this.baseColour) - clickchange);
        this.clickedColour.setBlue(this.p.blue(this.baseColour) - clickchange);
        /*this.hoverColour.setAlpha(0.87 * 255);
        this.clickedColour.setAlpha(0.87 * 255);
        this.baseColour.setAlpha(0.87 * 255);*/
        this.colour = this.baseColour;
    };

    this.display = function () {
        this.p.noStroke();
        this.p.fill(this.colour);
        this.p.strokeWeight(1);
        this.p.textSize(this.textSize);
        this.p.rect(this.x, this.y, this.width, this.height, this.smooth);
        if (this.clickable) {
            this.textColour = dark;
        } else {
            this.textColour = this.p.color(110);
        }
        this.p.stroke(this.textColour);
        this.p.fill(this.textColour);
        if (this.type === "A") {
            this.p.textAlign(this.p.LEFT, this.p.CENTER);
            this.p.text(this.text, this.x + 5, this.y + this.height / 2);
        } else if (this.type === "B") {
            this.p.textAlign(this.p.CENTER, this.p.CENTER);
            this.p.text(this.text, this.x + this.width / 2, this.y + this.height / 2);

        }
        this.p.textAlign(this.p.LEFT, this.p.BASELINE);
    };

    this.in = function () {
        let x = this.p.mouseX;
        let y = this.p.mouseY;
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.height + this.y))
    };

    this.clicked = function () {
        if (this.clickable) {
            this.colour = this.clickedColour;
            this.clickTimer = this.clickLinger;
            this.onClick();
        }
    };

    this.onClick = function () {

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
        this.p.textSize(this.textSize);
        if (this.type === "A") {
            this.width = this.p.textWidth(t);
        } else if (this.type === "B") {
            this.textSize = 45;
            while (this.p.textWidth(t) > (this.width - 10)) {
                this.textSize--;
                this.p.textSize(this.textSize);
            }
        }
    };
    this.hovered = function () {
        this.isHovered = true;
        this.colour = this.hoverColour;
    };

    this.unhovered = function () {
        this.isHovered = false;
        this.colour = this.baseColour;
    };

    this.displayHover = function () {
    }
}

function InterfaceText(p, x, y, colour, t, size, id, type, width) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.textSize = size;
    this.text = t;
    this.type = type;
    if (type === "B" || type === "C") {
        this.width = width;
    } else if (type === "A" || type === "D") {
        p.textSize(size);
        this.width = p.textWidth(t);
    }

    this.height = size;
    this.p = p;
    this.textColour = colour;
    this.hoverable = false;

    this.display = function () {
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
        this.p.stroke(this.textColour);
        this.p.strokeWeight(1);
        this.p.fill(this.textColour);
        this.p.textSize(this.textSize);
        if (this.type === "B" || this.type === "C") {
            if (this.type === "B") {
                this.p.textAlign(this.p.CENTER);
            } else {
                this.p.textAlign(this.p.LEFT);
                this.p.textLeading(33);
            }
            this.p.text(this.text, this.x, this.y, this.width);
        } else if (this.type === "A" || this.type === "D") {
            if (this.type === "A") {
                this.p.textAlign(this.p.LEFT);
            } else if (this.type === "D") {
                this.p.textAlign(this.p.RIGHT);
            }
            this.p.text(this.text, this.x, this.y);
        }
    };

    this.in = function () {
        let x = this.p.mouseX;
        let y = this.p.mouseY;
        return (this.x <= x && x <= (this.width + this.x) && y <= this.y && (this.y - this.height) <= y);
    };

    this.setColour = function (c) {
        this.textColour = this.p.color(c);
    };

    this.setText = function (t) {
        this.text = t;
        this.height = this.textSize;
        this.p.textSize(this.textSize);
        if (type === "B" || type === "C") {
            this.width = width;
        } else if (type === "A") {
            p.textSize(size);
            this.width = p.textWidth(t);
        }
    };

    this.displayHover = function () {
        if (this.hoverable) {
            displayStandardHoverBubble(this.hoverText, this.hoverLines)
        }
    };
}

function InterfaceImage(p, x, y, path, id, name, width, height, colour) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.path = path;
    if (!!colour) {
        this.colour = p.lerpColor(leftP.color(colour), light, 0.45);
    } else {
        this.colour = p.color(DARKC);
    }
    if (!path) {
        this.width = width;
        this.height = height;
        this.image = undefined;
        this.name = "";
    } else {
        this.width = width;
        this.height = height;
        this.image = p.loadImage(path);
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
        let img = new InterfaceImage(p, 0, 0, "", "", "", this.width, this.height);
        img.image = this.image;
        img.path = this.path;
        img.colour = this.colour;
        return img;

    };

    this.set = function (other) {
        this.id = other.id;
        this.x = other.x;
        this.y = other.y;
        this.path = other.path;
        this.colour = other.colour;
        if (!other.path) {
            this.width = other.width;
            this.height = other.height;
            this.image = undefined;
            this.name = other.name;
        } else {
            this.width = other.width;
            this.height = other.height;
            this.image = other.image;
            this.name = other.name;
        }
        this.hoverable = other.hoverable;
        this.clickable = other.clickable;
    };

    this.open = function (path, name, width, height) {
        this.x = this.x + (550 - width) / 2;
        this.width = width;
        this.height = height;
        this.image = p.loadImage(path);
        this.name = name;
    };

    this.display = function () {
        if (this.image) {
            p.image(this.image, this.x, this.y, this.width, this.height);
        }
    }
}

function InterfaceImageBox() {
    this.images = [];

    this.add = function (image) {
        this.images.push(image);
    };

    this.contains = function (name) {
        for (let image of this.images) {
            if (image.name === name) {
                return true;
            }
        }
        return false;
    };

    this.get = function (name) {
        for (let i = 0; i < this.images.length; i++) {
            if (this.images[i].name === name) {
                return this.images[i];
            }
        }
        return undefined;
    };


}

function interfaceCalculateLines(p, hoverText, width, size) {
    if (!width) {
        width = 290;
    }
    if (!size) {
        size = 15;
    }
    let height = 0;
    p.textSize(size);
    for (let line of hoverText.split("\n")) {
        height += 1;
        let x_pos = 0;
        for (let word of line.split(" ")) {
            /*strokeWeight(5);
            stroke(red2);
            point(mouseX+10+x_pos + change, mouseY+height*size);*/
            if (x_pos + p.textWidth(word + " ") < width) { //we are still on that line
                x_pos += p.textWidth(word + " ")
            } else { //start a new line
                height += 1;
                x_pos = p.textWidth(word + " ");
            }
        }
    }
    return height;
}