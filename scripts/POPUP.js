idCounter = 0;
transitionEvent = whichTransitionEvent();

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
    let popup = document.getElementById(id);
    if (popup.classList.contains("remove")) {
        popup.classList.remove("remove");
    }
    if (!popup.classList.contains("show")) {
        popup.classList.add("show");
    }
}

function setText(id, t) {
    let regExpEmoji = /:([a-zA-Z0-9_]+):/;
    let regExpNextEmojiName = /(?<=:)[a-zA-Z0-9_]+(?=:)/;
    let final_text = t;
    let emoji = final_text.match(regExpNextEmojiName);
    while (emoji) {
        final_text = final_text.replace(regExpEmoji, "<img class='notifImage' src='/images/locked/emojis/" + emoji + ".png'>");
        emoji = final_text.match(regExpNextEmojiName);
    }
    document.getElementById(id).innerHTML = final_text;
}

function addPopup(text, where, ma) {
    let id = "";
    if (!!ma) {
        id = ma;
    } else {
        id = "popup" + idCounter++;
    }
    let newNotif = document.createElement('span');
    let a1 = document.createAttribute("class");
    a1.value = "popup";
    let a2 = document.createAttribute("shouldredirect");
    a2.value = where;
    let a3 = document.createAttribute("id");
    a3.value = id;
    newNotif.setAttributeNode(a1);
    newNotif.setAttributeNode(a2);
    newNotif.setAttributeNode(a3);
    document.getElementById("popupWrapper").appendChild(newNotif);
    let elements = document.getElementsByClassName("popup");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("click", onClickForPopups, false);
        elements[i].WhatToDoWhenAPopupIsClicked = function () {
            let where = this.getAttribute("shouldredirect");
            if (where.length > 0) {
                console.log("where", where);
                location = where;
            }
        };
    }

    setText(id, text);
    showPopup(id);
}

function dieFirstPopup() {
    let wrapper = document.getElementById("popupWrapper");
    if (wrapper.children.length > 0) {
        let item = wrapper.children.item(0);
        if (!item.classList.contains("remove")) {
            item.classList.add("remove");
            item.addEventListener(transitionEvent, customFunction);
        }
        window.setTimeout("dieFirstPopup()", 1000 * POPUPLIFETIME);
    }
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

function getNotifications() {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/notifications', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                let response = JSON.parse(xhr.responseText);
                for (let notification of response) {
                    addPopup(notification[0], notification[1]);
                }
                if (response.length > 0) {
                    window.setTimeout("dieFirstPopup()", 1000 * POPUPLIFETIME);
                }
            } else {
                console.log(xhr.status, xhr.responseText);
            }
        }
    };
}

getNotifications();
function handleVisibilityChange() {
    if (!document.hidden) {
        getNotifications();
    }
}
document.addEventListener("visibilitychange", handleVisibilityChange, false);
