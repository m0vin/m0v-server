// Table
	var $th = '', $n = 10, $tr = [], $table, $i, $ii, $j, $pageCount;
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
			$table.insertAdjacentHTML("afterend","<div id='buttons' ></div");
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

	
	var tdn = document.getElementById("tdn").innerHTML; //tabledivname
	paginate(tdn);

