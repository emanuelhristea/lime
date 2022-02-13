$(document).ready(function () {
	$("#create_customer").click(function () {
		$.ajax({
			url: "/api/customer",
			type: "POST",
			data: $("#customer_form").serialize(),
			success: function () {
				location.replace("/admin/");
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != ""
						? err.msg
						: "Incorrect Data";
				$("#customer_result").html(`<p>${message}</p>`);
			},
		});
		return false;
	});

	$("a[id^='delete_customer']").click(function () {
		var id = this.id.replace("delete_customer_", "");
		if (!confirm(`Are you sure you want to delete ${this.title.replace("Delete ", "")}?`)) {
			return
		}
		$.ajax({
			url: `/api/customer/${id}`,
			type: "DELETE",
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != ""
						? err.msg
						: "Cannot delete customer";
				$("#customers_result").html("<p><b>" + message + "</p><b>");
			},
		});
		return false;
	});

	$("#update_customer").click(function () {
		var id = this.dataset.customerId;
		$.ajax({
			url: `/api/customer/${id}`,
			type: "PATCH",
			data: $("#customer_form").serialize(),
			success: function () {
				location.replace("/admin");
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != ""
						? err.msg
						: "Incorrect Data";
				$("#customer_result").html("<p>" + message + "</p>");
			},
		});
		return false;
	});

	$("a[id^='toggle_customer']").click(function () {
		var id = this.id.replace("toggle_customer_", "");
		$.ajax({
			url: `/api/customer/${id}`,
			type: "PATCH",
			data: {
				name: this.dataset.customer,
				status: this.dataset.status == "true" ? "" : "on",
			},
			success: function () {
				location.reload()
				// $.ajax({
				// 	url: `/admin/customer/${id}`,
				// 	type: "GET",
				// 	success: function (data, status, xhr) {
				// 		$(`#customer_row_${id}`).html(data);
				// 	}
				// });
				
				// $(`#activate_customer_${id}`).fadeOut(500, function(){
					
				// 	$(`#activate_customer_${id}`).fadeIn().delay(500)
				// })
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != ""
						? err.msg
						: "Cannot update customer";
				$("#customers_result").html("<p><b>" + message + "</p><b>");
			},
		});
		return false;
	});
});
