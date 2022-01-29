# duranzcli
 
 Package written in GO programming Language for storing/extracting stats to perform analysis from the ball by ball data available on cricsheet.
 
  - Dataset: Cricsheet [https://cricsheet.org/]
> Currently works with odi, t20, test, ipl data


Apply your cricketing brain and run queries on the extracted stats to find amazing insights. 

Common Stats Examples.

    - Batting Stats
        Runs Scored
        Balls Played
        Average
        Strike Rate
        Number of 6s, 4s, 3s, 2s, 1
        Out Types
    - Bowling Stats
        Runs Conceded
        Balls Bowled
        Average
        Economy
        6s, 4s, 3s, 2s, 1 Conceded
        Maiden
    - Fielding Stats
        Catches
        Run Outs
        Stumpings
    - Player Vs Player Stats
    - Batsman X vs Bowler Y
    - No. of times dismissal (by type, e.g, bowled 10 times)
    

> Check this out also : duranzapi (web application). 
> To be published soon and will be available on https://www.kunalduran.com

# process to run the application

##### # pre-requisites
- Go 
- mysql

> step1 : clone repository 
> step2 : download cricsheet data anywhere
> step3 : provide path to data in the module/process/cricket.go line:15
> step4 : run go build or go run main.go with appropriate arguments

### Available commands example:

- To create all the tables run *go run main.go initial*
- To empty the complete database *go run main.go delete*
- Run format wise commands *go run main.go all odi* 
- *go run main.go all ipl*
- Run process for a particular matchfile *go run main.go all test 42242.json*

    