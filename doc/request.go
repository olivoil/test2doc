package doc

import (
	"net/http"
	"strings"
	"text/template"
)

var (
	requestTmpl *template.Template
	requestFmt  = `{{if or .HasBody .HasHeader}}
+ Request {{if .HasContentType}}({{.Header.ContentType}}){{end}}{{with .Header}}

{{.Render}}{{end}}{{with .Body}}
{{.Render}}{{end}}{{end}}`
)

func init() {
	requestTmpl = template.Must(template.New("request").Parse(requestFmt))
}

const (
	TitleHeader = "X-Test2Doc-Title"
	DescriptionHeader = "X-Test2Doc-Description"
	FilterHeader = "X-Test2Doc-Filter"
)

type Request struct {
	Header   *Header
	Body     *Body
	Method   string
	Response *Response

	Title string
	Description string

	// TODO:
	// Attributes
	// Schema
}

func NewRequest(req *http.Request) (*Request, error) {
	desc := req.Header.Get(DescriptionHeader)
	req.Header.Del(DescriptionHeader)
	title := req.Header.Get(TitleHeader)
	req.Header.Del(TitleHeader)

	header := http.Header{}
	for k, _ := range req.Header {
		header.Set(k, req.Header.Get(k))
	}

	filter := req.Header.Get(FilterHeader)
	if filter != "" {
		f := strings.Split(filter, ";")
		for _, h := range f {
			if h != "" {
				header.Del(h)
			}
		}
	}
	req.Header.Del(FilterHeader)
	header.Del(FilterHeader)

	body1, body2, err := cloneBody(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = nopCloser{body1}

	b2bytes := body2.Bytes()
	contentType := req.Header.Get("Content-Type")

	return &Request{
		Header: NewHeader(header),
		Body:   NewBody(b2bytes, contentType),
		Method: req.Method,
		Description: desc,
		Title: title,
	}, nil
}

func (r *Request) Render() string {
	return render(requestTmpl, r)
}

func (r *Request) HasBody() bool {
	return r.Body != nil
}

func (r *Request) HasHeader() bool {
	return r.Header != nil && len(r.Header.DisplayHeader) > 0
}

func (r *Request) HasContentType() bool {
	return r.Header != nil && len(r.Header.ContentType) > 0
}
