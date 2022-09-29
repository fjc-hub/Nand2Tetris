
			@17
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@17
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@eq.true.0
			D;JEQ
			@SP
			A=M-1
			M=0
			@eq.skip.0
			0; JMP
			(eq.true.0)
			@SP
			A=M-1
			M=-1
			(eq.skip.0)
		
			@17
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@16
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@eq.true.1
			D;JEQ
			@SP
			A=M-1
			M=0
			@eq.skip.1
			0; JMP
			(eq.true.1)
			@SP
			A=M-1
			M=-1
			(eq.skip.1)
		
			@16
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@17
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@eq.true.2
			D;JEQ
			@SP
			A=M-1
			M=0
			@eq.skip.2
			0; JMP
			(eq.true.2)
			@SP
			A=M-1
			M=-1
			(eq.skip.2)
		
			@892
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@891
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@lt.true.3
			D;JLT
			@SP
			A=M-1
			M=0
			@lt.skip.3
			0; JMP
			(lt.true.3)
			@SP
			A=M-1
			M=-1
			(lt.skip.3)
		
			@891
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@892
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@lt.true.4
			D;JLT
			@SP
			A=M-1
			M=0
			@lt.skip.4
			0; JMP
			(lt.true.4)
			@SP
			A=M-1
			M=-1
			(lt.skip.4)
		
			@891
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@891
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@lt.true.5
			D;JLT
			@SP
			A=M-1
			M=0
			@lt.skip.5
			0; JMP
			(lt.true.5)
			@SP
			A=M-1
			M=-1
			(lt.skip.5)
		
			@32767
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@32766
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@gt.true.6
			D;JGT
			@SP
			A=M-1
			M=0
			@gt.skip.6
			0; JMP
			(gt.true.6)
			@SP
			A=M-1
			M=-1
			(gt.skip.6)
		
			@32766
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@32767
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@gt.true.7
			D;JGT
			@SP
			A=M-1
			M=0
			@gt.skip.7
			0; JMP
			(gt.true.7)
			@SP
			A=M-1
			M=-1
			(gt.skip.7)
		
			@32766
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@32766
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D
			@gt.true.8
			D;JGT
			@SP
			A=M-1
			M=0
			@gt.skip.8
			0; JMP
			(gt.true.8)
			@SP
			A=M-1
			M=-1
			(gt.skip.8)
		
			@57
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@31
			D=A
			@SP		// (*sp) = D
			A=M
			M=D
			@SP		// sp++
			M=M+1 
	
			@53
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
		
			@112
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
		
			@SP		// A=sp
			A=M-1	// A=(*sp)-1
			M=-M
		
			@SP		// A=sp
			AM=M-1	// A=(*sp)-1; (*sp)--
			D=M		// D=*(*sp)
			A=A-1	// A=(*sp)-1
			M=M&D
		
			@82
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
			M=M|D
		
			@SP
			A=M-1
			M=!M
		