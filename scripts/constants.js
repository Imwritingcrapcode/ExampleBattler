//SYSTEM COLOURS
BG = 230;
LEFTC = 50;
RIGHTC = "rgb(120,180,255)";
LIGHTC = 170;
OTHERLIGHTC = 135;
//DARKC = 80;
DARKC = LEFTC;
LVLUPC = 'rgb(255, 182, 13)';
CLICKABLEC = 'rgba(150, 150, 255, 1)';
ACTIVEC = 'rgba(130, 130, 255, 1)';
CLICKEDC = 'rgba(120, 120, 200, 1)';
HOVERC = 'rgba(255, 255, 255, 0.87)';
WINC = 'rgb(255, 182, 13)';

//SHORTCUTS FOR COLOURS
RED = 'rgba(255, 65, 51, 1)';
ORANGE = 'rgba(255, 131, 36, 1)';
YELLOW = 'rgba(200, 200, 0, 1)';
GREEN = 'rgba(0, 200, 0, 1)';
CYAN = 'rgba(0, 200, 200, 1)';
BLUE = 'rgba(36, 65, 255, 1)';
VIOLET = 'rgba(105, 36, 255, 1)';
PINK = 'rgba(210, 0, 170, 1)';
GREY = 'rgba(90, 90, 90, 1)';
BROWN = 'rgba(120, 50, 35, 1)';
BLACK = 'rgba(0, 0, 0, 1)';
WHITE = 'rgba(200, 200, 200, 1)';

//RARITY COLOURS
STCOLOUR = DARKC;
ADCOLOUR = BLUE;
SPCOLOUR = 'rgba(208, 200, 0, 1)';
RPCOLOUR = VIOLET;
LFCOLOUR = WINC;

FRAMESFORANIMATIONS = 60;
POPUPLIFETIME = 20;

EFFECTDESCRIPTIONS = [
    [0, "Doubles the damage you deal."],
    [2, "Prevents you from healing."],
    [3, "You can't use skills of the same colour you used last."],
    [4, "This turn, your opponent chooses which skills you use."],
    [6, "Your opponent has to deal more damage than the threshold, or no damage at all."],
    [7, "Your opponent's next attack will deal this much less damage."],
    [8, "Your opponent can't use debuffs on you. Your Royal Move and Composure become stronger."],
    [9, "This turn, you can use two skills but not your ultimate."],
    [10, "If you survive your opponent's next turn, fully heals you."],
    [12, "You heal from Euphoric Source at the end of each turn, but Source gets rapidly depleted."],
    [13, "Green tokens add Green damage to your Stab."],
    [14, "Black tokens add Black damage to your Stab."],
    [15, "Damage dealt by Royal Move. You can spend it on Composure heal or Pride damage."],
    [16, "Boosts your Electric Shock damage."],
    [17, "Boosts your Layers defense."],
    [18, "Euphoric Source gives your Pink Sphere additional damage as well as well as healing while in Euphoria."],
];

COLOURIDS = [
    ["Red", RED],
    ["Orange", ORANGE],
    ["Yellow", YELLOW],
    ["Green", GREEN],
    ["Cyan", CYAN],
    ["Blue", BLUE],
    ["Violet", VIOLET],
    ["Pink", PINK],
    ["Gray", GREY],
    ["Brown", BROWN],
    ["Black", BLACK],
    ["White", WHITE],
];

DUSTS = new Map([
    ["w", "White"],
    ["b", "Blue"],
    ["y", "Yellow"],
    ["p", "Purple"],
    ["s", "Star"],
]);

SKILLDESCRIPTIONS = new Map([
    ["Your Number", "Deal 10 + the remainder of your opponent's number divided by 7 Orange damage."],
    ["Your Colour", "Next turn, your opponent can't use the skills of the same colour they used last. Deal 15 damage of that colour.\nUnlocks when your opponent uses a skill.\nCooldown: 1."],
    ["Your Dream", "Heal for (your max HP - your opponent's number) / your turn number.\nIf your opponent's number is more than 83, subtract a flat number as if it was 83."],
    ["My Story", "Next turn, you decide which skills your opponent uses.\nUnlocks on turn 7.\nCooldown 1."],
    ["Dance", "Double all of your damage.\nLasts 2 turns."],
    ["Rage", "Deal 24 - 2 * your turn number Red damage."],
    ["Stop", "Every player can not heal until the end of their next turn. While this is active for you, .Execute becomes stronger.\nCooldown 1."],
    [".Execute", "If your opponent's at less than 10% of their max HP, defeat them instantly.\nWhile Stop effect is active, the threshold goes to 20% of opponent's max hp."],
    ["Scarcity", "Deal 12 Black damage, then set opponent's max HP to their current HP.\nCooldown 1."],
    ["Indifference", "If opponent's ultimate is not available yet, delay it for 1 turn. Can't be delayed later than their 10th turn.\nCooldown 2. Unlocks on turn 2."],
    ["Green Sphere", "Deal 15 - (opponent's max HP - opponent's current HP) green damage."],
    ["Despondency", "Deal 40 - (opponent's max HP - 70) Blue damage.\nUnlocks on turn 9."],
    ["Ampleness", "Increases Euphoric Source and everyone's max HP by 12.\nCooldown: 1."],
    ["Exuberance", "If your opponent's ultimate is not unlocked yet, add 10 to Euphoric Source, everyone's max HP and your current HP. Also, your opponent's ultimate unlocks 1 turn earlier.\nIf it already is unlocked, add 20 instead.\nCooldown 2."],
    ["Pink Sphere", "Deal 12 Pink damage. Also, increase everyone's max HP by 12."],
    ["Euphoria", "Heal everyone by the amount in Euphoric Source at the end each turn.\nEvery turn, Source gets depleted by 9. Lasts until the end of the game.\nStarting turn: 4."],
    ["Run", "Your opponent's next attack will deal 5 less damage. Gain a Green Token."],
    ["Weaken", "Reduce your opponent's defense to Green by 1. Gain a Black Token."],
    ["Speed", "Next turn, you'll use two skills but not your ultimate. Gain a Green Token."],
    ["Stab", "For each of your tokens, deal 6 Green&Black damage."],
    ["Royal Move", "Deal 12 Green damage and add that to Stolen HP.\nWith Mint Mist, deal 20 Green damage instead."],
    ["Composure", "Spend some Stolen HP to heal yourself for up to 20.\nWith Mint Mist, heal up to 30."],
    ["Mint Mist", "You become invisible, your opponent can't debuff you. Your Royal Move and Composure become stronger.\nLasts 2 turns.\nCooldown: 2."],
    ["Pride", "Spend all of your Stolen HP to deal as much Cyan damage.\nUnlocks on turn 8."],
    ["E-Shock", "Deal Cyan damage. Base damage is 5, gets to 10, 15 and 20 when boosted by I Boost."],
    ["I Boost", "Boost your S Layers threshold by 5 and E-Shock damage by 5. Can only be used three times in a match."],
    ["S Layers", "Next turn, your opponent can't damage you unless they deal more than a certain threshold.\nThresholds are 5, 10, 15 and 20.\nGain 1 Defense against all colours but Black."],
    ["Last Chance", "If you survive your opponent's next turn, fully heal.\nUnlocks on turn 7.\nCan only be used once per match."],
]);


MATCHES = new Map([
    ["ST", [1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 5, 5, 5, 5]],
    ["AD", [1, 1, 2, 2, 2, 2, 3, 3, 4, 4, 4, 4, 4, 5, 5, 6, 6, 6, 6]],
    ["SP", [1, 2, 4, 5, 5, 5, 6, 8, 9, 11, 11, 11, 11, 11, 13, 14, 16, 16, 16]],
    ["RP", [2, 6, 10, 13, 13, 14, 14, 16, 20, 23, 24, 24, 25, 27, 29, 31, 35, 37, 37]],
    ["LF", [4, 7, 12, 16, 17, 18, 20, 23, 27, 30, 31, 31, 31, 33, 34, 38, 41, 43, 44]],
]);

CONVERSIONRATE = {
    b:0.4,
    g:0.2,
    w:0.5,
    y:0.25
};

SECONDSPERCONVERSION = {
    w: 24,
    b: 30,
    y: 45,
    p: 60,
};


colourMap = new Map([
    ["Red", 0],
    ["Orange", 1],
    ["Yellow", 2],
    ["Green", 3],
    ["Cyan", 4],
    ["Blue", 5],
    ["Violet", 6],
    ["Pink", 7],
    ["Gray", 8],
    ["Brown", 9],
    ["Black", 10],
    ["White", 11],
]);

raritiesMap = new Map([
    ["ST", 0],
    ["AD", 1],
    ["SP", 2],
    ["RP", 3],
    ["LF", 4],
]);


function getResolution(num) {
    switch (num) {
        case 1:
            return [384, 550];
        case 8:
            return [211, 550];
        case 9:
            return [218, 550];
        case 10:
            return [535, 550];
        case 33:
            return [413, 550];
        case 51:
            return [350, 550];
        case 119:
            return [324, 550];
        default:
            return [0, 0];
    }
}

function existsPortrait(num) {
    switch (num) {
        case 1:
            return true;
        case 8:
            return true;
        case 9:
            return true;
        case 10:
            return true;
        case 33:
            return true;
        case 51:
            return true;
        case 119:
            return true;
        default:
            return false;
    }
}