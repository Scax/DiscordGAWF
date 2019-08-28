# Discord GAWF

Discord Game Activity for Wine and Firejail

This is a little tool inspired by [DiscordComp](https://github.com/null-von-sushi/Discordcomp).

If you are running some games or programs via wine discord wont let you choose these processes for game activity

This program looks for .exe process and creates a .dummy file and process you can select in discord

If you run discord in Firejail you this programm can start the dummies inside the discord firejail

## How to use

Download the files from releases or build them yourself

### If you use just wine

Copy the files `DiscordGAWF` and `dummy` to any directory and run the `DiscordGAWF`.
This will create a .dummy file for every .exe process running

### If you use firejail
Make sure you use `-use-firejail` and that the current directory is whitelisted in the discord profile or
specify a whitelisted directory using `-firejail-path-override` and make sure the firejail is discord is running in is named "discord" or set the name via `-firejail-name`

### Help
```
Usage of ./DicordGAWF:
  -dummies string
        The location of where the dummies will be stored (default "./.")
  -firejail-name string
        the name of the firejail discord is running in. (default "discord")
  -firejail-path-override string
        override the dummy storage passed to firejail 'firejail --join=discord override/something.dummy'
  -use-firejail
        Set if this discord is running in a firejail
  -wait int
        duration in seconds between scans for .exe processes (default 60)
```
## Note
If your system gets compromised this tool could be be used to further exploit you.

## TODO

Add the possibility to specify a exclude or include list
Add option for auto clean
