#!/usr/bin/python
from __future__ import print_function
from time import sleep
import os
import httplib
import json


# Greeter is a terminal application that greets old friends warmly,
#   and remembers new friends.


### FUNCTIONS ###

class HTTPError(Exception):
    def __init__(self, code, reason = ''):
        self.code = code
        self.reason = reason
    def __str__(self):
        return repr(self.code) + ":" + repr(self.reason)

conn = False

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

def makeGetRequest(host = "", path = ""):
    if not host:
        host = "localhost:8080"
    if not path:
        path = "/"
    global conn
    if not conn:
        conn = httplib.HTTPConnection("128.97.95.56:8080")
    conn.request("GET", path)
    resp = conn.getresponse()
    return resp


def getBoard(host = "", id = "0"):
    resp = makeGetRequest(host, path = "/game/" + str(id) + "/string")
    if resp.status != 200:
        raise HTTPError(resp.status, resp.reason)
        return
    data1 = resp.read()
    board = json.loads(data1)['Board']
    return board


def getGame():
    try:
        board = getBoard(id = 63714)
        for line in board:
            print(line)
    except HTTPError as e:
        print(str(e))

def main():
    getGame()

if __name__ == "__main__":
    main()
