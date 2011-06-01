package walk

// import "os"
import "fmt"
import . "go/ast"
import "go/token"

type DummyNode struct {
    start token.Pos
    end   token.Pos
    Name  string
    List  []Node
}

func NewDummyNode(name string, start, end token.Pos, list []Node) *DummyNode {
    return &DummyNode{start: start, end: end, Name: name, List: list}
}
func (self *DummyNode) Pos() token.Pos { return self.start }
func (self *DummyNode) End() token.Pos { return self.end }

// Helper functions.
func walkIdentList(v Visitor, list []*Ident) {
    for _, x := range list {
        GoAST_Walk(v, x)
    }
}


func walkExprList(v Visitor, list []Expr) {
    if len(list) > 0 {
        GoAST_Walk(v,
            NewDummyNode("ExprList",
                list[0].Pos(),
                list[len(list)-1].End(),
                func() []Node {
                    nodes := make([]Node, 0, len(list))
                    for _, c := range list {
                        nodes = append(nodes, Node(c))
                    }
                    return nodes
                }()))
    }
}


func walkStmtList(v Visitor, list []Stmt) {
    for _, x := range list {
        GoAST_Walk(v, x)
    }
}


func walkDeclList(v Visitor, list []Decl) {
    for _, x := range list {
        GoAST_Walk(v, x)
    }
}


// I am forking the version of walk in the go stdlib. Why this crazyness?
// because I want to insert nodes and this is the cleanest way to do it.

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func GoAST_Walk(v Visitor, node Node) {
    if v = v.Visit(node); v == nil {
        return
    }

    // walk children
    // (the order of the cases matches the order
    // of the corresponding node types in ast.go)
    switch n := node.(type) {
    // Comments and fields
    case *Comment:
        // nothing to do
    case *DummyNode:
        for _, c := range n.List {
            GoAST_Walk(v, c)
        }

    case *CommentGroup:
        for _, c := range n.List {
            GoAST_Walk(v, c)
        }

    case *Field:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        walkIdentList(v, n.Names)
        GoAST_Walk(v, n.Type)
        if n.Tag != nil {
            GoAST_Walk(v, n.Tag)
        }
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *FieldList:
        for _, f := range n.List {
            GoAST_Walk(v, f)
        }

    // Expressions
    case *BadExpr, *Ident, *BasicLit:
        // nothing to do

    case *Ellipsis:
        if n.Elt != nil {
            GoAST_Walk(v, n.Elt)
        }

    case *FuncLit:
        GoAST_Walk(v, n.Type)
        GoAST_Walk(v, n.Body)

    case *CompositeLit:
        if n.Type != nil {
            GoAST_Walk(v, n.Type)
        }
        walkExprList(v, n.Elts)

    case *ParenExpr:
        GoAST_Walk(v, n.X)

    case *SelectorExpr:
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Sel)

    case *IndexExpr:
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Index)

    case *SliceExpr:
        GoAST_Walk(v, n.X)
        if n.Low != nil {
            GoAST_Walk(v, n.Low)
        }
        if n.High != nil {
            GoAST_Walk(v, n.High)
        }

    case *TypeAssertExpr:
        GoAST_Walk(v, n.X)
        if n.Type != nil {
            GoAST_Walk(v, n.Type)
        }

    case *CallExpr:
        GoAST_Walk(v, n.Fun)
        walkExprList(v, n.Args)

    case *StarExpr:
        GoAST_Walk(v, n.X)

    case *UnaryExpr:
        GoAST_Walk(v, n.X)

    case *BinaryExpr:
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Y)

    case *KeyValueExpr:
        GoAST_Walk(v, n.Key)
        GoAST_Walk(v, n.Value)

    // Types
    case *ArrayType:
        if n.Len != nil {
            GoAST_Walk(v, n.Len)
        }
        GoAST_Walk(v, n.Elt)

    case *StructType:
        GoAST_Walk(v, n.Fields)

    case *FuncType:
        GoAST_Walk(v, n.Params)
        if n.Results != nil {
            GoAST_Walk(v, n.Results)
        }

    case *InterfaceType:
        GoAST_Walk(v, n.Methods)

    case *MapType:
        GoAST_Walk(v, n.Key)
        GoAST_Walk(v, n.Value)

    case *ChanType:
        GoAST_Walk(v, n.Value)

    // Statements
    case *BadStmt:
        // nothing to do

    case *DeclStmt:
        GoAST_Walk(v, n.Decl)

    case *EmptyStmt:
        // nothing to do

    case *LabeledStmt:
        GoAST_Walk(v, n.Label)
        GoAST_Walk(v, n.Stmt)

    case *ExprStmt:
        GoAST_Walk(v, n.X)

    case *SendStmt:
        GoAST_Walk(v, n.Chan)
        GoAST_Walk(v, n.Value)

    case *IncDecStmt:
        GoAST_Walk(v, n.X)

    case *AssignStmt:
        walkExprList(v, n.Lhs)
        walkExprList(v, n.Rhs)

    case *GoStmt:
        GoAST_Walk(v, n.Call)

    case *DeferStmt:
        GoAST_Walk(v, n.Call)

    case *ReturnStmt:
        walkExprList(v, n.Results)

    case *BranchStmt:
        if n.Label != nil {
            GoAST_Walk(v, n.Label)
        }

    case *BlockStmt:
        walkStmtList(v, n.List)

    case *IfStmt:
        if n.Init != nil {
            GoAST_Walk(v, n.Init)
        }
        GoAST_Walk(v, n.Cond)
        GoAST_Walk(v, n.Body)
        if n.Else != nil {
            GoAST_Walk(v, n.Else)
        }

    case *CaseClause:
        walkExprList(v, n.List)
        walkStmtList(v, n.Body)

    case *SwitchStmt:
        if n.Init != nil {
            GoAST_Walk(v, n.Init)
        }
        if n.Tag != nil {
            GoAST_Walk(v, n.Tag)
        }
        GoAST_Walk(v, n.Body)

    case *TypeSwitchStmt:
        if n.Init != nil {
            GoAST_Walk(v, n.Init)
        }
        GoAST_Walk(v, n.Assign)
        GoAST_Walk(v, n.Body)

    case *CommClause:
        if n.Comm != nil {
            GoAST_Walk(v, n.Comm)
        }
        walkStmtList(v, n.Body)

    case *SelectStmt:
        GoAST_Walk(v, n.Body)

    case *ForStmt:
        if n.Init != nil {
            GoAST_Walk(v, n.Init)
        }
        if n.Cond != nil {
            GoAST_Walk(v, n.Cond)
        }
        if n.Post != nil {
            GoAST_Walk(v, n.Post)
        }
        GoAST_Walk(v, n.Body)

    case *RangeStmt:
        GoAST_Walk(v, n.Key)
        if n.Value != nil {
            GoAST_Walk(v, n.Value)
        }
        GoAST_Walk(v, n.X)
        GoAST_Walk(v, n.Body)

    // Declarations
    case *ImportSpec:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        if n.Name != nil {
            GoAST_Walk(v, n.Name)
        }
        GoAST_Walk(v, n.Path)
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *ValueSpec:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        walkIdentList(v, n.Names)
        if n.Type != nil {
            GoAST_Walk(v, n.Type)
        }
        walkExprList(v, n.Values)
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *TypeSpec:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        GoAST_Walk(v, n.Name)
        GoAST_Walk(v, n.Type)
        if n.Comment != nil {
            GoAST_Walk(v, n.Comment)
        }

    case *BadDecl:
        // nothing to do

    case *GenDecl:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        for _, s := range n.Specs {
            GoAST_Walk(v, s)
        }

    case *FuncDecl:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        if n.Recv != nil {
            GoAST_Walk(v, n.Recv)
        }
        GoAST_Walk(v, n.Name)
        GoAST_Walk(v, n.Type)
        if n.Body != nil {
            GoAST_Walk(v, n.Body)
        }

    // Files and packages
    case *File:
        if n.Doc != nil {
            GoAST_Walk(v, n.Doc)
        }
        GoAST_Walk(v, n.Name)
        walkDeclList(v, n.Decls)
        for _, g := range n.Comments {
            GoAST_Walk(v, g)
        }
        // don't walk n.Comments - they have been
        // visited already through the individual
        // nodes

    case *Package:
        for _, f := range n.Files {
            GoAST_Walk(v, f)
        }

    default:
        fmt.Printf("ast.Walk: unexpected node type %T", n)
        panic("ast.Walk")
    }

    v.Visit(nil)
}