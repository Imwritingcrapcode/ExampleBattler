rightSketch = function (p) {
    rightP = p;
    p.setup = function () {
        imageBox = new InterfaceImageBox();
        rightobjects = [];
        can3 = p.createCanvas(768, 413);
        can3.parent('rightcanvas');
        bg_color2 = p.color(BG);
        let t1 = new InterfaceText(p, 5, 35, dark, "", 35, "name", "A");
        let t2 = new InterfaceText(p, 5, 61, dark, "", 21, "rarity", "A");
        let t3 = new InterfaceText(p, 5, 131, dark, "", 28, "level", "A");
        let t4 = new InterfaceText(p, 5, 164, dark, "", 28, "matchesPlayed", "A");
        let t5 = new InterfaceText(p, 5, 197, dark, "", 28, "matchesWon", "A");
        let r1 = {
            hoverable: false,
            clickable: false,
            x: 519.5,
            y: 5,
            width: 201,
            height: 268,
            display: function () {
                p.fill(light);
                p.noStroke();
                p.rect(this.x, this.y, this.width, this.height, 5);
            }
        };
        let i1 = new InterfaceImage(p, 519.5, 5, "", "girl", "", 201, 268, dark);
        p_screen = new LoadingScreen(p, 0.4*201 + 519.5, 5 + 0.4*268, 0.2*201, 0.2*201);
        p_screen.setColours(dark, bg_color2, other, right);
        let r2 = {
            hoverable: false,
            clickable: false,
            x: 519.5,
            y: 5,
            width: 201,
            height: 268,
            display: function () {
                p.stroke(getElementRight("girl").colour);
                p.strokeWeight(4);
                p.noFill();
                p.rect(this.x, this.y, this.width, this.height, 5);
            }
        };
        let t6 = new InterfaceText(p, 5, 267, dark, "", 28, "tags", "A");
        let t7 = new InterfaceText(p, 5, 300, dark, "", 28, "skills", "A");
        let t8 = new InterfaceText(p, 5, 333, dark, "", 28, "skillcolours", "A");
        let t9 = new InterfaceText(p, 5, 380, dark, "", 28, "description", "C", 743);
        let t10 = new InterfaceText(p, 514.5, 61, dark, "", 28, "set", "D");
        rightobjects.push(t1);
        rightobjects.push(t2);
        rightobjects.push(t3);
        rightobjects.push(t4);
        rightobjects.push(t5);
        rightobjects.push(r1);
        rightobjects.push(i1);
        rightobjects.push(r2);
        rightobjects.push(t6);
        rightobjects.push(t7);
        rightobjects.push(t8);
        rightobjects.push(t9);
        rightobjects.push(t10);
    };
    p.draw = function () {
        p.background(bg_color2);
        for (let obj of rightobjects) {
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
            } else if (obj.id === "girl" && obj.loaded() && p_screen.stopped < 1) {
                p_screen.stop();
                p_screen.clear();
            } else if (obj.id === "girl" && obj.loaded()) {
                obj.display();
            } else if (obj.id === "girl" && p_screen.stopped < 1) {
                p_screen.display();
            } else {
                obj.display();
            }
        }

    };

};

function getElementRight(id) {
    for (obj of rightobjects) {
        if (obj.id === id) {
            return obj
        }
    }
}