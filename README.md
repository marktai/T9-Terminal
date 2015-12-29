# Tic Tac Toe Squared Terminal Client

A basic client for the Tic Tac Toe Squared (AKA Meta Tic Tac Toe) server I created.  

This is written in Go, mostly because I was mad that Python websockets weren't working nicely and threading was easy in Go.

#### Some features:

 *  Amazing 20xx graphics
 *  Automatically rescales to terminal size
 *  Adjustable screen refresh time for slow computers
 *  Cross platform (tested on Ubuntu, OSX, and Windows)
 *  Utilizes all features from [server](https://github.com/marktai/T9-server)
 *  Scrollable information 
 *  Websockets to automatically refresh game on change

## How to Play:
 1. Run the executable (two options)
    1. [Install Golang](https://golang.org/dl/), clone repo, and use make run
    2. Download and run (after allowing it in your antivirus)
       * [Linux](https://www.marktai.com/upload/T9clientLinux)
       * [OSX](https://www.marktai.com/upload/T9clientOSX) (currently has problems accessing root certificates used in https)
       * [Windows](https://www.marktai.com/upload/T9clientWindows.exe)  
 2. Leave host empty for https://www.marktai.com/T9 or use your own host
 3. Either create or join a game
    1. If creating a game, enter any two arbitrary player ID's.
    2. If joining a game, either use one of the game ID's on the right and select a player 
 4. Use m to make a move and then select a box (1/9 of the board) and a square (1/9 of the box). 
 5. Send all bugs you find to mark@marktai.com
  
