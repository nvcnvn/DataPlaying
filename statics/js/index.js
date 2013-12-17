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
		    url: "/index/ajax",
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
})

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
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#histogram'+idx+'">Histogram</a>')));
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#quantileplot'+idx+'">QuantilePlot</a>')));
		ul.append($('<li></li>').append($('<a data-toggle="tab" href="#qqplot'+idx+'">QQPlot</a>')));

		var div = $('<div></div>', {
			class: 'tab-content'
		});

		var histogram = $('<div></div>', {
			class: "tab-pane active",
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

function QQPlot(el, item) {
	//generate normal sample
	var samples = jstat.seq(-3, 3, item.Data.length);
	samples.sort();

	var div = $('<div></div>', {height: 400, width: 400});
	var scatters = [];
	for(i = 0; i <= item.Data.length; i++){
		scatters.push([samples[i], item.Data[i]]);
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
   //  , function(chart) {
   //  	console.log(chart)
   //  	chart.renderer.path(['M', 0, 0, 'L', 400, 400])
   //  	.attr({
			// 'stroke-width': 2,
			// stroke: 'red'
   //      }).add();
   //  }
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