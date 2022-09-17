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

//pseudo code: 
// for true {
//     set v0 = -1
//     if KBD == 0 {
//         set v0 = 0
//     }
//     set addr = @SCREEN    // addr is a pointer (store memory address)
//     for @KBD > addr {
//         set M[*addr] = @v0
//         addr = addr + 1
//     }
// }

(OutLoop)
    @v0
    M=-1    //set @v0 = -1
    
    @KBD
    D=M     //load keyboard alue
    @Pressed
    D; JNE
    @v0
    M=0     //no pressed

(Pressed)
    @SCREEN
    D=A     //load @SCREEN address
    @addr
    M=D     //store @SCREEN address to variable addr (reset @addr)

    (InnerLoop)
        @addr
        D=M
        @KBD
        D=A-D
        @OutLoop
        D; JLE      // if @KBD - *(@addr) <= 0 then jump (out innerloop)
        @v0
        D=M
        @addr 
        A=M
        M=D         // set M[*(@addr)] = @v0
        @addr
        M=M+1       // @addr++
        @InnerLoop
        0; JMP
    
@OutLoop
0; JMP