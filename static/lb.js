
	var pubs = '[{"id":5,"latitude":12.97,"longitude":77.55,"hash":9876543,"created":"2020-08-07T16:33:41.646037Z","email":0},{"id":4,"latitude":12.93,"longitude":77.63,"hash":987654,"created":"2020-08-06T15:19:51.295684Z","email":0},{"id":1,"latitude":12.95,"longitude":77.6,"hash":729808954,"created":"1970-01-01T00:04:13Z","email":0},{"id":6,"latitude":13,"longitude":77.5,"hash":1260454696,"created":"1970-01-01T00:03:29Z","email":0}]';


	function loadPubsLocal() {
		var myObj = JSON.parse(pubs);
		//var $table = document.getElementById("pubs");
		var tbody = '';
		for(p in myObj) {
			tbody += '<tr>';
			tbody += '<td>' + myObj[p].hash + '</td>';
			tbody += '<td>' + myObj[p].latitude + '</td>';
			tbody += '<td>' + myObj[p].longitude + '</td>';
			tbody += '<td>' + myObj[p].created + '</td>';
			tbody += '</tr>';
		}
		document.getElementById("tbody").innerHTML = tbody;
	}


	function loadPubs() {
		var xhttp = new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				var myObj = JSON.parse(this.responseText);
				//console.log(myObj[0]);
				//document.getElementById("demo").innerHTML =
				//this.responseText;
				/*var thead = '<tr>'
				thead += '<th>Hash</th>';
				thead += '<th>Latitude</th>';
				thead += '<th>Longitude</th>';
				thead += '<th>Created</th>';
				thead += '</tr>';
				document.getElementById("thead").innerHTML = thead;*/
				var tbody = '';
				for(p in myObj) {
					tbody += '<tr>';
					var rank = 1 + p/1;
					tbody += '<td>' + rank + '</td>';
					tbody += '<td>' + myObj[p].hash + '</td>';
					tbody += '<td>' + myObj[p].longitude + '</td>';
					tbody += '<td>' + myObj[p].latitude + '</td>';
					tbody += '<td>' + myObj[p].created + '</td>';
					tbody += '<td>' + myObj[p].kwp + '</td>';
					tbody += '<td>' + myObj[p].kwlast + '</td>';
					tbody += '<td>' + myObj[p].kwhday + '</td>';
					tbody += '<td>' + myObj[p].kwhlife + '</td>';
					tbody += '</tr>';
				}
				document.getElementById("tbody").innerHTML = tbody;
				var tdn = document.getElementById("tdn").innerHTML; //tabledivname
				paginate(tdn);
			}
		};
		xhttp.open("GET", "/api/pubs", true);
		xhttp.send();
	}

// Table

	var $th = '', $n = 25, $tr = [], $table, $i, $ii, $j, $pageCount;
	function paginate(tdn) {
		$table = document.getElementById(tdn);
		var $rowCount = $table.rows.length,
		$firstRow = $table.rows[0].firstElementChild.tagName,
		$hasHead = ($firstRow === "TH");
		$i,$ii,$j = ($hasHead)?1:0;
		$th = ($hasHead?$table.rows[(0)].outerHTML:"");
		console.log($firstRow);
		$pageCount = Math.ceil($rowCount / $n);
		if ($pageCount > 1) {
			// assign each row outHTML (tag name & innerHTML) to the array
			for ($i = $j,$ii = 0; $i < $rowCount; $i++, $ii++)
			$tr[$ii] = $table.rows[$i].outerHTML;
			// create a div block to hold the buttons
			$table.insertAdjacentHTML("afterend","<div id='buttons' style='display:inline;'></div>");
			// the first sort, default page is the first one
			sort(1);
		}
	}
	// ($p) is the selected page number. it will be generated when a user clicks a button
	function sort($p) {
		/* create ($rows) a variable to hold the group of rows
		** to be displayed on the selected page,
		** ($s) the start point .. the first row in each page, Do The Math
		*/
		var $rows = $th,$s = (($n * $p)-$n);
		for ($i = $s; $i < ($s+$n) && $i < $tr.length; $i++)
		$rows += $tr[$i];

		// now the table has a processed group of rows ..
		$table.innerHTML = $rows;
		// create the pagination buttons
		document.getElementById("buttons").innerHTML = pageButtons($pageCount,$p);
		// CSS Stuff
		document.getElementById("id"+$p).setAttribute("class","active");
	}


	// ($pCount) : number of pages,($cur) : current page, the selected one ..
	function pageButtons($pCount,$cur) {
		/* this variables will disable the "Prev" button on 1st page
		and "next" button on the last one */
		var $prevDis = ($cur == 1)?"disabled":"",
		$nextDis = ($cur == $pCount)?"disabled":"",
		/* this ($buttons) will hold every single button needed
		** it will creates each button and sets the onclick attribute
		** to the "sort" function with a special ($p) number..
		*/
		$buttons = "<input type='button' class='button' value='&lt;&lt; Prev' onclick='sort("+($cur - 1)+")' "+$prevDis+">";
		for ($i=1; $i<=$pCount;$i++)
		$buttons += "<input type='button' class='button' id='id"+$i+"'value='"+$i+"' onclick='sort("+$i+")'>";
		$buttons += "<input type='button' class='button' value='Next &gt;&gt;' onclick='sort("+($cur + 1)+")' "+$nextDis+">";
		return $buttons;
	}

	
	paginate("pubs");
	//loadPubs();
	//loadPubsLocal();
