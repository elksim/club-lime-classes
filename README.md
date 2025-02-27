#

**_Not affiliated with Club Lime in any way._**

Club Lime provides a way to explore all of the classes that they have scheduled [(Club Lime classes)](https://www.clublime.com.au/classes/)*, but in my opinion it's quite clunky to explore if you're interested in multiple gym locations or session types. Fortunately, you can download all upcoming scheduled classes as a CSV file.

This website uses this CSV file and aims to makes it easier to explore the scheduled classes you're interested in.

## todo

- [ ] options above the table where we can select which locations and which workouts we are interested in
- [ ] filter the table by the chosen options
- [ ] want to generate the index page once per day when we fetch the new csv and then just serve this generated page
- [x] settings page where we can cut down the list of locations and workouts that get shown on the main page as options
- [ ] source the raw data CSV from club lime, register a fly volume and write and read the raw data from there
- [ ] might be cool to store the raw data CSV once per day
- [ ] unfortunately the csv doesn't include data on which classes are full and which are cancelled, but I think we might be able to get this information actually
- [ ] pottenntially possible to make the classes clickable and take you to club lime to view and book, but I'm not sure about that.
- the club lime phone app is actually quite nice for class exploration - makes sense for me to look to differenciate from it
  - [ ] could have buttons like 'yoga' 'spin' 'cardio'? 'weights'? which shows like, all the yoga classes for the week
