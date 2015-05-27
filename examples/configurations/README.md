Configuration files are plain YAML files used to pass parameters when templates are applied or "run". Their structure depends on each template. Configuration files may contain placeholders for additional command line parameters:

```
input: $0
output: $1
```

$0 and $1 refer to parameters that can be specified when the engine is invoked from the terminal. The main use of this feature is to allow specifying input and output files.
