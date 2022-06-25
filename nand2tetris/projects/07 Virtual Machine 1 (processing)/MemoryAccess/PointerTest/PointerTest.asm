// push constant 3030
	@3030
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// pop pointer 0
	@SP
	M=M-1
	A=M
	D=M
	@THIS
	M=D
// push constant 3040
	@3040
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// pop pointer 1
	@SP
	M=M-1
	A=M
	D=M
	@THAT
	M=D
// push constant 32
	@32
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// pop this 2
	@THIS
	D=M
	@2
	D=D+A
	@R13
	M=D
	@SP
	M=M-1
	A=M
	D=M
	@R13
	A=M
	M=D
// push constant 46
	@46
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// pop that 6
	@THAT
	D=M
	@6
	D=D+A
	@R13
	M=D
	@SP
	M=M-1
	A=M
	D=M
	@R13
	A=M
	M=D
// push pointer 0
	@THIS
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// push pointer 1
	@THAT
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// add
	@SP
	M=M-1
	A=M
	D=M
	A=A-1
	M=D+M
	D=M-D
	M=M-D
	M=D+M
// push this 2
	@THIS
	D=M
	@2
	A=D+A
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// sub
	@SP
	M=M-1
	A=M
	D=M
	A=A-1
	M=M-D
// push that 6
	@THAT
	D=M
	@6
	A=D+A
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// add
	@SP
	M=M-1
	A=M
	D=M
	A=A-1
	M=D+M
	D=M-D
	M=M-D
	M=D+M
(END)
	@END
	0;JMP
