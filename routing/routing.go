package routing

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"bytes"
	"sort"
)

type Router struct {
	Name          rune
	RoutingTable  map[rune]routingTableEntry
	adjacencyList map[rune]int
}

type routingTableEntry struct {
	direction rune
	distance  int
}

type updatingInformation struct {
	to    rune
	from  rune
	table map[rune]int
}

func (this Router) sendRoutingInformation() map[rune]int {
	res := make(map[rune]int, len(this.RoutingTable))

	for k, v := range this.RoutingTable {
		res[k] = v.distance
	}
	return res
}

func (this Router) getNeighbours() []rune {
	res := make([]rune, 0, len(this.adjacencyList))
	for n, _ := range this.adjacencyList {
		res = append(res, n)
	}
	return res
}

func (this *Router) updateInformation(source rune, externalTable map[rune]int) (changed bool) {
	changed = false
	incoming := this.RoutingTable[source]
	for router, distance := range externalTable {
		entry, known := this.RoutingTable[router]
		if known && distance + incoming.distance < entry.distance || ! known {

			this.RoutingTable[router] = routingTableEntry{incoming.direction, distance + incoming.distance}
			changed = true
		}

	}
	return changed
}

func ParseFile() (map[rune]Router, error) {

	result := make(map[rune]Router)
	var routerName rune
	file, err := os.Open("router.txt")
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if line[0] == "" {
			continue
		}
		if len(line) == 1 {
			lineSlice := []rune(line[0])
			routerName = lineSlice[0]
			var router Router
			router.adjacencyList = make(map[rune]int)
			router.RoutingTable = make(map[rune]routingTableEntry)
			router.Name = routerName
			result[routerName] = router

		} else if len(line) == 2 {
			dist, _ := strconv.Atoi(line[1])
			router := result[routerName]
			router.adjacencyList[[]rune(line[0])[0]] = dist
		}

	}
	for routerName, router := range result {
		router.RoutingTable[routerName] = routingTableEntry{routerName, 0}
		for neighbour, cost := range router.adjacencyList {
			router.RoutingTable[neighbour] = routingTableEntry{neighbour, cost}
		}
	}
	return result, nil
}

func DistanceVectorRoutingStep(routers map[rune]Router) (changed bool) {

	changed = false;
	jobs := make([]updatingInformation, len(routers) * 12)
	for routerName, router := range routers {
		routingTable := make(map[rune]int)
		for k, v := range router.RoutingTable {
			routingTable[k] = v.distance
		}
		for neighbour, _ := range router.adjacencyList {
			jobs = append(jobs, updatingInformation{neighbour, routerName, routingTable})
		}
	}

	for _, job := range jobs {
		router := routers[job.to]
		if router.updateInformation(job.from, job.table) {
			changed = true
		}
	}

	return changed
}

func (this Router)String() string {
	var routers []int

	for k, _ := range this.RoutingTable {
		routers = append(routers, int(k))
	}

	sort.Ints(routers)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Router %c: RoutingTable:{", this.Name))

	for _, k := range routers {

		v := this.RoutingTable[rune(k)]
		buffer.WriteString(fmt.Sprintf(" (%c, %c, %d)", k, v.direction, v.distance))
	}
	buffer.WriteString("} ")
	return buffer.String()
}
