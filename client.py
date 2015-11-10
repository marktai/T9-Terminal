#!/usr/bin/python
from __future__ import print_function
from __future__ import division
# from __future__ import input
from time import sleep
import os
import httplib
import json


# Greeter is a terminal application that greets old friends warmly,
#   and remembers new friends.

class clear:
    def __call__(self):
        if os.name==('ce','nt','dos'): os.system('cls')
        elif os.name=='posix': os.system('clear')
        else: print('\n'*120)
    def __neg__(self): self()
    def __repr__(self):
        self();
        return ''

clear=clear()

### FUNCTIONS ###

class HTTPError(Exception):
    def __init__(self, resp):
        self.status = resp.status
        self.reason = resp.reason
        if resp:
            try:
                self.error = json.loads(resp.read())['Error']
            except Exception as e:
                self.error = ''

    def __str__(self):
        if self.error:
            return repr(self.status) + ":" + str(self.error)
        return repr(self.status) + ":" + repr(self.reason)

    def __bool__(self):
        return self.status != 200


conn = False

boxToStringTranslator = {}
stringToBoxTranslator = {}

for i in range(9):
    out = ""
    height = i // 3
    if height == 0:
        out += "top-"
    elif height == 1:
        out += "middle-"
    elif height == 2:
        out += "bottom-"

    width = i % 3
    if width == 0:
        out += "left"
    elif width == 1:
        out += "middle"
    elif width == 2:
        out += "right"

    boxToStringTranslator[i] = out
    boxToStringTranslator[str(i)] = out
    stringToBoxTranslator[out] = str(i)

def display_title_bar():
    # Clears the terminal screen, and displays a title bar.
    os.system('clear')
              
    print("\t**********************************************")
    print("\t***  Greeter - Hello old and new friends!  ***")
    print("\t**********************************************")
    
def test():

    ### MAIN PROGRAM ###    

    # Print a bunch of information, in short intervals
    names = ['aaron', 'brenda', 'cyrene', 'david', 'eric']

    # Print each name 5 times.
    for name in names:
        display_title_bar()

        print("\n\n")
        for x in range(0,5):
            print(name.title())
        
        # Pause for 1 second between batches.
        sleep(1)

def makeGetRequest(host = "", path = "", type = "GET"):
    if not host:
        host = "localhost:8080"
    if not path:
        path = "/"
    global conn
    # if not conn:
    conn = httplib.HTTPConnection(host)
    # print("made connection")
    conn.request(type, path)
    resp = conn.getresponse()
    return resp


def makePostRequest(host = "", path = ""):
    return makeGetRequest(host, path, type = "POST")


def getBoardString(host = "", id = "0"):
    resp = makeGetRequest(host, path = "/game/" + str(id) + "/string")
    if resp.status != 200:
        raise HTTPError(resp)
        return
    data1 = resp.read()
    board = json.loads(data1)['Board']
    return board

def getGameInfo(host = "", id = "0"):
    resp = makeGetRequest(host, path = "/game/" + str(id))
    if resp.status != 200:
        raise HTTPError(resp)
        return
    data1 = resp.read()
    game = json.loads(data1)['Game']
    return game

def parseHeader(game, player):
    turnString = ""
    if str(game["Players"][game["Turn"] // 10]) == str(player):
        turnString = "Your Turn"
    else:
        turnString = "Not Your Turn"
    box = boxToStringTranslator[game["Turn"] % 10]
    out = "ID: %s | %s | Box: %s" % (game["GameID"], turnString, box)
    return out

def getGame(host = "", id = ""):
    try:
        game = getGameInfo(host, id = id)
        return game
    except HTTPError as e:
        print("Error retrieving board")
        print(str(e))
        return False

def normalizeSquare(inpString):
    if inpString in boxToStringTranslator:
        return inpString
    inpString = inpString.lower().replace("  ", " ").replace(" ", "-")
    if inpString in stringToBoxTranslator:
        return stringToBoxTranslator[inpString]
    return -1


def makeMove(host, id, player, box, square):
    query = "?Player=%s&Box=%s&Square=%s" % (str(player), str(box), str(square))
    return makePostRequest(host, path = "/game/" + str(id) + query)

def ui():
    userInput = ""
    host = raw_input("Host:")
    id = raw_input("ID:")
    player = raw_input("Player ID:")
    id = "63714"
    player = "0"
    printGameHeader = True
    printGameBoard = True
    printGameInfo = False

    afterText = ""

    stayInLoop = True

    while stayInLoop:
        clear()
        game = getGame(host, id)
        if game:
            if printGameHeader:
                print(parseHeader(game,player))
            if printGameInfo:
                for key in game:
                    print("%s: %s" % (key, game[key]))
            if printGameBoard:
                board = getBoardString(host, id = id)
                for line in board:
                    print(line)

            if afterText:
                print("\n" + afterText)
                afterText = ""

        first = True
        while first:
            first = False

            userInput = raw_input("\nCommand (r, h, i, m, p, s, q):")
            if userInput.lower() == "r" or userInput.lower() == "refresh":
                pass

            elif userInput.lower() == "h" or userInput.lower() == "header":
                printGameHeader = True
                printGameInfo = False
                printGameBoard = True
            elif userInput.lower() == "i" or userInput.lower() == "info":
                printGameHeader = False
                printGameInfo = True
                printGameBoard = True

            elif userInput.lower() == "m" or userInput.lower() == "move":
                square = -1
                while square == -1:
                    inputSquare = raw_input("Square:")
                    if inputSquare.lower() == "b" or inputSquare.lower() == "back":
                        break
                    square = normalizeSquare(inputSquare)
                if square != -1:
                    resp = makeMove(host, id, player, game["Turn"]%10, square)
                    e = HTTPError(resp)
                    if e:
                        afterText += str(e)

                
            elif userInput.lower() == "s" or userInput.lower() == "switch":
                tempid = raw_input("ID:")
                id = tempid if tempid else id
                tempplayer = raw_input("Player ID:")
                player = tempplayer if tempplayer else player

            elif userInput.lower() == "p" or userInput.lower() == "player":
                tempplayer = raw_input("Player ID:")
                player = tempplayer if tempplayer else player

            elif userInput.lower() == "q" or userInput.lower() == "quit":
                stayInLoop = False

            elif userInput.lower() == "e" or userInput.lower() == "exit":
                stayInLoop = False

            else:
                first = True


    clear()
    print("Thanks for playing!")




def main():
    ui()

if __name__ == "__main__":
    main()
