touch = is_touch_device4();
let current = undefined;

function setup() {
    bg_color = color(BG);
    dark = color(LEFTC);
    hoverc = color(HOVERC);
    light = color(LIGHTC);
    let can = createCanvas(1280, 550);
    can.parent('test');
    leftPanel = new Panel(0, 0, 550, 550, 5);
    rightPanel = new Panel(730, 0, 550, 550, 5);
    topPanel = new Panel(550, 0, 180, 230, 5);
    bottomPanel = new Panel(550, 230, 180, 320, 5);
    let bubble = new SpeechBubble(550, 230, 160, "And so we meet again, my eternal rival. Shall we begin?", 20, "bubble", [0, 0]);
    bottomPanel.addTopLayer(bubble);
    //TODO change the size and location, make sure you can create 1 of the 2 speech bubbles, make this.clicked
}

function draw() {
    clear();
    background(bg_color);
    leftPanel.display();
    rightPanel.display();
    topPanel.display();
    bottomPanel.display();
}


function SpeechBubble(x, y, w, t, size, id, target) {
    this.x = x;
    this.y = y;
    this.w = w + 20;
    this.text = t;
    let lines = calculateLines(t, w, size);
    this.h = lines * size + (lines - 1) * 5;
    this.textSize = size;
    this.id = id;
    this.target = target;
    this.lifeTime = 60;
    this.fadeTime = 5;
    this.maxFadeTime = 5;
    this.transparency = 255;
    //this.hoverable = true;

    this.display = function () {
        textAlign(LEFT);
        noStroke();
        let c = color(255);
        c.setAlpha(this.transparency);
        fill(c);
        rect(this.x, this.y, this.w, this.h + 10, 5);
        strokeWeight(0.5);
        let col = color(dark.toString());
        col.setAlpha(this.transparency);
        stroke(col);
        fill(col);
        textSize(this.textSize);
        textLeading(25);
        text(this.text, this.x + 10, this.y + 5, this.w - 10);
        console.log(this.lifeTime);
        if (this.fadeTime <= 0) {
            removeElement(this.id);
        } else if (this.lifeTime > 0) {
            this.lifeTime -= 1;
        } else {
            this.fadeTime -= 1;
            this.transparency = this.fadeTime / this.maxFadeTime * 255;
        }
    };

    this.setText = function (text) {
        this.text = text;
        let lines = calculateLines(t, w, size);
        this.h = lines * size + (lines - 1) * 5;
    };
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

function removeElement(id) {
    for (let i = 0; i < leftPanel.objects.length; i++) {
        let obj = leftPanel.objects[i];
        if (obj.id === id) {
            leftPanel.objects.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < rightPanel.objects.length; i++) {
        let obj = rightPanel.objects[i];
        if (obj.id === id) {
            rightPanel.objects.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < topPanel.objects.length; i++) {
        let obj = topPanel.objects[i];
        if (obj.id === id) {
            topPanel.objects.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < bottomPanel.objects.length; i++) {
        let obj = bottomPanel.objects[i];
        if (obj.id === id) {
            bottomPanel.objects.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < leftPanel.toplayerobjs.length; i++) {
        let obj = leftPanel.toplayerobjs[i];
        if (obj.id === id) {
            leftPanel.toplayerobjs.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < rightPanel.toplayerobjs.length; i++) {
        let obj = rightPanel.toplayerobjs[i];
        if (obj.id === id) {
            rightPanel.toplayerobjs.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < topPanel.toplayerobjs.length; i++) {
        let obj = topPanel.toplayerobjs[i];
        console.log(obj.id);
        if (obj.id === id) {
            topPanel.toplayerobjs.splice(i, 1);
            return;
        }
    }
    for (let i = 0; i < bottomPanel.toplayerobjs.length; i++) {
        let obj = bottomPanel.toplayerobjs[i];
        if (obj.id === id) {
            bottomPanel.toplayerobjs.splice(i, 1);
            return;
        }
    }
}

