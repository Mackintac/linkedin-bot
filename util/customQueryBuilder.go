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
	technologiesStruct := projectConfig.ChatGPTQueries.Technologies

	styleRng, styleSlice := rngForStruct(styleOfStruct)
	postTypeRng, postSlice := rngForStruct(postTypeStruct)
	techRng, techSlice := rngForStruct(technologiesStruct)

	var customQuery string = postSlice[postTypeRng] +
		styleSlice[styleRng] +
		techSlice[techRng]

	return customQuery
}

func rngForStruct(s interface{}) (int, []string) {

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	sValue := reflect.ValueOf(s)
	var structSlice []string

	// str, ok := sValue.(string)
	// if !ok {

	// }
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

	structRng := rng.Intn(len(structSlice) - 1)

	return structRng, structSlice
}
