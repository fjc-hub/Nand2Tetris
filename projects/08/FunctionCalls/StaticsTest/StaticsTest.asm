
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
	
			@6
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
	
			@Class1.set.return-address.1
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
			@2
			D=D-A
			@5
			D=D-A
			@ARG
			M=D 	
			@SP		// update LCL pointer
			D=M
			@LCL
			M=D
			@Class1.set		// transfer control
			0; JMP
		(Class1.set.return-address.1)
	
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
	
			@23
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@15
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@Class2.set.return-address.2
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
			@2
			D=D-A
			@5
			D=D-A
			@ARG
			M=D 	
			@SP		// update LCL pointer
			D=M
			@LCL
			M=D
			@Class2.set		// transfer control
			0; JMP
		(Class2.set.return-address.2)
	
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
	
			@Class1.get.return-address.3
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
			@Class1.get		// transfer control
			0; JMP
		(Class1.get.return-address.3)
	
			@Class2.get.return-address.4
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
			@Class2.get		// transfer control
			0; JMP
		(Class2.get.return-address.4)
	(Sys.init$WHILE)

			@Sys.init$WHILE
			0;JMP
	
		(Class1.set)
	
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
			D=M
			@Class1.0
			M=D
		
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
	
			@SP
			AM=M-1
			D=M
			@Class1.1
			M=D
		
			@0
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
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
	
		(Class1.get)
	
			@Class1.0
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@Class1.1
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
	
		(Class2.set)
	
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
			D=M
			@Class2.0
			M=D
		
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
	
			@SP
			AM=M-1
			D=M
			@Class2.1
			M=D
		
			@0
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
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
	
		(Class2.get)
	
			@Class2.0
			D=M
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@Class2.1
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
	