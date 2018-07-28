<!DOCTYPE html>
<html><body>
    <div>parent.tpl</div>
    <div style="padding-left: 2rem;">
        {{template "child01"}}{{include "parts/child01.tpl"}}
        {{template "child02-1"}}{{include "parts/child02.tpl"}}
        {{template "child02-2" .}}
    </div>
</body></html>