#!/usr/bin/python
from __future__ import print_function
from time import sleep
import os
import httplib


# Greeter is a terminal application that greets old friends warmly,
#   and remembers new friends.


### FUNCTIONS ###

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

def getGame():
        conn = httplib.HTTPSConnection("128.97.95.56:8080")
        #conn.request("GET", "/game/63714")
        #conn = httplib.HTTPSConnection("www.marktai.com")
        conn.request("GET", "/")
        r1 = conn.getresponse()
        print (r1.status, r1.reason)
        data1 = r1.read()
        print (data1)

def main():
    getGame()

if __name__ == "__main__":
    main()
