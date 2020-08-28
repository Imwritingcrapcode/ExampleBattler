from random import randint

costs = {'ST':300, 'AD':370, 'SP':775, 'RP':1435, 'LF':1795}
ST_CHANCE = 450
AD_CHANCE = 300
SP_CHANCE = 150
RP_CHANCE = 80
LF_CHANCE = 20

def sort_value(t):
    if t == 'ST':
        return 1
    elif t == 'AD':
        return 2
    elif t == 'SP':
        return 3
    elif t == 'RP':
        return 4
    elif t == 'LF':
        return 5

def generate_pack():
    n = randint(0, 999)
    pack = ['', '']
    if n < ST_CHANCE:
        pack[0] = 'ST'
    elif n < ST_CHANCE + AD_CHANCE:
        pack[0] = 'AD'
    elif n < ST_CHANCE + AD_CHANCE + SP_CHANCE:
        pack[0] = 'SP'
    elif n < ST_CHANCE + AD_CHANCE + SP_CHANCE + RP_CHANCE:
        pack[0] = 'RP'    
    else:
        pack[0] = 'LF'
    n = randint(0, 999)
    if n >= 1000 - LF_CHANCE and pack[0] == 'ST':
        pack[1] = 'LF'
    elif n >= 1000 - LF_CHANCE - RP_CHANCE and sort_value(pack[0]) <= sort_value('AD'):
        pack[1] = 'RP'
    elif n >= 1000 - LF_CHANCE - RP_CHANCE - SP_CHANCE and sort_value(pack[0]) <= sort_value('SP'):
        pack[1] = 'SP'
    elif  n >= 1000 - LF_CHANCE - RP_CHANCE - SP_CHANCE - AD_CHANCE and sort_value(pack[0]) <= sort_value('RP'):
        pack[1] = 'AD'    
    else:
        pack[1] = 'ST'
    return pack

def value(pack):
    return (costs[pack[0]] + costs[pack[1]])/2

def generate_packs(n):
    total = 0
    for i in range(n):
        pack = generate_pack()
        v = value(pack)
        total += v
        print(', '.join(sorted(pack, key=sort_value)) + ', value:', v)
    return total/n


print(generate_packs(20))
