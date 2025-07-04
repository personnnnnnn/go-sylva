grammar Sylva;
INT: [0-9]+;
FLOAT: INT '.' INT | INT '.' | '.' INT;
STRING:
	'"' (~(["\\\r\n]) | '\\' ~[\r\n])* '"'
	| '\'' (~(['\\\r\n]) | '\\' ~[\r\n])* '\'';

expr: value;
value: INT | FLOAT | STRING | '(' expr ')';