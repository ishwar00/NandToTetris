// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.
    // size = 8192 memory block size
    @8192
    D=A
    @size
    M=D
(INITIALISE-BLACKING)
// i = 0 
    @i
    M=0
(BLACKING)
    // KBD != 0
    @KBD
    D=M
    @INITIALISE-CLEANING
    D;JEQ
    // okay KBD != 0 true 
    // i != size 
    @size
    D=M
    @i
    D=D-M
    @BLACKING
    D;JEQ
    // i != size is true 
    @SCREEN
    D=A
    @i
    A=D+M
    M=-1 // RAM[SCREEN + i] = -1 
    @i
    M=M+1
    @BLACKING
    0;JMP
(INITIALISE-CLEANING)
    // i = 0 
    @i
    M=0
(CLEANING)
    @KBD
    D=M
    @INITIALISE-BLACKING
    D;JNE
    // KBD == 0 is true 
    @i
    D=M
    @size
    D=D-M
    @CLEANING
    D;JEQ
    // i != size is true 
    @SCREEN
    D=A
    @i
    A=D+M
    M=0 // RAM[SCREEN + i] = 0 
    @i
    M=M+1
    @CLEANING
    0;JMP