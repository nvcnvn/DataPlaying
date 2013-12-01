$(function() {
	$('#process').click(function(){
		var a = [];
		$.each($('.checkField:checked'), function(idx, item){
			a.push(parseInt($(item).attr('data-id')));
		});

		$.ajax({
		    type: "POST",
		    url: "/index/ajax",
		    data: JSON.stringify(a),
		    contentType: "application/json; charset=utf-8",
		    dataType: "json",
		    success: function(data){
		    	BoxPlot(data);
		    	Histogram(data);
		    },
		    failure: function(errMsg) {
		        alert(errMsg);
		    }
		});
	});
})

function BoxPlot(obj) {
	var categories = [];
	var data = [];

	$.each(obj, function(idx, item){
		categories.push(item.Name);
		data.push([item.Min, item.Quartiles[0], item.Quartiles[1], item.Quartiles[2], item.Max]);
	});

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

function Histogram(data) {
	$('#histogram').empty();
	$.each(data, function(idx, item){
		var div = $('<div></div>', {height: 400});
		$('#histogram').append(div);
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