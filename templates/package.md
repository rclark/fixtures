<!--
MIT License

Copyright (c) 2019 Jeff Principe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
-->

{{- template "import" . -}}
{{- spacer -}}

{{- if len .Doc.Blocks -}}
	{{- template "doc" .Doc -}}
	{{- spacer -}}
{{- end -}}

{{- range (iter .Examples) -}}
	{{- template "example" .Entry -}}
	{{- spacer -}}
{{- end -}}

{{- header (add .Level 1) "Index" -}}
{{- spacer -}}

{{- template "index" . -}}

{{- if len .Consts -}}
	{{- spacer -}}

	{{- header (add .Level 1) "Constants" -}}
	{{- spacer -}}

	{{- range (iter .Consts) -}}
		{{- template "value" .Entry -}}
		{{- if (not .Last) -}}{{- spacer -}}{{- end -}}
	{{- end -}}

{{- end -}}

{{- if len .Vars -}}
	{{- spacer -}}

	{{- header (add .Level 1) "Variables" -}}
	{{- spacer -}}

	{{- range (iter .Vars) -}}
		{{- template "value" .Entry -}}
		{{- if (not .Last) -}}{{- spacer -}}{{- end -}}
	{{- end -}}

{{- end -}}

{{- if len .Funcs -}}
	{{- spacer -}}

	{{- range (iter .Funcs) -}}
		{{- template "func" .Entry -}}
		{{- if (not .Last) -}}{{- spacer -}}{{- end -}}
	{{- end -}}
{{- end -}}

{{- if len .Types -}}
	{{- spacer -}}

	{{- range (iter .Types) -}}
		{{- template "type" .Entry -}}
		{{- if (not .Last) -}}{{- spacer -}}{{- end -}}
	{{- end -}}
{{- end -}}
