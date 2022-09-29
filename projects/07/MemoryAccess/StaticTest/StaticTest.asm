
			@111
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@333
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@888
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			@StaticTest.8
			M=D
		
			@SP
			AM=M-1
			D=M
			@StaticTest.3
			M=D
		
			@SP
			AM=M-1
			D=M
			@StaticTest.1
			M=D
		
			@StaticTest.3
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@StaticTest.1
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
		
			@StaticTest.8
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
		