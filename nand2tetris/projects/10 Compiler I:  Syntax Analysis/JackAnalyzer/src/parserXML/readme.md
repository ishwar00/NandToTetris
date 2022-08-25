## It may not produce xml output as expected on grouped expressions

_syntax tree_ produced in _xml_ file is done by walking through AST produced by `parser`. 
Problem is with expressions like `(4 + 2) + 2`, here there is no need of parentheses, it's AST and `4 + 2 + 2` 
are same. AST does not store information about grouped expressions explicitly.
So it not possible to detect such parentheses. Where as expressions like `4 + (2 + 2)`, 
it is possible to detect that parentheses were used to alter the precedence and order of evaluation. 

AST produced for `(4 + 2) + 2` and `4 + 2 + 2` are same, and looks like
```
          +
        /   \
       +     2
     /   \
    4     2
```
where as AST produced for `4 + (2 + 2)` is
```
             + 
           /   \
          4     +
              /   \
             2     2 
```
so AST produced in _xml_ may not contain these undetected parentheses.

only two were comapared successfully\
✔️  ArrayTest\
✔️  ExpressionLessSquare\
❌ Square\

As it is not possible detect parentheses which do no alter precedence.
__Parser works, don't get the wrong idea__
