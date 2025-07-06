grammar Sylva;
INT: [0-9]+;
FLOAT: INT '.' INT | INT '.' | '.' INT;
STRING:
	'"' (~(["\\\r\n]) | '\\' ~[\r\n])* '"'
	| '\'' (~(['\\\r\n]) | '\\' ~[\r\n])* '\'';
TRUE: 'true';
FALSE: 'false';
BOOL: TRUE | FALSE;
ID: [_]* [a-zA-Z][a-zA-Z0-9]* | '_';
WHITESPACE: [ \t\r\n] -> skip;

expr:
	(op = '+' | '-') expr			# UnaryOp
	| expr ('..') expr				# ConcatExpr
	| expr (op = '*' | '/') expr	# MulExpr
	| expr (op = '+' | '-') expr	# AddExpr
	| value							# ValueExpr;

value:
	INT									# IntValue
	| FLOAT								# FloatValue
	| STRING							# StringValue
	| BOOL								# BoolValue
	| '[' (expr (',' expr)* ','?)? ']'	# ListValue
	| '(' expr ')'						# ParensValue;