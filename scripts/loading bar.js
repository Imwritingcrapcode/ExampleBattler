function LoadingBar(x, y, initial_w, h, radius, id, c, c2) {
    this.id = id;
    this.clickable = false;
    this.hoverable = false;
    this.rectColour = c;
    this.otherColour = c2;
    this.stopColour = c;
    this.percentage = 0.0;
    this.newPercentage = 0.0;
    this.x = x;
    this.y = y;
    this.initial_w = initial_w;
    this.w = 0;
    this.h = h;
    this.radius = radius;

    this.display = function () {
        this.w = round(this.initial_w * this.percentage / 100);
        this.stopColour = lerpColor(this.rectColour, this.otherColour, this.percentage / 100);
        this.setGradient(this.x, this.y, this.w, this.h, this.radius, this.c, this.stopColour);
        stroke(color);
        strokeWeight(1);
        noFill();
        rect(this.x, this.y, this.w, this.h, this.radius);
    };

    this.setPercentage = function(perc) {
        this.percentage = perc;
    };

    this.setNewPercentage = function(new_perc) {
        this.newPercentage = new_perc;
    };

    this.setGradient = function(x, y, w, h, r, c1, c2) {
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
}