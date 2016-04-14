# StormCloud
This is a simplified version of Redis built in Go mostly for my amusement and to satisfy my need to reinvent the wheel.

#Current Features
- Key/List storage system
- Socket based communication
- Simple client tool to connect with system
- API to push, pop, list keys, empty keys, delete keys, as well as do some system maintenance. 
- Configurable persistant disk storage
- Configurable port setting


#Configuration
Example config.json
```{
  "port" : "6464",
  "usediskwriter" : true,
  "readfromdiskatstart" : true
}```

- **port** *number* - Port the StormCloud runs on
- **usediskwriter** *true/false* -  Active storage of StormCloud data to disk.
- **readfromdiskstart** *true/false* - Automatically load StormCloud data from disk on startup.

#API
 - **fpush** *key* *value* - Pushes to the front of the list for the given key
 - **bpush** *key* *value* - Pushes to the back of the list for the given key
 - **get** *key* - Gets all the values for the key
 - **keys** - Gets all the keys stored in the system
 - **fpop** *key* - Pops a value off the front of the list for the given key
 - **bpop** *key* - Pops a value off the back of the list for the given key
 - **empty** *key* - Empties the key of all its values
 - **deletekey** *key* - Deletes the key and its values
 - **savedata** - Saves current data in the StormCloud to disk
 - **loaddata** - Overwrites the current data in the StormCloud with data from the disk.
 - **autosave** *true/false*(optional) - Changes the usediskwriter setting from true or false.  If no parameter is provided it'll return the current state of usediskwriter.  **Note*** this change does not persist between restarts of the StormCloud.
