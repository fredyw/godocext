godocext
========

A small tool written in Go to do a global search for the whole Go API

How to build
------------
`go build godocext.go`

Usage
-----
    Usage: ./godocext
      -f=false: show functions only
      -h=false: show help
      -m=false: show methods only
      -t=false: show types only`

Sample output
-------------
### Filter by methods only ###
    archive/zip: func (w *Writer) Create(name string) (io.Writer, error)
    archive/zip: func (w *Writer) CreateHeader(fh *FileHeader) (io.Writer, error)
    ...
    bufio: func (b *Reader) Buffered() int
    bufio: func (b *Reader) Peek(n int) ([]byte, error)
    ...
    bytes: func (b *Buffer) Len() int
    bytes: func (b *Buffer) Next(n int) []byte
    ...    

### Filter by functions only ###
    archive/zip: func NewReader(r io.ReaderAt, size int64) (*Reader, error)
    archive/zip: func NewWriter(w io.Writer) *Writer
    ...
    bufio: func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error)
    bufio: func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error)
    ...
    bytes: func Compare(a, b []byte) int
    bytes: func Contains(b, subslice []byte) bool
    ...

### Filter by types only ###
    archive/zip: type Reader struct { ... }
    archive/zip: type Writer struct { ... }
    ...
    bufio: type ReadWriter struct { ... }
    bufio: type Reader struct { ... }
    ...
    builtin: type ComplexType complex64
    builtin: type FloatType float32
    ...

Examples
--------
### Search for types that implement Read function ###
    Linux
    -----
    ./godocext -m | egrep "Read\(.*\)"
    godoc bufio Read (to get more detailed information)

    Windows (this example uses Powershell)
    --------------------------------------
    .\godocext -m | Select-String -Pattern "Read\(.*\)"
    godoc bufio Read (to get more detailed information)

    Build static output for faster search
    =====================================
    Linux
    -----
    ./godocext -m > methods.txt
    ./godocext -f > functions.txt
    ./godocext -t > types.txt
    grep "Read\(.*\)" methods.txt

    Windows (this example uses Powershell)
    --------------------------------------
    .\godocext -m > methods.txt
    .\godocext -f > functions.txt
    .\godocext -t > types.txt
    Get-Content methods.txt | Select-String -Pattern "Read\(.*\)"

