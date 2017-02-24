//go:generate go run ../ttgen/main.go -filename templates/index.gtt -out template.go -package test

package test

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

const coreTemplate = `<!DOCTYPE html>
<html>
<body>
    <h1>{{ .Header }}</h1>
    <hr>
    {{range .Users }}
        <p>{{ .Name }} {{ .Lastname }}: {{ .Age }}</p>
    {{ end }}
</<body>
</html>`

var coreTpl *template.Template
var BenchmarkUsers []User

func init() {
	var err error
	coreTpl, err = template.New("test").Parse(coreTemplate)
	if err != nil {
		panic(err)
	}

	BenchmarkUsers = make([]User, 1000)
	for i := 0; i < 1000; i++ {
		BenchmarkUsers[i] = User{"Name", "Lastname", 32}
	}
}

func TestProcessTest(t *testing.T) {
	buf := &bytes.Buffer{}
	ProcessTest(buf, "<Header>", []User{
		{"Ivan", "Sidorov", 20},
		{"Petr", "Ivanov", 30},
		{"James", "Bond", 40},
		{"John", "Connor", 50},
		{"Sara", "Connor", 60},
	})

	assert.Contains(t, buf.String(), "<h1>&lt;Header&gt;</h1>")
	assert.Contains(t, buf.String(), "<p>James Bond: 40</p>")

	//buf.WriteTo(os.Stdout)
}

func BenchmarkProcessTest(b *testing.B) {
	buf := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		ProcessTest(buf, "<Header>", BenchmarkUsers)
	}
}

func BenchmarkCoreTemplate(b *testing.B) {
	buf := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		coreTpl.Execute(buf, struct {
			Header string
			Users  []User
		}{"<Header>", BenchmarkUsers})
	}
}
