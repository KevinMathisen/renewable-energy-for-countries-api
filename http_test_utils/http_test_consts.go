package http_test_utils

import (
	"strconv"
)

// Values for fucntions
const FLOAT_PRECISION = 12

// Values for building URLs
const BEGIN_YEAR = "1990"
const END_YEAR = "2010"
const COUNTRY_CODE = "NOR"
const COUNTRY_NAME = "norway"
const BEGIN = "begin=" + BEGIN_YEAR
const END = "end=" + END_YEAR
const SORT_BY = "sortByValue=true"
const NEIGHBOURS = "neighbours=true"
const MEAN = "mean=true"
const PARAM = "?"
const AND = "&"

// Values for checking
const ALL_COUNTRIES = 79      //All different countries in the dataset
const CURRENT_COUNTRIES = 72  //Amount of countries with data for year 2021
const EXPECTED_NEIGHBOURS = 4 //Amount of neighbours for Norway

var NEIGHBOURS_CODES = []string{"FIN", "NOR", "RUS", "SWE"}        //The codes for Norway's neighbours in the default order
var SORTED_NEIGHBOURS_CODES = []string{"NOR", "SWE", "FIN", "RUS"} //The codes for Norway's neighbours in sorted order

var INT_BEGIN_YEAR, _ = strconv.Atoi(BEGIN_YEAR) //Int value of BEGIN_YEAR
var INT_END_YEAR, _ = strconv.Atoi(END_YEAR)     //Int value of END_YEAR

const COUNTRY_OLDEST_PERCENTAGE = 67.87996  //Oldest percentage for Norway
const COUNTRY_LATEST_PERCENTAGE = 71.558365 //Latest percentage for Norway
const COUNTRY_EXPECTED_ENTRIES = 57         //Amount of entries Norway has in the dataset

const COUNTRY_BEGIN_PERCENTAGE = 72.44774 //Percentage for Norway in year BEGIN_YEAR
const COUNTRY_END_PERCENTAGE = 65.47019   //Percentage for Norway in year END_YEAR
const COUNTRY_BEGIN_END_ENTRIES = 21      //Amount of entries Norway has in the dataset between BEGIN_YEAR and END_YEAR
const COUNTRY_BEGIN_ENTRIES = 32          //Amount of entries Norway has in the dataset between BEGIN_YEAR and the end of the dataset
const COUNTRY_END_ENTRIES = 46            //Amount of entries Norway has in the dataset between the start of the dataset and END_YEAR

const COUNTRY_BEGIN_END_SORT_FIRST = 1990                //The year of the first object after sort
const COUNTRY_BEGIN_END_SORT_FIRST_PERCENTAGE = 72.44774 //The percentage of the first object after sort
const COUNTRY_BEGIN_END_SORT_LAST = 2003                 //The year of the last object after sort
const COUNTRY_BEGIN_END_SORT_LAST_PERCENTAGE = 63.816036 //The percentage of the last object after sort

const COUNTRY_MEAN = 68.01918892982457           //Mean percentage for Norway
const COUNTRY_BEGIN_END_MEAN = 68.63185428571428 //Mean percentage for Norway between BEGIN_YEAR and END_YEAR

const NEIGHBOUR_ENTRIES_AMOUNT = 208  //Amount of objects returned when calling for the neighbours of Norway
const NEIGHBOUR_BEGIN_END_AMOUNT = 84 //Amount of objects returned when calling for the nieghbours of Norway between BEGIN_YEAR and END_YEAR

const NEIGHBOURS_SORT_LAST_CODE = "RUS"          //The ISO code of the last country in the list when sorted by oercentage
const NEIGHBOURS_SORT_LAST_NAME = "russia"       //The name of the last country in the list when sorted by oercentage
const NEIGHBOURS_SORT_LAST_PERCENTAGE = 4.605263 //The percentage of the last country in the list when sorted by oercentage
const NEIGHBOURS_SORT_LAST_YEAR = 1989           //The year of the last country in the list when sorted by oercentage

const NEIGHBOURS_SORT_MEAN_LAST = 6.004957597297297 //Percentage of the last country recieved from the sorted mean of the neighbours

const ALL_SORT_LAST_CODE = "SAU"                       //The ISO code of the last country in the list when sorting all countries by percentage
const ALL_SORT_LAST_NAME = "Saudi Arabia"              //The name of the last country in the list when sorting all countries by percentage
const ALL_SORT_LAST_PERCENTAGE = 0.0013665377020357142 //The percentage of the last country in the list when sorting all countries by percentage
