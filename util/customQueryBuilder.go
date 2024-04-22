package projectUtil

import (
	"math/rand"
	"reflect"
	"time"
)

var projectConfig TProjectConfig = InitProjectConfig()

func CustomQueryBuilder() string {

	styleOfStruct := projectConfig.ChatGPTQueries.StyleOf
	postTypeStruct := projectConfig.ChatGPTQueries.PostType
	technologiesStruct := projectConfig.ChatGPTQueries.Subject

	styleRng, styleSlice := rngForStruct(styleOfStruct)
	postTypeRng, postSlice := rngForStruct(postTypeStruct)
	techRng, techSlice := rngForStruct(technologiesStruct)

	var customQuery string = postSlice[postTypeRng] +
		styleSlice[styleRng] +
		techSlice[techRng]

	return customQuery
}

func rngForStruct(s interface{}) (int, []string) {

	// create a fresh source to create rng from.
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// create a variable of type reflect.Value
	sValue := reflect.ValueOf(s)

	// create a string slice of name structSlice to append struct values to.
	var structSlice []string

	// loop through the values of the reflect.Value and appends them to slice as strings.
	// if the type of value is an Array, loop through arrary and append to slice
	for i := 0; i < sValue.NumField(); i++ {
		fieldValue := sValue.Field(i)
		if fieldValue.Kind() == reflect.Array {
			for j := 0; j < fieldValue.Len(); j++ {
				structSlice = append(structSlice, fieldValue.Index(j).String())
			}
		} else {
			structSlice = append(structSlice, fieldValue.String())
		}

	}

	// create an int based on size of slice, minus 1. ex: slice of size 4, random number will be between 0 and 3
	structRng := rng.Intn(len(structSlice) - 1)

	// return the rng int as well as the slice with values appended
	return structRng, structSlice
}
