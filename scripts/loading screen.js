loading = function(p) {

    p.setup = function () {
        dark = p.color(DARKC);
        light = p.color(LIGHTC);
        right = p.color(RIGHTC);
        last = p.color(LVLUPC);

        plscreen = new LoadingScreen(p, 0.4 * height, 0.4 * height, 0.2 * height, 0.2 * height);
        plscreen.setColours(dark, light, right, last);
        oppscreen = new LoadingScreen(p, 1.4 * height, 0.4 * height, 0.2 * height, 0.2 * height);
        oppscreen.setColours(dark, light, right, last);
        //r = 2sin(2phi)
        //x = rcos(phi)
        //y = rsin(phi)
        console.log()
    };

    p.draw = function () {
        background(bg_color);
        plscreen.display();
        oppscreen.display();
        /*let x = width / 2;
        let y = height / 2;
        strokeWeight((width * 0.4) / 150);
        stroke(dark);
        for (let phi = 0; phi <= PI * 2; phi += 0.01) {
            let r = 2 * sin(2 * phi) * (lscreen.w/2);
            let new_x = r * cos(phi) + width / 2;
            let new_y = r * sin(phi) + height / 2;
            line(x, y, new_x, new_y);
            x = new_x;
            y = new_y;
            let inter = map(phi, 0.0, 360.0, 0, 1);
            let c = lerpColor(dark, right, inter);
            stroke(c);
        }*/

    };

};
