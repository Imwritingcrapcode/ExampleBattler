idCounter = 0;
transitionEvent = whichTransitionEvent();
document.body += "<div align=\"right\" id=\"popupWrapper\" class=\"popupWrapper\"></div>";

function whichTransitionEvent() {
    let t;
    let el = document.createElement('fakeelement');
    const animations = {
        'animation': 'animationend',
        'Oanimation': 'oanimationEnd',
        'Mozanimation': 'animationend',
        'Webkitanimation': 'webkitanimationEnd'
    };

    for (t in animations) {
        if (el.style[t] !== undefined) {
            return animations[t];
        }
    }
}

function showPopup(id) {
    console.log("show", id);
    let popup = document.getElementById(id);
    if (popup.classList.contains("remove")) {
        popup.classList.remove("remove");
    }
    if (!popup.classList.contains("show")) {
        popup.classList.add("show");
    }
}

function setText(id, text) {
    document.getElementById(id).innerText = text;
}

function addPopup(text, where, ma) {
    let id = "";
    if (!!ma) {
        id = ma;
    } else {
        id = "popup" + idCounter++;
    }
    document.getElementById("popupWrapper").innerHTML += "<span class=\"popup\" shouldredirect=\"" + where + "\" id=\"" + id + "\"></span>";


    let elements = document.getElementsByClassName("popup");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("click", onClickForPopups, false);
        elements[i].WhatToDoWhenAPopupIsClicked = function () {
            console.log(this.getAttribute("shouldredirect"));
        };

    }

    setText(id, text);
    showPopup(id);
}


function onClickForPopups(ev) {
    this.WhatToDoWhenAPopupIsClicked();
    if (!this.classList.contains("remove")) {
        this.classList.add("remove");
        this.addEventListener(transitionEvent, customFunction);
    }
}

function customFunction(ev) {
    this.removeEventListener(transitionEvent, customFunction);
    let elements = document.getElementById("popupWrapper");
    if (elements.contains(this)) {
        elements.removeChild(this);
    }
}