package static

var indexHtml = `
<html>
<head>
	<title>Vertigo</title>
	<!-- Latest compiled and minified CSS -->
	<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">

	<!-- Optional theme -->
	<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap-theme.min.css">

	<link rel="stylesheet" href="/static/style.css">

	<!-- Latest compiled and minified JavaScript -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	<script src="//netdna.bootstrapcdn.com/bootstrap/3.1.1/js/bootstrap.min.js"></script>
	<script type="text/javascript" src="https://www.google.com/jsapi"></script>

	<script type="text/javascript" src="/static/scripts.js"></script>
</head>
<body>
	<div class="col-sm-8 class="container theme-showcase">
		<div class="page-header">
			<h1>Vertigo <small>Vertical scaling of Docker containers</small></h1>
		</div>
		<div class="panel panel-primary">
			<div class="panel-heading">
				<h2 class="panel-title">Instances</h2>
			</div>
			<div class="panel-body">Running instances answering service queries.</div>
			<div id="instances"></div>
		</div>
		<div class="panel panel-primary">
			<div class="panel-heading">
				<h2 class="panel-title">Service</h2>
			</div>
			<div class="panel-body">Controls and stats for the running service.</div>
			<table id="service" class="table">
				<tr>
					<td>Uptime</td>
					<td id="service-uptime" style="font-style:italic">Updating...</td>
				</tr>
				<tr>
					<td>Latency</td>
					<td id="service-latency" style="font-style:italic">Updating...</td>
				</tr>
				<tr>
					<td id="service-uptime">QPS</td>
					<td>
						<div class="input-group">
							<input id="service-qps" type="text" value="1" class="form-control">
							<span class="input-group-btn">
								<button class="btn btn-default" type="button" onclick="saveQps()">Set</button>
							</span>
						</div>
					</td>
				</tr>
			</table>
		</div>
	</div>
	<script type="text/javascript">
		startPage();
	</script>
</body>
</html>
`
