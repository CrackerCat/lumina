package main

import (
    "bufio"
    "bytes"
    "go/ast"
    "go/importer"
    "go/parser"
    "go/token"
    "go/types"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
)

type Metadata struct {
    Type            string
    ResponseType    string
}

type Tasks map[string]Metadata

func extractAutoGenTask(filename string) Tasks {
    tasks := Tasks{}
    f, err := os.Open(filename)
    if err != nil {
        log.Fatalf("unable to open %s, %v", filename, err)
    }
    r := bufio.NewReader(f)
    for {
        line, err := r.ReadString('\n')
        if err != nil { break }
        if !strings.HasPrefix(line, "//autogen ") { break }
        fields := strings.Fields(line[10:])
        if len(fields) > 0 {
            structName := fields[0]
            md := Metadata{}
            if len(fields) > 1 {
                md.Type = fields[1]
                if len(fields) > 2 {
                    md.ResponseType = fields[2]
                }
            }
            tasks[structName] = md
        }
    }
    f.Close()
    return tasks
}

type Context struct {
    depth       int
    counter     uint64
    r           bytes.Buffer
    w           bytes.Buffer
}

func (ctx *Context) unique() string {
    ctx.counter++
    return "v"+strconv.FormatUint(ctx.counter, 10)
}

func (ctx *Context) pad(w *bytes.Buffer) {
    write(w, strings.Repeat("    ", ctx.depth))
}

func (ctx *Context) both(s string) {
    ctx.read(s)
    ctx.write(s)
}

func (ctx *Context) read(s string) {
    ctx.pad(&ctx.r); write(&ctx.r, s); write(&ctx.r, "\n")
}

func (ctx *Context) write(s string) {
    ctx.pad(&ctx.w); write(&ctx.w, s); write(&ctx.w, "\n")
}

func (ctx *Context) readVar(t, v string) {
    ctx.read("if "+v+", err = read"+t+"(r); err != nil { return }")
}

func (ctx *Context) writeVar(t, v string) {
    ctx.write("if err = write"+t+"(w, "+v+"); err != nil { return }")
}

func (ctx *Context) readTypedVar(t, v, r string) {
    tmp := ctx.unique()
    ctx.read("var "+tmp+" "+strings.ToLower(t))
    ctx.read("if "+tmp+", err = read"+t+"(r); err != nil { return }")
    ctx.read(v+" = "+r+"("+tmp+")")
}

func (ctx *Context) writeTypedVar(t, v string) {
    ctx.write("if err = write"+t+"(w, "+strings.ToLower(t)+"("+v+")); err != nil { return }")
}

func (ctx *Context) readCall(s string) {
    ctx.read("if err = "+s+".ReadFrom(r); err != nil { return }")
}

func (ctx *Context) writeCall(s string) {
    ctx.write("if err = "+s+".WriteTo(w); err != nil { return }")
}

func resolveType(t types.Type) (typeName string, realType types.Type) {
    typeName = ""
    realType = t
    for {
        named, ok := realType.(*types.Named)
        if !ok { break }
        if typeName == "" { typeName = named.String() }
        realType = named.Underlying()
    }
    if typeName == "" {
        switch realType.(type) {
        case *types.Basic:
            typeName = realType.(*types.Basic).Name()
        case *types.Slice:
            // typeName is resolved later
        default:
            log.Fatalf("unable to resolve type %#v", t)
        }
    }
    return
}

func (ctx *Context) Process(varName string, varType types.Type) {
    ctx.depth += 1
    ctx.both("// Field "+varName)
    varTypeName, varRealType := resolveType(varType)
    switch varRealType.(type) {
    case *types.Basic:
        basicType := varRealType.(*types.Basic)
        ctx.both("// Basic "+basicType.Name())
        switch basicType.Kind() {
        case types.Uint64:
            if varTypeName == "uint64" {
                ctx.readVar("Uint64", varName)
                ctx.writeVar("Uint64", varName)
            } else {
                ctx.both("// Typed "+varTypeName)
                ctx.readTypedVar("Uint64", varName, varTypeName)
                ctx.writeTypedVar("Uint64", varName)
            }
        case types.Uint32:
            if varTypeName == "uint32" {
                ctx.readVar("Uint32", varName)
                ctx.writeVar("Uint32", varName)
            } else {
                ctx.both("// Typed "+varTypeName)
                ctx.readTypedVar("Uint32", varName, varTypeName)
                ctx.writeTypedVar("Uint32", varName)
            }
        case types.Int32:
            if varTypeName == "int32" {
                ctx.readVar("Int32", varName)
                ctx.writeVar("Int32", varName)
            } else {
                ctx.both("// Typed "+varTypeName)
                ctx.readTypedVar("Int32", varName, varTypeName)
                ctx.writeTypedVar("Int32", varName)
            }
        case types.Uint8:
            if varTypeName == "uint8" || varTypeName == "byte" {
                ctx.read("if "+varName+", err = r.ReadByte(); err != nil { return }")
                ctx.write("if err = w.WriteByte("+varName+"); err != nil { return }")
            } else {
                ctx.both("// Typed "+varTypeName)
                tmp := ctx.unique()
                ctx.read("var "+tmp+" uint8")
                ctx.read("if "+tmp+", err = "+"r.ReadByte(); err != nil { return }")
                ctx.read(varName+" = "+varTypeName+"("+tmp+")")
                ctx.write("if err = w.WriteByte(uint8("+varName+")); err != nil { return }")
            }
        case types.Bool:
            if varTypeName != "bool" {
                log.Fatal("typed bool is not supported")
            }
            ctx.readVar("Bool", varName)
            ctx.writeVar("Bool", varName)
        case types.String:
            if varTypeName != "string" {
                log.Fatal("typed string is not supported")
            }
            ctx.readVar("String", varName)
            ctx.writeVar("String", varName)
        default:
            log.Fatalf("basic type %#v is unsupported", basicType)
        }
    case *types.Struct:
        ctx.both("// Struct "+varTypeName)
        ctx.readCall(varName)
        ctx.writeCall(varName)
    case *types.Array:
        arrayType := varRealType.(*types.Array)
        if arrayType.Len() > 0x7FFFFFFF {
            log.Fatalf("array too long")
        }
        elemType := arrayType.Elem()
        elemTypeName, _ := resolveType(elemType)
        if elemTypeName == "" {
            log.Fatalf("array of slice is not supported directly")
        }
        length := strconv.FormatInt(arrayType.Len(), 10)
        ctx.both("// Array ["+length+"]"+elemTypeName)
        if elemBasicType, ok := elemType.(*types.Basic); ok && elemBasicType.Kind() == types.Byte {
            ctx.read("if err = readBytes(r, "+varName+"[:]); err != nil { return }")
            ctx.write("if err = writeBytes(w, "+varName+"[:]); err != nil { return }")
        } else {
            ctx.both("for i := uint32(0); i < "+length+"; i++ {")
            ctx.Process(varName+"[i]", elemType)
            ctx.both("}")
        }
    case *types.Slice:
        sliceType := varRealType.(*types.Slice)
        elemType := sliceType.Elem()
        elemTypeName, _ := resolveType(elemType)
        if elemTypeName == "" {
            log.Fatalf("slice of slice is not supported directly")
        }
        tmp := ctx.unique()
        ctx.both("// Slice []"+elemTypeName)
        ctx.read("var "+tmp+" uint32")
        ctx.readVar("Uint32", tmp)
        ctx.read(varName+" = make([]"+elemTypeName+", "+tmp+")")
        ctx.write("if len("+varName+") > 0x7FFFFFFF { err = errTooLong; return }")
        ctx.write("var "+tmp+" = uint32(len("+varName+"))")
        ctx.writeVar("Uint32", tmp)
        if elemBasicType, ok := elemType.(*types.Basic); ok && elemBasicType.Kind() == types.Byte {
            ctx.read("if err = readBytes(r, "+varName+"); err != nil { return }")
            ctx.write("if err = writeBytes(w, "+varName+"); err != nil { return }")
        } else {
            ctx.both("for i := uint32(0); i < "+tmp+"; i++ {")
            ctx.Process(varName+"[i]", elemType)
            ctx.both("}")
        }
    default:
        log.Fatalf("unknown type %#v", varRealType)
    }
    ctx.depth -= 1
}

func write(w io.Writer, s string) {
    _, err := io.WriteString(w, s)
    if err != nil {
        log.Fatal("unable to write: ", err)
    }
}

func writeFile(f *os.File, s string) {
    _, err := f.WriteString(s)
    if err != nil {
        log.Fatal("unable to write to ", f.Name(), ": ", err)
    }
}

func sourceFilter(fi os.FileInfo) bool {
    name := fi.Name()
    return !strings.HasSuffix(name, "_test.go") &&
           !strings.HasSuffix(name, "_autogen.go")
}

func main() {
    fset := token.NewFileSet()
    pkgs, err := parser.ParseDir(fset, ".", sourceFilter, parser.AllErrors)
    if err != nil {
        log.Fatal("parse error at ", err)
    }
    if len(pkgs) != 1 {
        log.Fatal("found multiple package")
    }
    var pkg *ast.Package
    for _, pkg = range pkgs {}
    conf := types.Config{
        IgnoreFuncBodies: true,
        Importer: importer.Default(),
    }
    files := []*ast.File{}
    allTasks := map[*ast.File]Tasks{}
    for filename, file := range pkg.Files {
        files = append(files, file)
        // TODO: type dependency inference, automatically add to autogen
        tasks := extractAutoGenTask(filename)
        if len(tasks) > 0 {
            allTasks[file] = tasks
        }
    }

    info := &types.Info{
        Types: map[ast.Expr]types.TypeAndValue{},
        // Defs:  map[*ast.Ident]types.Object{},
        Uses:  map[*ast.Ident]types.Object{},
    }
    _, err = conf.Check("", fset, files, info)
    if err != nil {
        log.Fatal("type check failed, ", err)
    }

    for filename, file := range pkg.Files {
        tasks, ok := allTasks[file]
        if !ok { continue }
        outputFilename := filename[:len(filename)-3]+"_autogen.go"
        output, err := os.Create(outputFilename)
        if err != nil {
            log.Fatal("unable to create ", outputFilename)
        }
        packageName := file.Name.Name
        writeFile(output, `// This file is automatically generated. DO NOT MODIFY!

package ` + packageName + `
`)
        for _, decl := range file.Decls {
            genDecl, _ := decl.(*ast.GenDecl)
            if genDecl == nil { continue }
            if genDecl.Tok != token.TYPE { continue }
            for _, spec := range genDecl.Specs {
                typeSpec := spec.(*ast.TypeSpec)
                astStructType, _ := typeSpec.Type.(*ast.StructType)
                if astStructType == nil { continue }
                structName := typeSpec.Name.Name
                task, ok := tasks[structName]
                if !ok { continue }
                delete(tasks, structName)
                if task.Type != "" {
                    writeFile(output, `
func (*` + structName + `) getType() PacketType {
    return ` + task.Type + `
}
`)
                }
                if task.ResponseType != "" {
                    writeFile(output, `
func (*` + structName + `) getResponseType() PacketType {
    return ` + task.ResponseType + `
}
`)
                }

                log.Print("struct ", structName)
                ctx := Context{}
                structType := info.TypeOf(astStructType).(*types.Struct)
                for i := 0; i < structType.NumFields(); i++ {
                    field := structType.Field(i)
                    if !field.Exported() { continue }
                    ctx.Process("this." + field.Name(), field.Type())
                }

                writeFile(output, `
func (this *` + structName + `) ReadFrom(r Reader) (err error) {
` + ctx.r.String() + `    return
}
`)
                writeFile(output, `
func (this *` + structName + `) WriteTo(w Writer) (err error) {
` + ctx.w.String() + `    return
}
`)
            }
        }
        if len(tasks) > 0 {
            log.Fatal("autogen struct is missing")
        }
        output.Close()
    }
}