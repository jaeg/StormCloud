# StormCloud
This is a simplified version of Redis built in Go mostly for my amusement and to satisfy my need to reinvent the wheel.

#Current Features
- An interface to interact with the system.
 - fpush *key* *value* - Pushes to the front of the list for the given key
 - bpush *key* *value* - Pushes to the back of the list for the given key
 - get *key* - Gets all the values for the key
 - keys - Gets all the keys stored in the system
 - fpop *key* - Pops a value off the front of the list for the given key
 - bpop *key* - Pops a value off the back of the list for the given key
 - empty *key* - Empties the key of all its values
 - deletekey *key* - Deletes the key and its values
