
		(SimpleFunction.test)
	
			@SP
			A=M
			M=0
			@SP
			M=M+1
		
			@SP
			A=M
			M=0
			@SP
			M=M+1
		
			@0
			D=A
			@LCL
			A=M+D
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@1
			D=A
			@LCL
			A=M+D
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		
			@SP
			A=M-1
			M=!M
		
			@0
			D=A
			@ARG
			A=M+D
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		
			@1
			D=A
			@ARG
			A=M+D
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP		// A=sp
			AM=M-1	// A=(*sp)-1; (*sp)--
			D=M		// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M-D
		
			@LCL 	// store return address to avoid being overwitten
			D=M
			@5
			A=D-A
			D=M
			@R13
			M=D
			@SP 	// transfer top-most value of stack onto ARG[0]
			A=M-1
			D=M
			@ARG
			A=M
			M=D
			D=A		// reset SP
			@SP
			M=D+1
					// restore all preserved states of caller
			@LCL	// restore THAT
			D=M
			@1
			A=D-A
			D=M
			@THAT
			M=D
			@LCL	// resotre THIS
			D=M
			@2
			A=D-A
			D=M
			@THIS
			M=D
			@LCL 	// restore ARG
			D=M
			@3
			A=D-A
			D=M
			@ARG
			M=D
			@LCL	// restore LCL
			D=M
			@4
			A=D-A
			D=M
			@LCL
			M=D
			@R13 	// jmp to ret-addr
			A=M
			0; JMP
	