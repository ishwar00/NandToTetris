	@261
	D=A
	@SP
	M=D
	@1234
	D=A
	@256
	M=D
	@261
	D=A
	@LCL
	M=D
	@265
	D=A
	@ARG
	M=D
	@THIS
	M=-1
	@THAT
	M=-1
	@Sys.init
	0;JMP
// function Main.fibonacci 0
(Main.fibonacci)
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
// push constant 2
	@2
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// lt
	@SP
	M=M-1
	A=M
	D=M
	A=A-1
	D=M-D
	@TRUE_0
	D;JLT
	D=0
	@DONE_0
	0;JMP
(TRUE_0)
	D=-1
(DONE_0)
	@SP
	A=M
	A=A-1
	M=D
// if-goto IF_TRUE
	@SP
	M=M-1
	A=M
	D=M
	@Main.fibonacci$IF_TRUE
	D;JNE
// goto IF_FALSE
	@Main.fibonacci$IF_FALSE
	0;JMP
// label IF_TRUE
(Main.fibonacci$IF_TRUE)
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
// return
	@LCL
	D=M
	@R13
	M=D
	@R13
	D=M
	@5
	A=D-A
	D=M
	@R14
	M=D
	@SP
	M=M-1
	A=M
	D=M
	@ARG
	A=M
	M=D
	@ARG
	M=M+1
	D=M
	@SP
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@THAT
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@THIS
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@ARG
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@LCL
	M=D
	@R14
	A=M
	0;JMP
// label IF_FALSE
(Main.fibonacci$IF_FALSE)
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
// push constant 2
	@2
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
// call Main.fibonacci 1
	@Main.fibonacci$ret.1
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@LCL
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@ARG
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THIS
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THAT
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@SP
	D=M
	@5
	D=D-A
	@1
	D=D-A
	@ARG
	M=D
	@SP
	D=M
	@LCL
	M=D
	@Main.fibonacci
	0;JMP
(Main.fibonacci$ret.1)
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
// call Main.fibonacci 1
	@Main.fibonacci$ret.2
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@LCL
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@ARG
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THIS
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THAT
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@SP
	D=M
	@5
	D=D-A
	@1
	D=D-A
	@ARG
	M=D
	@SP
	D=M
	@LCL
	M=D
	@Main.fibonacci
	0;JMP
(Main.fibonacci$ret.2)
// add
	@SP
	M=M-1
	A=M
	D=M
	A=A-1
	M=D+M
// return
	@LCL
	D=M
	@R13
	M=D
	@R13
	D=M
	@5
	A=D-A
	D=M
	@R14
	M=D
	@SP
	M=M-1
	A=M
	D=M
	@ARG
	A=M
	M=D
	@ARG
	M=M+1
	D=M
	@SP
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@THAT
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@THIS
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@ARG
	M=D
	@R13
	M=M-1
	A=M
	D=M
	@LCL
	M=D
	@R14
	A=M
	0;JMP
// function Sys.init 0
(Sys.init)
// push constant 4
	@4
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
// call Main.fibonacci 1
	@Sys.init$ret.1
	D=A
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@LCL
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@ARG
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THIS
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THAT
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@SP
	D=M
	@5
	D=D-A
	@1
	D=D-A
	@ARG
	M=D
	@SP
	D=M
	@LCL
	M=D
	@Main.fibonacci
	0;JMP
(Sys.init$ret.1)
// label WHILE
(Sys.init$WHILE)
// goto WHILE
	@Sys.init$WHILE
	0;JMP
