package sylva

type Value any

/*
   a `Value` can be a:

   int -> int
   float -> float64
   boolean -> bool
   string -> string

   nil -> nil (but should never be actually used)

   function -> func(args Value...) Value, error

   table ({ ...items }) -> map[Value]Value
   list ([...items]) -> []Value

   struct (${...items}) -> *Struct
   tuple ( (x = 0, y = 0, ...) ) -> Tuple

   userdata -> anything not on this list

   TODO: enums (rust style)
*/
