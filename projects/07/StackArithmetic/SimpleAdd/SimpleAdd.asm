
			@7
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@8
			D=A
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
		