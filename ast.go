package template

import (
	"io"
	"strings"
)

type iAstNode interface {
	WriteGo(io.Writer, *GenGoOpts)
}

type GenGoOpts struct {
	PackageName string
	Imports     []string
	TemplateId  string
}

type astList struct {
	children []iAstNode
}

func (n *astList) WriteGo(w io.Writer, opts *GenGoOpts) {
	for _, child := range n.children {
		if child != nil {
			child.WriteGo(w, opts)
		}
	}
}

func (n *astList) Add(node ...iAstNode) {
	n.children = append(n.children, node...)
}

type astFile struct {
	header, macroses iAstNode
}

func (n *astFile) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "package "+opts.PackageName+"\n")

	if n.header == nil {
		n.header = &astHeader{}
	}

	n.header.WriteGo(w, opts)

	if n.macroses != nil {
		n.macroses.WriteGo(w, opts)
	}
}

type astHeader struct {
	imports *astList
}

func (n *astHeader) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "import (\n")
	io.WriteString(w, "\"fmt\"\n")
	io.WriteString(w, "\"io\"\n")
	io.WriteString(w, "\"github.com/go-qbit/template/filter\"\n")
	if n.imports != nil {
		for _, pkgName := range n.imports.children {
			pkgName.WriteGo(w, opts)
		}
	}
	io.WriteString(w, ")\n")
}

type astImport struct {
	pkgName string
}

func (n *astImport) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, n.pkgName+"\n")
}

type astTemplate struct {
	name    string
	vars    *astList
	wrapper *astUseWrapper
	body    *astList
}

func (n *astTemplate) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "func Process"+n.name+"(w io.Writer")

	for _, v := range n.vars.children {
		io.WriteString(w, ", ")
		v.WriteGo(w, opts)
	}

	io.WriteString(w, ") {\n")

	if n.wrapper != nil {
		io.WriteString(w, "Wrapper"+n.wrapper.name+"(w, func() {\n")
	}

	if n.body != nil {
		n.body.WriteGo(w, opts)
	}

	if n.wrapper != nil {
		io.WriteString(w, "}")
		for _, p := range n.wrapper.params.children {
			io.WriteString(w, ", ")
			p.WriteGo(w, opts)
		}
		io.WriteString(w, ")\n")
	}

	io.WriteString(w, "}\n\n")
}

type astUseWrapper struct {
	name   string
	params *astList
}

func (n *astUseWrapper) WriteGo(w io.Writer, opts *GenGoOpts) {}

type astWrapper struct {
	name string
	vars *astList
	body iAstNode
}

func (n *astWrapper) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "func Wrapper"+n.name+"(w io.Writer, tplClbF func()")
	for _, v := range n.vars.children {
		io.WriteString(w, ", ")
		v.WriteGo(w, opts)
	}
	io.WriteString(w, ") {\n")

	if n.body != nil {
		n.body.WriteGo(w, opts)
	}

	io.WriteString(w, "}\n\n")
}

type astWriteContent struct{}

func (n *astWriteContent) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "tplClbF()\n")
}

type astVariableDef struct {
	vName, vType string
}

func (n *astVariableDef) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, n.vName+" "+n.vType)
}

type astWriteValue struct {
	value iAstNode
}

func (n *astWriteValue) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "io.WriteString(w, fmt.Sprint(")
	n.value.WriteGo(w, opts)
	io.WriteString(w, "))\n")
}

type astWriteString struct {
	value iAstNode
}

func (n *astWriteString) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "io.WriteString(w, ")
	n.value.WriteGo(w, opts)
	io.WriteString(w, ")\n")
}

type astString struct {
	value string
}

func (n *astString) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, strings.Replace(n.value, "\n", `\n`, -1))
}

type astValue struct {
	name string
}

func (n *astValue) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, n.name)
}

type astFilter struct {
	name  string
	value iAstNode
}

func (n *astFilter) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "filter."+n.name+"(")
	n.value.WriteGo(w, opts)
	io.WriteString(w, ")")
}

type astLoop struct {
	localVariable string
	loopVariable  iAstNode
	body          iAstNode
}

func (n *astLoop) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "for _, "+n.localVariable+":= range ")
	n.loopVariable.WriteGo(w, opts)
	io.WriteString(w, "{\n")
	n.body.WriteGo(w, opts)
	io.WriteString(w, "}\n")
}

type astCondition struct {
	condition iAstNode
	ifBody    iAstNode
	elseBody  iAstNode
}

func (n *astCondition) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "if ")
	n.condition.WriteGo(w, opts)
	io.WriteString(w, " {\n")
	n.ifBody.WriteGo(w, opts)
	io.WriteString(w, "}")

	if n.elseBody != nil {
		io.WriteString(w, " else {\n")
		n.elseBody.WriteGo(w, opts)
		io.WriteString(w, "}")
	}

	io.WriteString(w, "\n")
}
