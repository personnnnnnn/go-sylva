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

   values can implement the "Object" interface which
   gives them the ability to have static properties (like
   methods, by using the ":" like in lua), but not
   be directly stored in the object to conserve memory

   userdata -> anything not on this list

   TODO: enums (rust style)
*/

type Object interface {
	GetStatic(key string) (Value, error)
}
