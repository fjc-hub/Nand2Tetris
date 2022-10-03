
		@256	// set *SP = 256
		D=A
		@SP
		M=D
	
			@Sys.init.return-address.0
			D=A		// push return-address
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@LCL	// store LCL state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@ARG 	// store ARG state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@THIS	// store THIS state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@THAT	// store THAT state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@SP 	// update new ARG
			D=M
			@0
			D=D-A
			@5
			D=D-A
			@ARG
			M=D 	
			@SP		// update LCL pointer
			D=M
			@LCL
			M=D
			@Sys.init		// transfer control
			0; JMP
		(Sys.init.return-address.0)
	
		(Sys.init)
	
			@4000
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
	
			@5000
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
	
			@Sys.main.return-address.1
			D=A		// push return-address
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@LCL	// store LCL state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@ARG 	// store ARG state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@THIS	// store THIS state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@THAT	// store THAT state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@SP 	// update new ARG
			D=M
			@0
			D=D-A
			@5
			D=D-A
			@ARG
			M=D 	
			@SP		// update LCL pointer
			D=M
			@LCL
			M=D
			@Sys.main		// transfer control
			0; JMP
		(Sys.main.return-address.1)
	
			@1
			D=A
			@R5
			D=A+D
		
			@R13
			M=D
			@SP
			AM=M-1
			D=M
			@R13
			A=M
			M=D
	(Sys.init$LOOP)

			@Sys.init$LOOP
			0;JMP
	
		(Sys.main)
	
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
		
			@SP
			A=M
			M=0
			@SP
			M=M+1
		
			@4001
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
	
			@5001
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
	
			@200
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@1
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
	
			@40
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@2
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
	
			@6
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@3
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
	
			@123
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@Sys.add12.return-address.2
			D=A		// push return-address
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@LCL	// store LCL state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@ARG 	// store ARG state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@THIS	// store THIS state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@THAT	// store THAT state of caller
			D=M
			@SP 
			A=M
			M=D
			@SP
			M=M+1
			@SP 	// update new ARG
			D=M
			@1
			D=D-A
			@5
			D=D-A
			@ARG
			M=D 	
			@SP		// update LCL pointer
			D=M
			@LCL
			M=D
			@Sys.add12		// transfer control
			0; JMP
		(Sys.add12.return-address.2)
	
			@0
			D=A
			@R5
			D=A+D
		
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
	
			@2
			D=A
			@LCL
			A=M+D
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@3
			D=A
			@LCL
			A=M+D
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@4
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
		
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		
			@SP 	// A=sp
			AM=M-1 	// A=(*sp)-1; (*sp)--
			D=M 	// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M+D
		
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
	
		(Sys.add12)
	
			@4002
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
	
			@5002
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
	
			@12
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
	