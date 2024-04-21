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

	styleOfValue := reflect.ValueOf(styleOfStruct)
	postTypeValue := reflect.ValueOf(postTypeStruct)
	technologiesValue := reflect.ValueOf(technologiesStruct)

	var styleOfSlice []interface{}
	var postTypeSlice []interface{}
	var technologiesSlice []interface{}

	for i := 0; i < styleOfValue.NumField(); i++ {
		fieldValue := styleOfValue.Field(i).Interface()
		styleOfSlice = append(styleOfSlice, fieldValue)
	}
	source := rand.NewSource(time.Now().UnixNano())

	rng := rand.New(source)

	styleRng := rng.Intn(len(styleOfSlice) - 1)

	customQuery :=
		projectConfig.ChatGPTQueries.PostType.Guide +
			projectConfig.ChatGPTQueries.StyleOf.Manager +
			projectConfig.ChatGPTQueries.GeneralTopic[styleRng]

	return customQuery
}
