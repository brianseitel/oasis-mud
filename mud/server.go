package mud

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	// "github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

type Server struct {
	connections []connection
}

func (server *Server) Handle(c *connection) {
	c.player = login(c)
	newAction(c.player, c, "look")
	for {
		c.player.ShowStatusBar()
		input, err := c.buffer.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.Trim(input, "\r\n")

		if len(input) > 0 {
			err := newActionWithInput(&action{player: c.player, conn: c, args: strings.Split(input, " ")})
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
			go server.Handle(newConnection(conn))
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
			fmt.Printf("")
			break
		case <-tick.C:
			var rooms []room
			db.Find(&rooms)
			for _, r := range rooms {
				for _, m := range r.Mobs {
					m.wander()
				}
			}
			break
		}
	}
}

var (
	db *gorm.DB
)

func (server *Server) init() {
	initializeDatabase()
}

func initializeDatabase() {
	var err error
	db, _ = gorm.Open("mysql", "homestead:secret@tcp(api.guidebox.brian:3306)/mud?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	// db.LogMode(true)
	// defer db.Close()
	db.AutoMigrate(&mob{}, &job{}, &race{}, &item{}, &area{}, &room{}, &exit{}, &player{})

	newJobDatabase()
	newRaceDatabase()
	newItemDatabase()
	newMobDatabase()
	newRoomDatabase()
	newPlayerDatabase()
}
