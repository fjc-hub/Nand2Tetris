
			@3030
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@THIS
			D=A
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	
			@3040
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@THAT
			D=A
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	
			@32
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@2
			D=A
			@THIS
			D=M+D
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	
			@46
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@6
			D=A
			@THAT
			D=M+D
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	
			@THIS
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@THAT
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
		
			@2
			D=A
			@THIS
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
		
			@6
			D=A
			@THAT
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
		