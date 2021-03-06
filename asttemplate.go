package template

import "io"

type astTemplate struct {
	name    string
	vars    *astList
	wrapper *astUseWrapper
	body    *astList
}

func (n *astTemplate) GetImports() []string {
	res := []string{"context", "io"}
	if n.body != nil {
		res = append(res, n.body.GetImports()...)
	}

	return res
}

func (n *astTemplate) GetStrings() []string {
	res := []string{}
	if n.body != nil {
		res = append(res, n.body.GetStrings()...)
	}

	return res
}

func (n *astTemplate) WriteGo(w io.Writer, opts *GenGoOpts) {
	io.WriteString(w, "func Process"+n.name+"(ctx context.Context, w io.Writer")

	if n.vars != nil {
		for _, v := range n.vars.children {
			if v != nil {
				io.WriteString(w, ", ")
				v.WriteGo(w, opts)
			}
		}
	}

	io.WriteString(w, ") {\n")

	if n.wrapper != nil {
		if n.wrapper.pkgName != "" {
			io.WriteString(w, n.wrapper.pkgName+".")
		}
		io.WriteString(w, "Wrapper"+n.wrapper.name+"(ctx, w, func() {\n")
	}

	if n.body != nil {
		n.body.WriteGo(w, opts)
	}

	if n.wrapper != nil {
		io.WriteString(w, "}")
		if n.wrapper.params != nil {
			for _, p := range n.wrapper.params.children {
				io.WriteString(w, ", ")
				p.WriteGo(w, opts)
			}
		}
		io.WriteString(w, ")\n")
	}

	io.WriteString(w, "}\n\n")
}
