
			@0
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@0
			D=A
			@LCL
			D=M+D
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	(BasicLoop.LOOP_START)

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
	
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		
			@0
			D=A
			@LCL
			D=M+D
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	
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
	
			@1
			D=A
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
		
			@0
			D=A
			@ARG
			D=M+D
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	
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
	
			@SP
			AM=M-1
			D=M			// pop the top-most value of the stack onto D-register
			@BasicLoop.LOOP_START
			D;JNE
	
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
	