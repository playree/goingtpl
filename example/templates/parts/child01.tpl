{{define "child01"}}
<div>
    child01.tpl {{.Date}} {{.Time}}
    <div style="padding-left: 2rem;">
        {{template "child03-1"}}{{include "parts/child03.tpl"}}
    </div>
</div>
{{end}}