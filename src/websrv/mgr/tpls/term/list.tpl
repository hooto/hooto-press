
<table width="100%">
  <tr>
    <td>
      <form id="lps-infols-qry" action="#" class="form-inlines">
        <input id="lps-infols-qry-text" type="text"
          class="form-control l5smgr-query-input" 
          placeholder="Press Enter to Search" 
          value="">
      </form>
    </td>
    <td align="right">
      <button type="button" 
        class="btn btn-primary btn-sm" 
        onclick="l5sTerm.Set()">
        New Term
      </button>
    </td>
  </tr>
</table>
<div id="l5smgr-termls"></div>

<script id="l5smgr-termls-tpl" type="text/html">  
<table class="table table-hover">
  <thead>
    <tr>
      <th>ID</th>
      <th>Title</th>
      {[ if (it.model.type == "taxonomy") { ]}
      <th>Weight</th>
      {[ } ]}
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody>
  {[~it.items :v]}
    <tr>
      <td class="l5smgr-font-wfix">{[=v.id]}</td>
      <td>{[=v.title]}</td>
      {[ if (it.model.type == "taxonomy") { ]}
      <td>{[=v.weight]}</td>
      {[ } ]}
      <td>{[=v.created]}</td>
      <td>{[=v.updated]}</td>
      <td align="right"><a class="term-item" modname="{[=it.modname]}" modelid="{[=it.modelid]}" href="#{[=v.id]}">Edit</a>
    </tr>
  {[~]}
  </tbody>
</table>
</script>

<script type="text/javascript">

$("#l5smgr-termls").on("click", ".term-item", function() {
    var id = $(this).attr("href").substr(1);
    l5sTerm.Set($(this).attr("modname"), $(this).attr("modelid"), id);
});

</script>
