package array

//go:generate gotemplate "bitbucket.org/7phs/fastgotext/wrapper/native/template/array" "IntArray(int, C.int, C.sizeof_int)"
//go:generate gotemplate "bitbucket.org/7phs/fastgotext/wrapper/native/template/matrix" "IntMatrix(int, C.int, C.sizeof_int)"
//go:generate gotemplate "bitbucket.org/7phs/fastgotext/wrapper/native/template/array" "FloatArray(float32, C.float, C.sizeof_float)"
//go:generate gotemplate "bitbucket.org/7phs/fastgotext/wrapper/native/template/matrix" "FloatMatrix(float32, C.float, C.sizeof_float)"
