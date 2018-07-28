{{define "child03-1"}}
<i>child03.tpl - 1</i>
{{end}}

{{define "child03-2"}}
<i>child03.tpl - 2 </i><br>
Arguments = {{.Date}} {{.Time}}<br>
Func now = {{now}}<br>
Func repeat = {{repeat "A" 5}}
{{end}}