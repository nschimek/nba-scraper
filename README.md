# nba-scraper
This command-line application will scrape NBA data from [Basketball Reference (BR)](https://www.basketball-reference.com) and load it into a relational MySQL Database.

## Initial Setup

### Create Database
This application requires a MySQL database to run.  Simply run the `sql/init.sql` file to create the required tables and keys.  Take note of your connection information as it will be needed next.

### Default Configuration File
In the `config` folder, there is a file called `sample.yaml`.  You should create a copy of this file and rename it `default.yaml`, then update it with your database connection information, among other settings.  They are explainted below.

## Configuration
The application can be configured via a YAML file or environment variables.  Some configuration settings can also be overriden at run-time by command-line parameters (described below).  

By default, the application expects a file called `default.yaml` in a `config` folder.  This is also generally the easiest way to configure it.

### YAML Configuration
By default, the application will look for a file called `config/default.yaml` for configuration.  There is a parameter to override this file if you wish.  The file must contain the following entries:

| Section | Attribute | Default | Description |
| :------ | :-------- | :-----: | :---------- |
| (root) | season | 2023 | The NBA season you are scraping in YYYY format.  Note that this should be the *finishing* year of the season (ex: the 2021-22 season would be entered as 2022) |
| (root) | debug | false | Enable or disable debug logging, which will output addtional log entries for troubleshooting |
| database | user | username | The username for connecting to your database |
| database | password | password | The password for connecting to your database |
| database | location | localhost | The location of your database |
| database | port | 3306 | The port used to connect to your database |
| database | name | nba | The name of your database |
| suppression | team | 30 | Teams that have been scraped more recently than this value (in days) will *not* be re-scraped during schedule, game, and standing scrapes.  (They can still be forcibly scraped by ID using the `teams` command.) |
| suppression | player | 365 | Players that have been scraped more recently than this value (in days) will *not* be re-scraped during schedule, game, and team scrapes.  (They can still be forcibly scraped by ID using the `players` command.) |

#### Sample
Here is a sample YAML configuration file:
```
season: 2023
debug: false
database:
  user: username
  password: mysupermegasecurepassword
  location: localhost
  port: 3306
  name: nba
suppression:
  team: 30
  player: 365
```

### Environment Variables
Coming soon.

## Running and Scraping
The scraper is run via the command line (Command Prompt in Windows).  Generally, you use the following format:
```
nba-scraper [-n|-d|-c] [schedule|games|teams|players|injuries|standings]
```

For example, `nba-scraper injuries` would scrape the injuires page.  **Note**: in Windows, `nba-scraper` will need to be `nba-scraper.exe`.

### Help
You can always get help with either the entire application or a specific command by using the `--help` parameter.

For example, `nba-scraper --help` will display help for the entire application, whereas `nba-scraper schedule --help` will display help for the schedule command.

### Global Parameters
These global parameters can be applied to any command and are strictly optional.  If they are not specified, the value from the configuration file will be used.
| Parameter | Flag | Description | Example |
| :-------- | :--: | :---------- | :------------ |
| season | `-n` | Override the season in the configuration file for this run | `nba-scraper -n 2022 <command>` |
| debug | `-d` | Enable debug mode | `nba-scraper -d <command>` |
| config | `-c` | Specify a different configuration file | `nba-scraper -c config/dev.yaml <command>` |

### Commands
Commands are generally what you will use to run the scraper.  Each command is dedicated to a different page.  Given the relational nature of the statistics, most commands will result in additional scrapes occuring to support the foreign key relationships.

| Command | Arguments/Flags | Description | Example |
| :------ | --------------- | :---------- | :------ | 
| schedule | <ul><li>`-s` Start Date in YYYY-MM-DD format</li><li>`-e` End Date in YYYY-MM-DD format</li><li>`-t` Also scrape the standings page for the current season</li><li>`-j` Also scrape the injuries page and load it for the current season</ul></ul> | Scrape games via the NBA schedule by a provided date range.  If no date range is provided, it defaults to yesterday. | <ul><li>`nba-scraper schedule -s 2021-11-01 -e 2021-11-30 -t -j` *(with injuries and standings)*</li><li>`nba-scraper schedule -s 2021-12-01 -e 2022-02-28` *(without injuries and standings)*</li></ul> | 
| games | (space-delimited BR game IDs) | Scrape individual games by BR IDs.  | `nba-scraper games 202102040LAL 202103150DEN`
| teams | (space-delimited BR team IDs) | Scrape individual teams by BR IDs. *Note: team suppression settings are ignored when using this command.* | `nba-scraper teams CHI GSW LAL` | 
| players | (space-delimited BR player IDs) | Scrape individual players by BR IDs. *Note: player suppression settings are ignored when using this command.* | `nba-scraper players lavinza01 jamesle01 curryst01` | 
| standings | *(none)* | Scrape the standings page for the currently configured season. | `nba-scraper standings` |
| injuries | *(none)* | Scrape the injuries page and load it for the current season. | `nba-scraper injuries` |
| version | *(none)* | Displays the current version.  Will also test DB connectivity and configured season. | `nba-scraper version` |

 > **Special Note on Standings**: standings for historical seasons are as of the last day of the season.

 > **Special Note on Injuries**: there are no historical injuires available, and only the current injuries can be scraped at any given time.  Take care when scraping historical seasons to avoid this page.

### Additional Scrapes
Due to the relational nature of the statistics, additional scrapes will occur for some commands in support of the foriegn keys:

 - `schedule`: games, teams, players, standings <sup>(if enabled)</sup>, injuries <sup>(if enabled)</sup> 
 - `games`: teams, players
 - `teams`: players
 - `standings`: teams, players
 - `injuries`: teams, players

### Running Without Commands
Running the scraper without commands will result in the `schedule` command being run with yesterday's date, along with the current standings and injuries.

## Special Note: Scraping a Historical Season
To scrape a historical season, you should do the following: 

 - Set the season parameter in your YAML configuration to the appropriate value.  Remember, use the year the season *ends* in.
 - Temporarily set `suppression.team` to 0 in the YAML configuration as well.
 - Run `nba-scraper standings` to scrape the standings as of the last day of the season.  Due to the suppression setting, this will scrape all teams in the standings (all of them), along with their rosters and salary data for that season.
 - Revert `suppression.team` back to 30, or whatever value you normally use
 - Run `nba-scraper schedule` with the `-s` and `-e` commands to scrape date ranges appropriate for the season you are scraping.  You can scrape across months and years, but *not* multiple seasons.  

 When done, be sure to revert the season value back to the current season if you are also in the process of scraping a season in-progress!