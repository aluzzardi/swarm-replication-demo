var previous = {}

function worker() {
	$.ajax("/status", {
		success: function(data, text, jq) {
			$("#status").empty();
			_.each(data, function(value, key) {
				var tpl = _.template($("#card-tpl").html());
				value.Address = key;
				card = $(tpl(value));
				$("#status").append(card);
				if (_.has(previous, key)) {
					if (value.Error && !previous[key].Error) {
						card.transition('flash');
					}
					if (value.Role == 'primary' && previous[key].Role == 'replica') {
						card.transition('tada');
					}
				}
			});

			previous = data;
		},
		complete: function() {
			_.delay(worker, 1500)
		}
	});
}

$(document).ready(function() {
	worker();
});
