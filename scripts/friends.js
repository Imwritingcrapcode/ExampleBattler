friends = [];
incoming = [];
pending = [];

function init() {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', '/friendlist', true);
    xhr.send();
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            let response = JSON.parse(xhr.responseText);
            console.log(response);
            setFriends(response.friends);
            setIncoming(response.incoming);
            setPending(response.pending);
        }
    };
}

function openTab(evt, tabName) {
    let i, tabcontent, tablinks;
    //hiding all the tab contents
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    //deactivating all buttons
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    //displaying current tab
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}

function setFriends(friends) {
    let innerHTML = "";
    if (!!friends) {
        friends = friends.sort((a, b) => (a[0] > b[0]) ? 1 : ((b[0] > a[0]) ? -1 : 0));
        innerHTML += "<table width='100%'>";
        for (let friend of friends) {
            innerHTML += "<tr><td width=\"33%\">";
            console.log("you have friends");
            innerHTML += friend[0] + "</td>";
            let middle = "<td>" + friend[1];
            if (friend.length === 3) {
                if (friend[1] === "Playing as") {
                    middle += " " + friend[2];
                } else if (friend[1] === "Offline") {
                    let time = parseSeconds(parseInt(friend[2]), true);
                    if (time === "-1") {
                        console.log("so long omg", friend[0])
                    } else {
                        middle += " for " + time;
                    }
                }
            }
            middle += "</td>";
            innerHTML += middle;
            let postfix = "<div align='right'><button class=\"tablinks\" onclick='removeFriend(\"" + friend[0] + "\")'>✖</button></div>";
            innerHTML += "<td width='30px'>" + postfix + "</td>";
            innerHTML += "</tr>";
        }
        innerHTML += "</table>";
    } else {
        console.log("you have 0 friends uwu");
        innerHTML = "<div align=\"center\">You have 0 friends uwu</div>";
    }
    document.getElementById("friends").innerHTML = innerHTML;
}

function setIncoming(incoming) {
    let innerHTML = "";
    if (!!incoming) {
        console.log("you have incoming");
        incoming = incoming.sort((a, b) => (a[0] > b[0]) ? 1 : ((b[0] > a[0]) ? -1 : 0));
        innerHTML += "<table width='100%'>";
        for (let friend of incoming) {
            innerHTML += "<tr><td width=\"33%\">" + friend + "</td><td>";
            innerHTML += "<div align='right'><button class=\"tablinks\" onclick='addFriend(\"" + friend + "\")'>✔</button>\t<button class=\"tablinks\" onclick='removeFriend(\"" + friend + "\")'>✖</button></div></td>";
            innerHTML += "</tr>"
        }
        innerHTML += "</table>";
    } else {
        console.log("you have 0 inc uwu");
        innerHTML = "<div align=\"center\">No incoming requests.</div>"
    }
    document.getElementById("incoming").innerHTML = innerHTML;

}

function setPending(pending) {
    let innerHTML = "";
    if (!!pending) {
        console.log("you have pending");
        pending = pending.sort((a, b) => (a[0] > b[0]) ? 1 : ((b[0] > a[0]) ? -1 : 0));
        innerHTML += "<table width='100%'>";
        for (let friend of pending) {
            innerHTML += "<tr><td width=\"33%\">" + friend + "</td><td>";
            innerHTML += "<div align='right'><button class=\"tablinks\" onclick='removeFriend(\"" + friend + "\")'>✖</button></div></td></tr>";
        }
        innerHTML += "</table>";
    } else {
        console.log("you have 0 pending uwu");
        innerHTML = "<div align=\"center\">No pending requests. Add someone! :></div>";
    }
    document.getElementById("pending").innerHTML = innerHTML;
}

function removeFriend(name) {
    console.log("remove: " + name);
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/friendlist', true);
    xhr.send(JSON.stringify(["Remove", name]));
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                alert(xhr.responseText);
                init();
            } else {
                alert(xhr.responseText);
                alert(response);

            }
        }
    };
}

function addFriend(name) {
    console.log("add: " + name);
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/friendlist', true);
    xhr.send(JSON.stringify(["Add", name]));
    xhr.onreadystatechange = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                alert(xhr.responseText);
                init();
            } else {
                alert(xhr.responseText);
            }
        }
    };
}
