package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
	//"github.com/Clinet/discordgo-embed"
)

type Channel struct {
	Id     string `json:"id"`
	Prefix string `json:"prefix"`
}

//Checks wether message author has administrator permissions
func IsAdmin(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	perm, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return false, err
	}
	return perm&discordgo.PermissionAdministrator != 0, nil
}

//Creating buckets
func CreateBuckets() {
	database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Chans"))
		if err != nil {
			return fmt.Errorf("creation Error: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Guilds"))
		if err != nil {
			return fmt.Errorf("creation Error: %s", err)
		}
		return nil
	})
}

//Loading channels from database
func LoadChannels() {
	channels = make(map[string]Channel)
	database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Chans"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ch Channel
			err := json.Unmarshal(v, &ch)
			if err != nil {
				consoleLog.Printf("[SETUP] Error getting channel %s", k)
				continue
			}
			channels[string(k)] = ch
		}
		return nil
	})
}

//Loading guilds from database
func LoadGuilds() {
	guilds = make(map[string]bool)
	database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Guilds"))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			guilds[string(k)] = true
		}
		return nil
	})
}

//Adding channel to database
func CreateChan(id, prefix string) (Channel, error) {
	ch := Channel{id, prefix}
	v, _ := json.Marshal(ch)
	err := database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Chans"))
		err := b.Put([]byte(id), v)
		return err
	})
	return ch, err
}

//Adding guild to database
func CreateGuild(id string) error {
	v, _ := json.Marshal(true)
	err := database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Guilds"))
		err := b.Put([]byte(id), v)
		return err
	})
	return err
}

//Loading library from json to memory
func LoadLibrary() {
	library = make(map[string]discordgo.MessageEmbed)
	jsonFile, err := os.Open("data.json")
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &library)
}