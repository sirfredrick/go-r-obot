# go-r-obot
>  Discord Bot that tracks how many Subreddits mentioned, rewritten in go

r-obot tracks how many times `r/<insert text here>` is mentioned in one of the messages on your server. It responds with various jokes related to the total message count.

Newly added in this go rewrite are unit tests in `main_test.go` as well as better joke generation.

## Running
First create a `.secret` text file beside the `main.go` and place your Discord API secret key in it.

Then make sure go is installed and on your path then run
```sh
go run main.go
```
To test run
```sh
go test main.go main_test.go
```
The `file.go` file is used to create a data file with a certain count.
To use it replace `336` in the `fmt.Sprint()` statement in `file.go` with the initial count you want and then run
```sh
go run file.go
```
You should then see the `data.txt` file

## Developing

Make sure go and git are installed and on your path then run
```sh
git clone https://github.com/sirfredrick/go-r-obot.git
```
then follow the above steps for running it.

## Contributing

If you'd like to contribute, please fork the repository and use a feature
branch. Pull requests are warmly welcome.

## Links

- Repository: https://github.com/sirfredrick/go-r-obot
- Issue tracker: https://github.com/sirfredrick/go-r-obot/issues
- Related projects:
   - Thanks to [discordgo](https://github.com/bwmarrin/discordgo) for making a Discord API bindings for go!

## License

Copyright © 2021 Sir Fredrick

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
