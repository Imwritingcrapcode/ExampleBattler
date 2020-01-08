let cnv;
let d;
function setup() {
    cnv = createCanvas(100, 100);
    cnv.mouseOut(changeD);
    d = 10;
}
function draw() {
    ellipse(width / 2, height / 2, d, d);
}

function changeD() {
    d = d + 10;
}