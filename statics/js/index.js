$(function() {
	$('#btProcess').click(function(){
		var a = [];
		$.each($('.checkField:checked'), function(idx, item){
			a.push(parseInt($(item).attr('data-id')));
		});

		if(a.length == 0){
			$('#alert').append(Alert.warning('you must select one or more attrs'));
			return;
		}

		$.ajax({
		    type: "POST",
		    url: "/ajax/summary",
		    data: JSON.stringify(a),
		    contentType: "application/json; charset=utf-8",
		    dataType: "json",
		    success: function(data){
		    	BoxPlot(data);
		    	Graph(data);
		    },
		    failure: function(errMsg) {
		        alert(errMsg);
		    }
		});
	});

	$('#btHideSlected').click(function(){
		$.each($('.checkField:checked'), function(idx, item){
			var dataId = parseInt($(item).attr('data-id')) + 1;
			$('#dataset th:nth-child('+dataId+')').hide();
			$('#dataset td:nth-child('+dataId+')').hide();
			$(item).prop('checked', false);
		});	
	});

	$('#btShowAll').click(function(){
		$('#dataset th').show();
		$('#dataset td').show();
	});

	$('#btHideTable').click(function(){
		if($(this).val() == 'Hide Table'){
			$(this).val('Show Table');
		}else{
			$(this).val('Hide Table');

		}
		$('#dataset').toggle();
	});

	$('#btShowScatter').click(function(){
		var a = [];
		$.each($('.checkField:checked'), function(idx, item){
			a.push(parseInt($(item).attr('data-id')));
		});

		if(a.length != 2){
			$('#alert').append(Alert.warning('you must select 2 attrs'));
			return;
		}

		$.ajax({
		    type: "POST",
		    url: "/ajax/data",
		    data: JSON.stringify(a),
		    contentType: "application/json; charset=utf-8",
		    dataType: "json",
		    success: function(data){
		    	ScatterChart(data);
		    },
		    failure: function(errMsg) {
		        alert(errMsg);
		    }
		});
	});

	$('#btShowQQ').click(function(){
		var a = [];
		$.each($('.checkField:checked'), function(idx, item){
			a.push(parseInt($(item).attr('data-id')));
		});

		if(a.length != 2){
			$('#alert').append(Alert.warning('you must select 2 attrs'));
			return;
		}

		$.ajax({
		    type: "POST",
		    url: "/ajax/sorteddata",
		    data: JSON.stringify(a),
		    contentType: "application/json; charset=utf-8",
		    dataType: "json",
		    success: function(data){
		    	$('#qqModal').modal();
		    	$('#qqChart').empty();
		    	drawQQPlot($('#qqChart'), data[0].Data, data[1].Data);
		    },
		    failure: function(errMsg) {
		        alert(errMsg);
		    }
		});
	});
});


function drawScatter(el, a, b) {
	el.empty();
	var scatters = [];
	for(i = 0; i<a.Data.length; i++){
		scatters.push([a.Data[i], b.Data[i]]);
	}

    el.highcharts({
        chart: {
            type: 'scatter',
            zoomType: 'xy'
        },
        title: {
            text: a.Name+' Versus '+b.Name
        },
        xAxis: {
            title: {
                enabled: true,
                text: a.Name
            },
            startOnTick: true,
            endOnTick: true,
            showLastLabel: true
        },
        yAxis: {
            title: {
                text: b.Name
            }
        },
        plotOptions: {
            scatter: {
                marker: {
                    radius: 5,
                    states: {
                        hover: {
                            enabled: true,
                            lineColor: 'rgb(100,100,100)'
                        }
                    }
                },
                states: {
                    hover: {
                        marker: {
                            enabled: false
                        }
                    }
                },
                tooltip: {
                    pointFormat: '{point.x}, {point.y}'
                }
            }
        },
        series: [{
            name: '',
            color: 'rgba(223, 83, 83, .5)',
            data: scatters
        }]
    });
}

function ScatterChart(data) {
	var a = data[0];
	var b = data[1];
	drawScatter($('#scatterChart0'), a, b);
	drawScatter($('#scatterChart1'), b, b);
	drawScatter($('#scatterChart2'), a, a);
	drawScatter($('#scatterChart3'), b, a);
	$('#scatterTable').show();
}

function BoxPlot(obj) {
	var categories = [];
	var data = [];

	$.each(obj, function(idx, item){
		categories.push(item.Name);
		data.push([item.Min, item.Quartiles[1], item.Quartiles[2], item.Quartiles[3], item.Max]);
	});

	$('#boxplot').width($('#dataset').width());
    $('#boxplot').highcharts({
        chart: {
            type: 'boxplot'
        },
        title: {
            text: 'Box Plot'
        },
        legend: {
            enabled: false
        },
        xAxis: {
            categories: categories,
            title: {
                text: '----------------------'
            }
        },
        series: [{
            name: 'Observations',
            data: data,
            tooltip: {
                headerFormat: '<em>{point.key}</em><br/>'
            }
        }, ]
    });
}

function Graph(data) {
	$('#graph').empty();
	$.each(data, function(idx, item){
		var ul = $('<ul data-tabs="tabs" class="nav nav-tabs"></ul>');
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#summary'+idx+'">Summary</a>')));
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#histogram'+idx+'">Histogram</a>')));
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#quantileplot'+idx+'">QuantilePlot</a>')));
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#qqplot'+idx+'">QQPlot</a>')));

		var div = $('<div></div>', {
			class: 'tab-content'
		});

		var summary = $('<div></div>', {
			class: "tab-pane active",
			id: 'summary'+idx
		});
		div.append(summary);
		summary.append('<p>Sorted:<div style="width:100%;max-height:2em;overflow-x:scroll;">'+item.Data+'</div></p>');
		summary.append('<p>Min:'+item.Min+'</p>');
		summary.append('<p>Max:'+item.Max+'</p>');
		summary.append('<p>Range:'+(item.Max-item.Min)+'</p>');
		summary.append('<p>Mean:'+item.Mean+'</p>');
		summary.append('<p>Mode:'+item.Mode+'</p>');
		summary.append('<p>Median:'+item.Median+'</p>');
		summary.append('<p>Variance:'+item.Variance+'</p>');
		summary.append('<p>StdDev:'+ Math.sqrt(item.Variance) +'</p>');
		summary.append('<p>Quartiles: Q1: '+item.Quartiles[1]+', Q2: ' +item.Quartiles[2]+ ', Q3: ' +item.Quartiles[3]+ '</p>');
		summary.append('<p>Interquartiles Range:'+(item.Quartiles[3]-item.Quartiles[1])+'</p>');


		var histogram = $('<div></div>', {
			class: "tab-pane",
			id: 'histogram'+idx
		});
		div.append(histogram);
		var quantileplot = $('<div></div>', {
			class: "tab-pane",
			id: 'quantileplot'+idx
		});
		div.append(quantileplot);
		var qqplot = $('<div></div>', {
			class: "tab-pane",
			id: 'qqplot'+idx
		});
		div.append(qqplot);

		$('#graph').append($('<div class="hero-unit"></div>').append('<h1>'+item.Name+'</h1>').append(ul).append(div));
		Histogram(histogram, item);
		QuantilePlot(quantileplot, item);
		QQPlot(qqplot, item);
	});
}

function drawQQPlot(el, a, b) {
	var div = $('<div></div>', {height: 400, width: 400});
	var scatters = [];
	for(i = 0; i <= a.length; i++){
		scatters.push([a[i], b[i]]);
	}

	$(el).append(div);
	div.highcharts({
        chart: {
            type: 'scatter',
            zoomType: 'xy'
        },
        series: [{
            name: 'Q-Q Plot',
            color: 'rgba(223, 83, 83, .5)',
            data: scatters
        }]
    });
}

function QQPlot(el, item) {
	//generate normal sample
	var samples = jstat.seq(-3, 3, item.Data.length);
	samples.sort();

	drawQQPlot(el, samples, item.Data);
}

function QuantilePlot(el, item) {
	var div = $('<div></div>', {height: 400, width: 400});
	var scatters = [];
	for(i = 0; i <= item.Data.length; i++){
		scatters.push([i/item.Data.length, item.Data[i]]);
	}

	$(el).append(div);
	div.highcharts({
        chart: {
            type: 'scatter',
            zoomType: 'xy'
        },
        series: [{
            name: 'Quantile Plot',
            color: 'rgba(223, 83, 83, .5)',
            data: scatters
        }]
    });
}

function Histogram(el, item) {
	var div = $('<div></div>', {height: 400});
	$(el).append(div);

	var series = [];
	for(var i = 0; i < item.Lookup.Value.length; i++){
		series.push({
			name: item.Lookup.Value[i].toString(),
			data: [item.Lookup.Count[i]],
            tooltip: {
            	headerFormat: ''
        	}
		});
	}

    div.highcharts({
        chart: {
            type: 'column'
        },
        xAxis: {
            categories: '',
            title: {
                text: null
            }
        },
        title: {
        	text: 'Histogram ' + item.Name,
        },
        yAxis: {
            min: 0,
            title: {
                text: 'appears (times)',
                align: 'high'
            },
            labels: {
                overflow: 'justify'
            }
        },
        tooltip: {
            valueSuffix: ' times'
        },
        plotOptions: {
            bar: {
                dataLabels: {
                    enabled: true
                }
            }
        },
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'top',
            x: -40,
            y: 100,
            floating: true,
            borderWidth: 1,
            backgroundColor: '#FFFFFF',
            shadow: true
        },
        credits: {
            enabled: false
        },
        series: series
    });
}