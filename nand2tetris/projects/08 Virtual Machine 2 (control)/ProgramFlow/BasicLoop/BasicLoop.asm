// push constant 0
	@0
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// pop local 0
	@LCL
	D=M
	@0
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
// label LOOP_START
(null$LOOP_START)
// push argument 0
	@ARG
	D=M
	@0
	A=D+A
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// push local 0
	@LCL
	D=M
	@0
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
// pop local 0
	@LCL
	D=M
	@0
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
// push argument 0
	@ARG
	D=M
	@0
	A=D+A
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// push constant 1
	@1
	D=A
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
// pop argument 0
	@ARG
	D=M
	@0
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
// push argument 0
	@ARG
	D=M
	@0
	A=D+A
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
// if-goto LOOP_START
	@SP
	M=M-1
	A=M
	D=M
	@null$LOOP_START
	D;JNE
// push local 0
	@LCL
	D=M
	@0
	A=D+A
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
