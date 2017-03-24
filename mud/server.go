package mud

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var gameServer Server

type Server struct {
	connections []connection
}

func (server *Server) handle(c *connection) {
	c.mob = login(c)

	err := registerConnection(c)
	if err != nil {
		return
	}

	newAction(c.mob, c, "look")
	for {
		c.mob.statusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			err := newActionWithInput(&action{mob: c.mob, conn: c, args: strings.Split(input, " ")})
			if err != nil {
				return // we're quitting
			}
		}

	}
}

func (server *Server) Serve(port int) {
	server.init()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("cannot start server: %s\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("waiting for connections on %s\n", listener.Addr())

	go server.timing()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("could not accept: %s\n", err)
		} else {
			fmt.Printf("connected: %s\n", conn.RemoteAddr())
			go server.handle(newConnection(conn))
		}
	}
}

func (server *Server) timing() {
	const (
		tickLen time.Duration = 5
	)

	pulse := time.NewTicker(time.Second)
	tick := time.NewTicker(time.Second * tickLen)

	for {
		select {
		case <-pulse.C:
			fmt.Printf(".")
			for e := mobList.Front(); e != nil; e = e.Next() {
				m := e.Value.(*mob)
				f := m.Fight
				if f != nil && m.Status == fighting {
					f.turn(m)
				}

				if m.wait < 0 {
					m.wait--
				}

				for _, af := range m.Affects {
					if af.duration > 0 {
						af.duration--
					} else {
						m.removeAffect(af)
					}
				}
			}
			break
		case <-tick.C:
			fmt.Printf("o")
			for e := mobList.Front(); e != nil; e = e.Next() {
				m := e.Value.(*mob)
				if m.Playable == false {
					m.wander()
				}
				m.notify(helpers.Newline)
				m.statusBar()
				m.regen()
			}

			for e := roomList.Front(); e != nil; e = e.Next() {
				r := e.Value.(*room)
				r.decayItems()
			}
			break
		}
	}
}

func registerConnection(c *connection) error {
	for _, oc := range gameServer.connections {
		if c.mob.ID == oc.mob.ID {
			c.SendString(fmt.Sprintf("This user is already playing. Bye! %s", helpers.Newline))
			c.end()
			return errors.New("this user is already playing")
		}
	}

	gameServer.connections = append(gameServer.connections, *c)
	return nil
}

var (
	db *gorm.DB
)

func (server *Server) init() {
	initializeDatabase()
}

func getConnections() []connection {
	return gameServer.connections
}

func initializeDatabase() {
	var err error
	db, _ = gorm.Open("mysql", "homestead:secret@tcp(api.guidebox.brian:3306)/mud?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&mob{}, &job{}, &race{}, &item{}, &area{}, &room{}, &exit{}, &skill{}, &mobSkill{}, &skillLevel{})

	newSkillDatabase()
	newJobDatabase()
	newRaceDatabase()
	newItemDatabase()
	newMobDatabase()
	newRoomDatabase()

	fmt.Println("itemGlow ", itemGlow)
	fmt.Println("itemHum ", itemHum)
	fmt.Println("itemDark ", itemDark)
	fmt.Println("itemLock ", itemLock)
	fmt.Println("itemEvil ", itemEvil)
	fmt.Println("itemInvis ", itemInvis)
	fmt.Println("itemMagic ", itemMagic)

	os.Exit(1)
}
