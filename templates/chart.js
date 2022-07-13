{{define "chartjs"}}
	<script src="/static/code/highcharts.js"></script>
	<script src="/static/code/modules/data.js"></script>
	<script src="/static/code/modules/exporting.js"></script>
	<script src="/static/code/modules/export-data.js"></script>
	<script src="/static/code/modules/accessibility.js"></script>
    <script type="text/javascript">
        Highcharts.getJSON(
            'http://localhost:8080/api/summary/{{.Freq}}/{{.Id}}/{{.Start}}/{{.End}}',
            function (data) {
                var processed_json = new Array();
                var processed_json1 = new Array();
                //var d = []; //new Array();
                for (i = 0; i < data.length; i++) {
                    //let d = new Date(Date.parse(data[i].timestamp));
                    //d.push((new Date(data[i].timestamp)).getTime());
                    //console.log(new Date(data[i].timestamp).getTime() / 1000);
                    processed_json.push({x: (new Date(data[i].timestamp)).getTime(), y: parseFloat(data[i].exportactiveenergy)});
                    processed_json1.push({x: (new Date(data[i].timestamp)).getTime(), y: parseFloat(data[i].activepwrmax)/1000.0});
                    //processed_json.push([d, parseFloat(data[i].exportactiveenergy)]);
                    //processed_json.push([(data[i].timestamp), parseFloat(data[i].exportactiveenergy)]);
                }
                Highcharts.chart('container', {
                    chart: {
                        zoomType: 'xy'
                    },
                    title: {
                        text: 'Power & Energy ' + {{.Freq}}
                    },
                    subtitle: {
                        text: document.ontouchstart === undefined ?
                            'Click and drag in the plot area to zoom in' : 'Pinch the chart to zoom in'
                    },
                    xAxis: {
                        type: 'datetime',
                        crosshair: true
                    },
                    yAxis: [{
                        labels: {
                            format: '{value} kwh',
                            style: {
                                color: Highcharts.getOptions().colors[0]
                            }
                        },
                        title: {
                            text: 'Energy',
                            style: {
                                color: Highcharts.getOptions().colors[0]
                            }
                        }
                    }, { // secondary
                        title: {
                            text: 'Power',
                            style: {
                                color: Highcharts.getOptions().colors[1]
                            }
                        },
                        labels: {
                            format: '{value} kw',
                            style: {
                                color: Highcharts.getOptions().colors[1]
                            }
                        },
                        opposite: true
                    }],
                    tooltip: {
                        shared: true
                    },
                    legend: {
                        /*enabled: false*/
                        layout: 'vertical',
                        align: 'left',
                        verticalAlign: 'top',
                        floating: true,
                        x: 100,
                        y: 60
                    },
                    /*plotOptions: {
                        area: {
                            fillColor: {
                                linearGradient: {
                                    x1: 0,
                                    y1: 0,
                                    x2: 0,
                                    y2: 1
                                },
                                stops: [
                                    [0, Highcharts.getOptions().colors[0]],
                                    [1, Highcharts.color(Highcharts.getOptions().colors[0]).setOpacity(0).get('rgba')]
                                ]
                            },
                            marker: {
                                radius: 2
                            },
                            lineWidth: 1,
                            states: {
                                hover: {
                                    lineWidth: 1
                                }
                            },
                            threshold: null
                        }
                    },*/

                    series: [{
                        type: 'area',
                        name: 'Energy',
                        tooltip: {
                            valueSuffix: ' kwh'
                        },
                        yAxis: 0,
                        data: processed_json //[{
                            //x: d,
                            //y: processed_json //data
                        //}]
                    }, {
                        type: 'column',
                        name: 'Active Power',
                        tooltip: {
                            valueSuffix: ' kw'
                        },
                        yAxis: 1,
                        data: processed_json1 //[{
                            //x: d,
                            //y: processed_json //data
                        //}]
                    }]
                });
            }
        );
    </script>
{{end}}
