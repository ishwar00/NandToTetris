// push constant 111
	@111
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// push constant 333
	@333
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// push constant 888
	@888
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// pop static 8
	@SP
	M=M-1
	A=M
	D=M
	@StaticTest.8
	M=D
// popp static 3
