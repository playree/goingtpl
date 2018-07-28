{{define "child02-1"}}
<div>child02.tpl - 1</div>
{{end}}

{{define "child02-2"}}
<div>
    child02.tpl - 2
    <div style="padding-left: 2rem;">
        {{template "child03-2" .}}{{include "parts/child03.tpl"}}
    </div>
</div>
{{end}}