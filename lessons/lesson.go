package lessons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func Init(filePath string) (Lesson, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("ERROR!!")
		fmt.Println(err)
		return Lesson{}, err
	}

	fmt.Println("Successfully Opened " + filePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var lesson Lesson

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &lesson)

	return lesson, nil
}

type Lesson struct {
	Descriptions       []DescriptionCard   `json:"descriptions"`
	InteractiveActions []InteractiveAction `json:"interactiveActions"`
	Quiz               []Question          `json:"quiz"`
}

type DescriptionCard struct {
	MainHeader string `json:"mainHeader"`
	Header     string `json:"header"`
	Data       string `json:"data"`
}

type InteractiveAction struct {
}

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}
