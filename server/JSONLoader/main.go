package JSONLoader

import (
	"encoding/json"
	"log"
	"os"
)

// World are worlds Obi Wan Kenobi plans to visit
type World struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//DarkJedi are an array struct of dark jedi
type DarkJedi struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Homeworld  World  `json:"Homeworld"`
	Master     int    `json:"master"`
	Apprentice int    `json:"apprentice"`
}

//AugmentURL contains relational URLs for DarkJedi
type AugmentURL struct {
	URL string `json:"url"`
	ID  int    `json:"id"`
}

//AugmentedDarkJedi is a DarkJedi with relational URLs
type AugmentedDarkJedi struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Homeworld  World      `json:"Homeworld"`
	Master     AugmentURL `json:"master"`
	Apprentice AugmentURL `json:"apprentice"`
}

//LoadJSON marshals decodes dark jedis and worlds JSONs into structs
func LoadJSON() ([]World, []DarkJedi) {
	worldsFile, err := os.Open("worlds.json")
	if err != nil {
		log.Fatal(err)
	}
	defer worldsFile.Close()

	darkJediFile, err := os.Open("dark-jedis.json")
	if err != nil {
		log.Fatal(err)
	}
	defer darkJediFile.Close()

	var worlds []World
	if err := json.NewDecoder(worldsFile).Decode(&worlds); err != nil {
		log.Fatal(err)
	}

	var darkJedis []DarkJedi
	if err := json.NewDecoder(darkJediFile).Decode(&darkJedis); err != nil {
		log.Fatal(err)
	}
	return worlds, darkJedis
}
