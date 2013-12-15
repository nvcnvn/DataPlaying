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
	Histogram($('#hostogram'), data);
	QuantilePlot($('#quantileplot'), data);
	QQPlot($('#qqplot'), data);
}

function normal_sample() {
	var x = 0, y = 0, rds, c;
	do {
		x = Math.random() * 2 - 1;
		y = Math.random() * 2 - 1;
		rds = x * x + y * y;
	} while (rds == 0 || rds > 1);
	c = Math.sqrt(-2 * Math.log(rds) / rds); // Box-Muller transform
	return x * c; // throw away extra sample y * c
}

function QQPlot(el, data) {
	$.each(data, function(idx, item){
		//generate normal sample
		var samples = [];
		for(i = 0; i <= item.Data.length; i++){
			samples.push(item.Mean + Math.sqrt(item.Variance)*normal_sample())
		}
		samples.sort();

		var div = $('<div></div>', {height: 400});
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
                name: 'test',
                color: 'rgba(223, 83, 83, .5)',
                data: scatters
            }]
        });
	});
}

function QuantilePlot(el, data) {
	$(el).empty();
	$.each(data, function(idx, item){
		var div = $('<div></div>', {height: 400});
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
                name: 'test',
                color: 'rgba(223, 83, 83, .5)',
                data: scatters
            }]
        });
	});
}

function Histogram(el, data) {
	$(el).empty();

	$.each(data, function(idx, item){
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
	});
}