package main

import (
	gsdbus "github.com/coreos/go-systemd/dbus"
	"log"
	// "reflect"
	"strings"
	"time"
)

// running in container for now... run on host CoreOS later:
// docker run -v /var/run/dbus:/var/run/dbus -v /run/systemd:/run/systemd -v /usr/bin/systemctl:/usr/bin/systemctl -v /etc/systemd/system:/etc/systemd/system -d --net=host 127.0.0.1:5000/done-go /usr/sbin/sshd -D

// subscribe and print systemd events

func main() {

	log.Println("connecting to dbus")
	conn, err := gsdbus.New()
	if err != nil {
		log.Fatal("error connecting to dbus: ", err)
	}

	log.Println(conn)

	// connType := reflect.TypeOf(conn)
	// for i := 0; i < connType.NumMethod(); i++ {
	// 	method := connType.Method(i)
	// 	log.Println(method.Name)
	// }

	log.Println("subscribing")
	err = conn.Subscribe()
	if err != nil {
		log.Fatal("error subscribing: ", err)
	}

	log.Println("calling SubscribeUnits")
	evChan, errChan := conn.SubscribeUnits(time.Second)

	log.Println(evChan, errChan)

	// SubscribeUnits evChan returns a map of unit_name -> *UnitStatus:
	// type UnitStatus struct {
	// 	Name        string          // The primary unit name as string
	// 	Description string          // The human readable description string
	// 	LoadState   string          // The load state (i.e. whether the unit file has been loaded successfully)
	// 	ActiveState string          // The active state (i.e. whether the unit is currently started or not)
	// 	SubState    string          // The sub state (a more fine-grained version of the active state that is specific to the unit type, which the active state is not)
	// 	Followed    string          // A unit that is being followed in its state by this unit, if there is any, otherwise the empty string.
	// 	Path        dbus.ObjectPath // The unit object path
	// 	JobId       uint32          // If there is a job queued for the job unit the numeric job id, 0 otherwise
	// 	JobType     string          // The job type as string
	// 	JobPath     dbus.ObjectPath // The job object path

	for {
		select {
		case changes := <-evChan:
			// log.Println("got something:  ", changes)

			for unit_name, unit_status := range changes {
				if strings.Contains(unit_name, ".service") {
					if unit_status.ActiveState == "active" {
						short_unit_name := strings.Replace(unit_name, ".service", "", 1)
						log.Println("name =", short_unit_name, ", state =", unit_status.ActiveState)
						// XXX need to filter out non-network services... not enough info available from UnitStatus?????
						// docker should be able to tell if a container is listening on a port... maybe that's why registrator would work
					}
				}
			}

			// tCh, ok := changes[target]

			// Just continue until we see our event.
			// if !ok {
			//	log.Println("not ok????")
			//	continue
			// }

			// if tCh.ActiveState == "active" && tCh.Name == target {

		case err = <-errChan:
			log.Fatal(err)
		}
	}
}
