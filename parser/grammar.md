```
expression     → equality ;
assignment     → IDENTIFIER "=" expr 
               | equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")"
               | IDENTIFIER ;


program        → declaration* EOF ;
declaration    → varDecl 
               | statement;
statement      → exprStmt
               | printStmt
               | blockStmt 
               | ifStmt ;
               | loopStmt ;
               | forLoopStmt ;

ifStmt         → "if" expression blockStmt 
                ("elif" expression blockStmt)* 
                ("else" blockStmt)? ;

loopStmt       → "loop" block
forLoopStmt    → "for" (varDecl | exprStmt)? "," expression? "," expression? block ;

exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
varDecl        → "var" IDENTIFIER "=" expr ;
blockStmt      → "{" declaration* "}" ;
```
