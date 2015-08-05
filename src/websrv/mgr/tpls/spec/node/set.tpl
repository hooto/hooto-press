<style>
.l5smgr-spec-node-field-attr-item td {
  padding: 0 2px 4px;
}
</style>

<div id="l5smgr-spec-nodeset-alert"></div>

<div id="l5smgr-spec-nodeset">

  <input type="hidden" name="modname" value="{[=it._modname]}">

  {[? it.meta.name]}
  <input type="hidden" name="name" value="{[=it.meta.name]}">
  {[??]}
  <div class="form-group">
    <label>Name</label>
    <input type="text" class="form-control input-sm" name="name" 
      placeholder="Node Name" value="{[=it.meta.name]}">
  </div>
  {[?]}

  <div class="form-group">
    <label>Title</label>
    <input type="text" class="form-control input-sm" name="title" 
      placeholder="Title" value="{[=it.title]}">
  </div>

  <div class="form-group">
    <label>Fields</label>
    <div>
      <table class="table table-condensed" width="100%">
      <thead>
        <tr>
          <th>Name</th>
          <th>Title</th>
          <th>Type</th>
          <th>Length</th>
          <th>Index Type</th>
          <th>Extended attributes</th>
          <th></th>
        </tr>
      </thead>
      <tbody id="l5smgr-spec-node-fields">
        {[~it.fields :v]}
        <tr id="field-seq-{[=v._seqid]}" class="l5smgr-spec-node-field-item">
          <td><input type="text" class="form-control input-sm" name="field_name" size="10" value="{[=v.name]}" readonly></td>
          <td><input type="text" class="form-control input-sm" name="field_title" size="20" value="{[=v.title]}"></td>
          <td>
            <select class="form-control input-sm" name="field_type">
            {[~it._field_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.type) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
          <td><input class="form-control input-sm" type="text" name="field_length" size="3" value="{[=v.length]}"></td>
          <td>
            <select class="form-control input-sm" name="field_index_type">
            {[~it._field_idx_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.indexType) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
          <td>
            <table class="l5smgr-spec-node-field-attrs">
              {[~v.attrs :atv]}
              <tr class="l5smgr-spec-node-field-attr-item">
                <td><input type="text" class="form-control input-sm" name="field_attr_key" size="8" value="{[=atv.key]}"></td>
                <td><input type="text" class="form-control input-sm" name="field_attr_value" size="16" value="{[=atv.value]}"></td>
              </tr>
              {[~]}
            </table>
          </td>
          <td>
            <button class="btn btn-default btn-sm" onclick="l5sSpec.NodeSetFieldAttrAppend('{[=v._seqid]}')">+ Attribute</button>
          </td>
        </tr>
        {[~]}
      </tbody>
      </table>
    </div>
  </div>

  <div class="form-group">
    <label>Terms</label>
    <div>
      <table class="table table-condensed" width="100%">
      <thead>
        <tr>
          <th>Name</th>
          <th>Title</th>
          <th>Type</th>
        </tr>
      </thead>
      <tbody id="l5smgr-spec-node-terms">
        {[~it.terms :v]}
        <tr id="field-seq-{[=v._seqid]}" class="l5smgr-spec-node-term-item">
          <td><input type="text" class="form-control input-sm" name="term_name" size="20" value="{[=v.meta.name]}" readonly></td>
          <td><input type="text" class="form-control input-sm" name="term_title" size="30" value="{[=v.title]}"></td>
          <td>
            <select class="form-control input-sm" name="term_type">
            {[~it._term_typedef :fv]}
            <option value="{[=fv.type]}" {[ if (fv.type == v.type) { ]}selected{[ } ]}>{[=fv.name]}</option>
            {[~]}
            </select>
          </td>
        {[~]}
      </table>
    </div>
  </div>

</div>

<script id="l5smgr-spec-node-field-item-tpl" type="text/html">
  <tr id="field-seq-{[=it._seqid]}" class="l5smgr-spec-node-field-item">
    <td><input type="text" class="form-control input-sm" name="field_name" size="10" value=""></td>
    <td><input type="text" class="form-control input-sm" name="field_title" size="16" value=""></td>
    <td>
      <select class="form-control input-sm" name="field_type">
      {[~it._field_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == it._type) { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
    <td><input type="text" class="form-control input-sm" name="field_length" size="5" value="0"></td>
    <td>
      <select class="form-control input-sm" name="field_index_type">
      {[~it._field_idx_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == it._indexType) { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
    <td>
      <table class="l5smgr-spec-node-field-attrs"></table>
    </td>
    <td>
      <button class="btn btn-default btn-sm" onclick="l5sSpec.NodeSetFieldAttrAppend('{[=it._seqid]}')">+ Attribute</button>
    </td>
  </tr>
</script>

<script id="l5smgr-spec-node-field-attr-item-tpl" type="text/html">
  <tr class="l5smgr-spec-node-field-attr-item">
    <td><input type="text" class="form-control input-sm" name="field_attr_key" size="8" value=""></td>
    <td><input type="text" class="form-control input-sm" name="field_attr_value" size="12" value=""></td>
  </tr>
</script>

<script id="l5smgr-spec-node-term-item-tpl" type="text/html">
  <tr id="field-seq-{[=it._seqid]}" class="l5smgr-spec-node-term-item">
    <td><input type="text" class="form-control input-sm" name="term_name" size="20" value=""></td>
    <td><input type="text" class="form-control input-sm" name="term_title" size="30" value=""></td>
    <td>
      <select class="form-control input-sm" name="term_type">
      {[~it._term_typedef :fv]}
        <option value="{[=fv.type]}" {[ if (fv.type == it._type) { ]}selected{[ } ]}>{[=fv.name]}</option>
      {[~]}
      </select>
    </td>
  </tr>
</script>
