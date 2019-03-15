package lessons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func Init(filePath string) {
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("ERROR!!")
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened " + filePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	fmt.Println(string(byteValue))
	// we initialize our Users array
	var lesson lesson

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &lesson)
	fmt.Println(lesson)

	showDescriptions(lesson.Descriptions)
	showInteractions(lesson.InteractiveActions)
	showQuizes(lesson.Quiz)
}
func showQuizes(questions []question) {

}
func showInteractions(actions []interactiveAction) {

}
func showDescriptions(cards []descriptionCard) {

}

type lesson struct {
	Descriptions       []descriptionCard   `json:"descriptions"`
	InteractiveActions []interactiveAction `json:"interactiveActions"`
	Quiz               []question          `json:"quiz"`
}

type descriptionCard struct {
	MainHeader string `json:"mainHeader"`
	Header     string `json:"header"`
	Data       string `json:"data"`
}

type interactiveAction struct {
}

type question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}
