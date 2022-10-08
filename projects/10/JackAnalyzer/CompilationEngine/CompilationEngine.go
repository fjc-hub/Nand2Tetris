package compilationengine

import "bufio"

type CompilationEngine struct {
	output *bufio.Writer
}

func Constructor(output *bufio.Writer) CompilationEngine {
	return CompilationEngine{
		output: output,
	}
}

func (c *CompilationEngine) CompileClass() {

}

func (c *CompilationEngine) CompileClassVarDec() {

}

func (c *CompilationEngine) CompileSubroutine() {

}

func (c *CompilationEngine) CompileParameterList() {

}

func (c *CompilationEngine) CompileVarDec() {

}

func (c *CompilationEngine) CompileStatements() {

}

func (c *CompilationEngine) CompileDo() {

}

func (c *CompilationEngine) CompileLet() {

}

func (c *CompilationEngine) CompileWhile() {

}

func (c *CompilationEngine) CompileReturn() {

}

func (c *CompilationEngine) CompileIf() {

}

func (c *CompilationEngine) CompileExpression() {

}

func (c *CompilationEngine) CompileTerm() {

}

func (c *CompilationEngine) CompileExpressionList() {

}

func (c *CompilationEngine) Flush() {
	c.output.Flush()
}
