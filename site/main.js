console.log("Starting")

updates = new EventSource("/updates")

updates.addEventListener('update', updateLight)
updates.addEventListener('message', updateLight)
updates.addEventListener('error', disconnect)


var current = "num-5"
var lastEvent = 0

function updateLight(e) {
	console.log("Got event")
	if(e.lastEventId == null || e.id <= lastEvent){
		console.log("skipping duplicate")
		return
	}

	lastEvent = e.lastEventId

	target = `num-${e.data}`
	console.log(`Setting target to ${target}`)

	document.getElementById(current).classList.remove("selected")
	document.getElementById(target).classList.add("selected")

	current = target
}

function disconnect(e) {
	// Assume the server died and the lastEvent will be reset
	lastEvent = 0
	document.getElementById(current).classList.remove("selected")
}

