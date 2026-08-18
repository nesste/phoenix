package devrepo

var emptyArgs = []byte(`{"type":"object","additionalProperties":false}`)

var pathArgs = []byte(`{
  "type":"object","additionalProperties":false,"required":["path"],
  "properties":{"path":{"type":"string","minLength":1}}
}`)

var editArgs = []byte(`{
  "type":"object","additionalProperties":false,"required":["path","old","new"],
  "properties":{
    "path":{"type":"string","minLength":1},
    "old":{"type":"string","minLength":1},
    "new":{"type":"string"}
  }
}`)

var queryArgs = []byte(`{
  "type":"object","additionalProperties":false,"required":["query"],
  "properties":{"query":{"type":"string","minLength":1}}
}`)

var focusArgs = []byte(`{
  "type":"object","additionalProperties":false,"required":["test"],
  "properties":{"test":{"type":"string","minLength":1}}
}`)

var commitArgs = []byte(`{
  "type":"object","additionalProperties":false,"required":["message"],
  "properties":{"message":{"type":"string","minLength":1,"maxLength":200}}
}`)

var readResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["path","content","digest"],
  "properties":{"path":{"type":"string"},"content":{"type":"string"},"digest":{"type":"string"}}
}`)

var editResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["path","changed","digest"],
  "properties":{"path":{"type":"string"},"changed":{"type":"boolean"},"digest":{"type":"string"}}
}`)

var commandResultSchema = []byte(`{
  "type":"object","additionalProperties":false,"required":["exit_code","stdout","stderr"],
  "properties":{"exit_code":{"type":"integer"},"stdout":{"type":"string"},"stderr":{"type":"string"}}
}`)

var findResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["matches"],
  "properties":{"matches":{"type":"array","items":{"type":"string"}}}
}`)

var repoStatusResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["branch","clean"],
  "properties":{"branch":{"type":"string"},"clean":{"type":"boolean"}}
}`)

var testsRunResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["exit_code","passed","failed","stdout","stderr"],
  "properties":{
    "exit_code":{"type":"integer"},"passed":{"type":"integer","minimum":0},"failed":{"type":"integer","minimum":0},
    "stdout":{"type":"string"},"stderr":{"type":"string"}
  }
}`)

var testsListResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["tests"],
  "properties":{"tests":{"type":"array","items":{"type":"string"}}}
}`)

var testsFocusResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["test","passed","exit_code","stdout","stderr"],
  "properties":{
    "test":{"type":"string"},"passed":{"type":"boolean"},"exit_code":{"type":"integer"},
    "stdout":{"type":"string"},"stderr":{"type":"string"}
  }
}`)

var gitStatusResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["branch","clean","output"],
  "properties":{"branch":{"type":"string"},"clean":{"type":"boolean"},"output":{"type":"string"}}
}`)

var gitDiffResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["diff"],
  "properties":{"diff":{"type":"string"}}
}`)

var gitCommitResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["committed","stdout","stderr"],
  "properties":{"committed":{"type":"boolean"},"stdout":{"type":"string"},"stderr":{"type":"string"}}
}`)

var recallResult = []byte(`{
  "type":"object","additionalProperties":false,"required":["episodes"],
  "properties":{"episodes":{"type":"array","items":{
    "type":"object","additionalProperties":false,"required":["episode_id","act_id"],
    "properties":{"episode_id":{"type":"string"},"act_id":{"type":"string"}}
  }}}
}`)
