package mud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

var bootTime time.Time

func startAPI() {
	bootTime = time.Now()
	http.HandleFunc("/", addDefaultHeaders(welcome))

	http.HandleFunc("/items", addDefaultHeaders(listItems))
	http.HandleFunc("/mobs", addDefaultHeaders(listMobs))
	http.HandleFunc("/players", addDefaultHeaders(listPlayers))
	http.HandleFunc("/players/stats", addDefaultHeaders(getPlayersStats))
	http.HandleFunc("/rooms", addDefaultHeaders(listRooms))
	http.HandleFunc("/uptime", addDefaultHeaders(uptime))
	http.ListenAndServe(":8080", nil)
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}

func uptime(w http.ResponseWriter, r *http.Request) {
	type result struct {
		BootedAt time.Time
		Uptime   int
	}

	uptime := time.Now().Sub(bootTime).Seconds()
	results := &result{BootedAt: bootTime, Uptime: int(uptime)}
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the MUD API!")
}

type APIArea struct {
	ID         int
	Name       string
	Age        int
	NumPlayers int
}

type APIExit struct {
	ID          int
	Keyword     string
	Description string
	Dir         string
	RoomID      int
	Key         *itemIndex
	Flags       []string
}

type APIItem struct {
	ID               int
	Name             string
	Description      string
	ShortDescription string
	ItemType         string
	ExtraFlags       string
	WearFlags        string
	WearLocation     string
	Weight           int
	Cost             int
	Level            int
	Timer            int
	Value            int
	Min              int
	Max              int
	Charges          int

	ClosedFlags int
}

type APIMob struct {
	ID     int
	Name   string
	Level  int
	Job    *job
	Race   *race
	Room   *APIRoom
	RoomID int
}

type APIRoom struct {
	ID   int
	Name string

	Area        APIArea
	AreaID      int
	Description string
	Exits       []*APIExit
	Items       []*item
	Mobs        []*APIMob

	Light      int
	RoomFlags  int
	SectorType int
}

type APIPlayer struct {
	ID      int
	SavedAt string

	//Mob information
	Name            string
	Description     string
	LongDescription string
	Title           string

	Affects    []*affect /* list of affects, incl durations */
	AffectedBy []string
	Act        []string

	Skills    []*mobSkill
	Inventory []*item
	Equipped  []*item
	Room      *APIRoom

	ExitVerb string
	Bamfout  string
	Bamfin   string

	Hitpoints    int
	MaxHitpoints int
	Mana         int
	MaxMana      int
	Movement     int
	MaxMovement  int

	Armor   int
	Hitroll int
	Damroll int

	Exp       int
	Level     int
	Alignment int
	Practices int
	Gold      int
	Trust     int

	Carrying       int
	CarryMax       int
	CarryWeight    int
	CarryWeightMax int

	Job    *job
	Race   *race
	Gender string

	Attributes         *attributeSet
	ModifiedAttributes *attributeSet

	Status       status
	RecallRoomID int
}

func getPlayersStats(w http.ResponseWriter, r *http.Request) {
	playerFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/players/*.json"))

	var dates []string

	then := time.Now()
	for i := 0; i < 30; i++ {
		then = then.AddDate(0, 0, -1)
		dates = append(dates, then.Format("2006-01-02"))
	}

	var players []*mobIndex
	for _, playerFile := range playerFiles {
		file, err := ioutil.ReadFile(playerFile)
		if err != nil {
			panic(err)
		}

		var player *mobIndex
		err = json.Unmarshal(file, &player)
		if err != nil {
			panic(err)
		}

		players = append(players, player)
	}

	createdAtResults := make(map[string]int)
	for _, p := range players {
		date, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", p.CreatedAt)
		if _, ok := createdAtResults[p.CreatedAt]; !ok {
			createdAtResults[date.Format("2006-01-02")] = 1
		} else {
			createdAtResults[date.Format("2006-01-02")]++
		}
	}

	savedAtResults := make(map[string]int)
	for _, p := range players {
		date, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", p.SavedAt)
		if _, ok := savedAtResults[p.SavedAt]; !ok {
			savedAtResults[date.Format("2006-01-02")] = 1
		} else {
			savedAtResults[date.Format("2006-01-02")]++
		}
	}

	lastSeenAtResults := make(map[string]int)
	for _, p := range players {
		date, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", p.LastSeenAt)
		if _, ok := lastSeenAtResults[p.LastSeenAt]; !ok {
			lastSeenAtResults[date.Format("2006-01-02")] = 1
		} else {
			lastSeenAtResults[date.Format("2006-01-02")]++
		}
	}

	createdAtFinal := make(map[string]int)
	for d, count := range createdAtResults {
		for _, date := range dates {
			if d == date {
				createdAtFinal[date] = count
			} else {
				createdAtFinal[date] = 0
			}
		}
	}

	savedAtFinal := make(map[string]int)
	for d, count := range savedAtResults {
		for _, date := range dates {
			if d == date {
				savedAtFinal[date] = count
			} else {
				savedAtFinal[date] = 0
			}
		}
	}

	lastSeenAtFinal := make(map[string]int)
	for d, count := range lastSeenAtResults {
		for _, date := range dates {
			if d == date {
				lastSeenAtFinal[date] = count
			} else {
				lastSeenAtFinal[date] = 0
			}
		}
	}

	type Results struct {
		Results struct {
			CreatedAt  map[string]int
			SavedAt    map[string]int
			LastSeenAt map[string]int
		}
	}

	results := &Results{}
	results.Results.CreatedAt = createdAtFinal
	results.Results.SavedAt = savedAtFinal
	results.Results.LastSeenAt = lastSeenAtFinal

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

func listItems(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")

	var items []*APIItem
	if count == "" {
		for e := itemList.Front(); e != nil; e = e.Next() {
			i := e.Value.(*item)

			oi := &APIItem{}
			oi.ID = i.ID
			oi.Name = i.Name
			oi.Description = i.Description
			oi.ShortDescription = i.ShortDescription
			oi.ItemType = itemTypeName(i)
			oi.ExtraFlags = extraBitName(i.ExtraFlags)
			oi.Weight = i.Weight
			oi.Cost = i.Cost
			oi.Level = i.Level
			oi.Timer = i.Timer
			oi.Min = i.Min
			oi.Max = i.Max
			oi.Value = i.Value
			oi.Charges = i.Charges

			items = append(items, oi)
		}
	}

	type Result struct {
		TotalActive int
		TotalIndex  int
		Results     []*APIItem
	}

	results := &Result{}
	results.TotalIndex = itemIndexList.Len()
	results.TotalActive = itemList.Len()
	results.Results = items

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

func listMobs(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")

	var mobs []*APIMob
	if count == "" {
		for e := mobList.Front(); e != nil; e = e.Next() {
			m := e.Value.(*mob)
			if m.client == nil || m.isNPC() {
				om := &APIMob{}
				om.ID = m.ID
				om.Name = m.Name
				om.Race = m.Race
				om.Job = m.Job
				om.Level = m.Level

				om.Room = &APIRoom{}
				if m.Room != nil {
					om.Room.ID = m.Room.ID
					om.Room.Area = APIArea{ID: m.Room.Area.ID, Name: m.Room.Area.Name, Age: m.Room.Area.age, NumPlayers: m.Room.Area.numPlayers}
					om.Room.Name = m.Room.Name
					om.Room.Description = m.Room.Description
					for _, x := range m.Room.Exits {
						ex := &APIExit{}
						ex.ID = x.ID
						ex.Keyword = x.Keyword
						ex.Description = x.Description
						ex.Dir = x.Dir
						ex.RoomID = x.RoomID
						ex.Key = getItem(x.Key)
						ex.Flags = getExitFlags(x.Flags)
						om.Room.Exits = append(om.Room.Exits, ex)
					}

					for _, mx := range m.Room.Mobs {
						m1 := &APIMob{}
						m1.ID = mx.ID
						m1.Name = mx.Name
						m1.Level = mx.Level
						m1.Job = mx.Job
						m1.Race = mx.Race
						m1.RoomID = mx.Room.ID
						om.Room.Mobs = append(om.Room.Mobs, m1)
					}
				}
				mobs = append(mobs, om)
			}
		}
	}

	type Result struct {
		TotalActive int
		TotalIndex  int
		Results     []*APIMob
	}
	results := &Result{}
	results.TotalIndex = mobIndexList.Len()
	results.TotalActive = mobList.Len()
	results.Results = mobs

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

func listPlayers(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")

	var players []*APIPlayer

	if count != "" {
		for e := mobList.Front(); e != nil; e = e.Next() {
			p := e.Value.(*mob)
			if p.client != nil && !p.isNPC() {
				players = append(players, &APIPlayer{})
			}
		}
	} else {
		for e := mobList.Front(); e != nil; e = e.Next() {
			p := e.Value.(*mob)
			if p.client != nil && !p.isNPC() {
				op := &APIPlayer{}
				op.ID = p.ID
				op.Name = p.Name
				op.SavedAt = p.SavedAt
				op.Description = p.Description
				op.LongDescription = p.LongDescription
				op.Title = p.Title

				op.Affects = p.Affects
				op.AffectedBy = getAffectNames(p.AffectedBy)
				op.Act = getActFlags(p.Act)

				op.Inventory = p.Inventory
				op.Equipped = p.Equipped

				op.Room = &APIRoom{}
				if p.Room != nil {
					op.Room.ID = p.Room.ID
					op.Room.Area = APIArea{ID: p.Room.Area.ID, Name: p.Room.Area.Name, Age: p.Room.Area.age, NumPlayers: p.Room.Area.numPlayers}
					op.Room.Name = p.Room.Name
					op.Room.Description = p.Room.Description
					for _, x := range p.Room.Exits {
						ex := &APIExit{}
						ex.ID = x.ID
						ex.Keyword = x.Keyword
						ex.Description = x.Description
						ex.Dir = x.Dir
						ex.RoomID = x.RoomID
						ex.Key = getItem(x.Key)
						ex.Flags = getExitFlags(x.Flags)
						op.Room.Exits = append(op.Room.Exits, ex)
					}

					for _, m := range p.Room.Mobs {
						m1 := &APIMob{}
						m1.ID = m.ID
						m1.Name = m.Name
						m1.Level = m.Level
						m1.Job = m.Job
						m1.Race = m.Race
						m1.RoomID = m.Room.ID
						op.Room.Mobs = append(op.Room.Mobs, m1)
					}
				}

				op.Skills = p.Skills
				op.ExitVerb = p.ExitVerb
				op.Bamfout = p.Bamfout
				op.Bamfin = p.Bamfin
				op.Hitpoints = p.Hitpoints
				op.MaxHitpoints = p.MaxHitpoints
				op.Mana = p.Mana
				op.MaxMana = p.MaxMana
				op.Movement = p.Movement
				op.MaxMovement = p.MaxMovement
				op.Armor = p.Armor
				op.Hitroll = p.Hitroll
				op.Damroll = p.Damroll
				op.Exp = p.Exp
				op.Level = p.Level
				op.Alignment = p.Alignment
				op.Practices = p.Practices
				op.Gold = p.Gold
				op.Trust = p.Trust
				op.Job = p.Job
				op.Race = p.Race
				switch p.Gender {
				case 0:
					op.Gender = "Male"
					break
				case 1:
					op.Gender = "Female"
					break
				default:
					op.Gender = "Neutral"
					break
				}

				op.Attributes = p.Attributes
				op.ModifiedAttributes = p.ModifiedAttributes

				players = append(players, op)
			}
		}
	}

	type Result struct {
		Total   int
		Results []*APIPlayer
	}
	results := &Result{Total: len(players), Results: players}

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}

func listRooms(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")

	var rooms []*APIRoom
	if count == "" {
		for e := roomList.Front(); e != nil; e = e.Next() {
			r := e.Value.(*room)
			or := &APIRoom{}
			or.ID = r.ID
			or.Name = r.Name
			or.Area = APIArea{}
			or.Area.ID = r.Area.ID
			or.Area.Name = r.Area.Name
			or.Area.NumPlayers = r.Area.numPlayers
			or.Area.Age = r.Area.age

			or.Description = r.Description
			for _, x := range r.Exits {
				ex := &APIExit{}
				ex.ID = x.ID
				ex.Keyword = x.Keyword
				ex.Description = x.Description
				ex.Dir = x.Dir
				ex.RoomID = x.RoomID
				ex.Key = getItem(x.Key)
				ex.Flags = getExitFlags(x.Flags)
				or.Exits = append(or.Exits, ex)
			}

			for _, mx := range r.Mobs {
				m1 := &APIMob{}
				m1.ID = mx.ID
				m1.Name = mx.Name
				m1.Level = mx.Level
				m1.Job = mx.Job
				m1.Race = mx.Race
				m1.RoomID = mx.Room.ID
				or.Mobs = append(or.Mobs, m1)
			}

			rooms = append(rooms, or)
		}
	}

	type Result struct {
		TotalRooms int
		TotalAreas int
		Results    []*APIRoom
	}
	results := &Result{}
	results.Results = rooms
	results.TotalRooms = roomList.Len()
	results.TotalAreas = areaList.Len()

	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(output))
}
