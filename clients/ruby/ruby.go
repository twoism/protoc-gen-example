package ruby

import (
	"fmt"
	"text/template"

	"strings"

	"bytes"
	"log"

	"unicode"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

var ClientTemplate = `require_relative "./protos/{{.BaseName}}"

class {{.ClassName}}Client
{{ with .Methods }}{{range .}}  def {{ .Name }}(request)
    return {{.Output}}.new
  end{{ end }}{{ end }}
end
`

var TestTemplate = `require 'minitest/autorun'
require_relative './{{.BaseName}}_client'
require_relative './protos/{{.BaseName}}'

include {{.Namespace}}

class TestBlog < Minitest::Test

  def setup
    @client = {{.ClassName}}Client.new
  end

{{ with .Methods }}{{range .}}  def test_{{ .Name }}
    assert_equal @client.{{.Name}}({{.Input}}.new), {{.Output}}.new
  end{{ end }}{{ end }}

end`

// Method impl for python methods
type Method struct {
	m *descriptor.MethodDescriptorProto
}

func (m *Method) Name() string {
	return ToSnake(m.m.GetName())
}

func (m *Method) Input() string {
	return TrimType(m.m.GetInputType())
}

func (m *Method) Output() string {
	return TrimType(m.m.GetOutputType())
}

type RubyClient struct {
	srv  *descriptor.ServiceDescriptorProto
	file *descriptor.FileDescriptorProto

	methods []*Method
}

func New(srv *descriptor.ServiceDescriptorProto, file *descriptor.FileDescriptorProto) (c *RubyClient) {
	c = &RubyClient{
		srv:  srv,
		file: file,
	}

	for _, m := range srv.GetMethod() {
		c.AppendMethod(m)
	}

	return
}

func (c *RubyClient) AppendMethod(m *descriptor.MethodDescriptorProto) {
	c.methods = append(c.methods, &Method{m})
}

func (c *RubyClient) ClassName() string {
	return *c.srv.Name
}

func (c *RubyClient) Methods() []*Method {
	return c.methods
}

func (c *RubyClient) BaseName() string {
	parts := strings.Split(c.file.GetName(), "/")

	if len(parts) == 0 {
		return ""
	}

	return strings.Replace(parts[len(parts)-1], ".proto", "", 1)

}

func (c *RubyClient) Namespace() string {
	parts := strings.Split(c.file.GetPackage(), ".")

	if len(parts) == 0 {
		return ""
	}

	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	return strings.Join(parts, "::")
}

func (c *RubyClient) FileName() *string {
	return proto.String(
		fmt.Sprintf(
			"%s/%s_client.rb", strings.Replace(c.file.GetPackage(), ".", "/", -1),
			c.BaseName(),
		),
	)
}

func (c *RubyClient) TestFileName() *string {
	return proto.String(
		fmt.Sprintf(
			"%s/%s_client_test.rb", strings.Replace(c.file.GetPackage(), ".", "/", -1),
			c.BaseName(),
		),
	)
}

func (c *RubyClient) Content() *string {
	t := template.New("ruby-tpl")
	tpl, err := t.Parse(ClientTemplate)

	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	tpl.Execute(&buf, c)

	return proto.String(buf.String())
}

func (c *RubyClient) TestContent() *string {
	t := template.New("ruby-test-tpl")
	tpl, err := t.Parse(TestTemplate)

	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	tpl.Execute(&buf, c)

	return proto.String(buf.String())
}

func (c *RubyClient) File() *plugin_go.CodeGeneratorResponse_File {
	return &plugin_go.CodeGeneratorResponse_File{
		Content: c.Content(),
		Name:    c.FileName(),
	}
}

func (c *RubyClient) TestFile() *plugin_go.CodeGeneratorResponse_File {
	return &plugin_go.CodeGeneratorResponse_File{
		Content: c.TestContent(),
		Name:    c.TestFileName(),
	}
}

// ToSnake snake_cases CamelCase strings
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// TrimType returns the last part of a dotted namespace
func TrimType(it string) string {
	parts := strings.Split(it, ".")

	if len(parts) == 0 {
		return ""
	}

	if len(parts) == 1 {
		return parts[0]
	}

	return parts[len(parts)-1]
}
