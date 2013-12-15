function Alert () {}

Alert.warning = function(message) {
	var div = $('<div class="alert"></div>');
	div.append('<button type="button" class="close" data-dismiss="alert">&times;</button>');
	div.append('<strong>Warning!</strong>');
	div.append(message);
	return div;
}