package main

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/cluster"
	"github.com/hazelcast/hazelcast-go-client/logger"
	"github.com/hazelcast/hazelcast-go-client/serialization"
	"log"
	"time"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("Person: %s (%d)", p.Name, p.Age)
}

func (p Person) WritePortable(writer serialization.PortableWriter) error {
	writer.WriteString("name", p.Name)
	writer.WriteInt16("age", int16(p.Age))
	return nil
}

func (p *Person) ReadPortable(reader serialization.PortableReader) error {
	p.Name = reader.ReadString("name")
	p.Age = int(reader.ReadInt16("age"))
	return nil
}

func (p *Person) FactoryID() int32 {
	return 824811
}

func (p *Person) ClassID() int32 {
	return 1
}

type PersonFactory struct {
}

func (p PersonFactory) Create(classID int32) serialization.Portable {
	return &Person{}
}

func (p PersonFactory) FactoryID() int32 {
	return 824811
}

func lifecycleStateChangeHandler(event hazelcast.LifecycleStateChanged) {
	var state string
	switch event.State {
	case hazelcast.LifecycleStateStarting:
		state = "STARTING"
	case hazelcast.LifecycleStateStarted:
		state = "STARTED"
	case hazelcast.LifecycleStateShuttingDown:
		state = "SHUTTING DOWN"
	case hazelcast.LifecycleStateShutDown:
		state = "SHUT DOWN"
	case hazelcast.LifecycleStateClientConnected:
		state = "CLIENT CONNECTED"
	case hazelcast.LifecycleStateClientDisconnected:
		state = "CLIENT DISCONNECTED"
	default:
		state = fmt.Sprintf("UNKNOWN STATE: %d", event.State)
	}
	log.Printf("State Changed: %s", state)
}

func memberStateChangeHandler(event cluster.MembershipStateChanged) {
	switch event.State {
	case cluster.MembershipStateAdded:
		log.Printf("Member Added: %s @%s", event.Member.UUID(), event.Member.Address())
	case cluster.MembershipStateRemoved:
		log.Printf("Member Removed: %s @%s", event.Member.UUID(), event.Member.Address())
	}
}

func main() {
	configBuilder := hazelcast.NewConfigBuilder()
	configBuilder.Logger().
		SetLevel(logger.TraceLevel)
	configBuilder.Cluster().
		SetAddrs("localhost:5701")
	configBuilder.Serialization().
		AddPortableFactory(&PersonFactory{})
	configBuilder.AddLifecycleListener(lifecycleStateChangeHandler)
	configBuilder.AddMembershipListener(memberStateChangeHandler)
	client, err := hazelcast.StartNewClientWithConfig(configBuilder)
	if err != nil {
		log.Fatal(err)
	}
	// get a map
	people, err := client.GetMap("people")
	if err != nil {
		log.Fatal(err)
	}
	personName := "Jane Doe"
	// set a value in the map
	if err = people.Set(personName, &Person{personName, 30}); err != nil {
		log.Fatal(err)
	}
	// get a value from the map
	person, err := people.Get(personName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved %s.\n", person)

	time.Sleep(5 * time.Second)
	client.Shutdown()
	time.Sleep(1 * time.Second)
}
