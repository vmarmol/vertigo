package static

var scriptsJs = `
function getStateNode(state) {
	if (state == "ok") {
		return $("<span>").addClass("label label-success").text("ok");
	} else if (state == "warming-up") {
		return $("<span>").addClass("label label-warning").text("warming up");
	} else if (state == "migrating") {
		return $("<span>").addClass("label label-default").text("migrating");
	} else {
		return $("<span>").addClass("label label-danger").text("unknown");
	}
}

// Gets and updates the Vertigo instances.
function updateInstances() {
	$.getJSON("/api/instances", function(data) {
		// Add header.
		var table = $("<table>").addClass("table")
			.append($("<tr>")
				.append($("<th>").text("Name"))
				.append($("<th>").text("State"))
				.append($("<th>").text("CPU"))
				.append($("<th>").text("Memory")));

		// Add a row per instance.
		for (var i = 0; i < data.length; i++) {
			var instance = data[i];
			table.append($("<tr>")
				.append($("<td>").text(instance.name))
				.append($("<td>").append(getStateNode(instance.state)))
				.append($("<td>").text(instance.cpu_usage + "%"))
				.append($("<td>").text(instance.memory_usage + "%")));
		}

		$("#instances").empty().append(table);
	});
}

// Update the uptime.
function updateUptime() {
	$.getJSON("/api/uptime", function(data) {
		$("#service-uptime").empty().text(data.uptime);
	}).fail(function(jqxhr, textStatus, error) {console.log("failed: " + textStatus + ", " + error)});
}

// Starts execution of the page.
function startPage() {
	setInterval(function() {
		updateInstances();
		updateUptime();
	}, 1000);
}

// Saves a change in QPS value.
function saveQps() {
	var qps = $("#service-qps").val();
	$.getJSON("/api/qps?qps=" + qps, function(data) {
		console.log("Set QPS to:  " + qps);
	});
}
`
