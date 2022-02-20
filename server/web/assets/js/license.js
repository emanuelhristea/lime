$(document).ready(function () {
	$("input[id^='create_license']").click(function () {
		var subscriptionId = this.id.replace("create_license_", "");
		var customerId = this.dataset.customerId;
		$.ajax({
			url: `/api/licenses/${subscriptionId}`,
			type: "POST",
			data: $("#license_form").serialize(),
			success: function () {
				location.replace(`/admin/customer/${customerId}/subscriptions/`);
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Incorrect Data";
				$("#license_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});

    $("a[id^='toggle_license_']").click(function () {
		var id = this.id.replace("toggle_license_", "");
		$.ajax({
			url: `/api/license/${id}`,
			type: "PATCH",
			data: {
                status: this.dataset.status == "true" ? "" : "on", 
            },
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot modify license";
				$("#subscriptions_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});

	$("a[id^='delete_license_']").click(function () {
		var id = this.id.replace("delete_license_", "");
		$.ajax({
			url: `/api/license/${id}`,
			type: "DELETE",
			success: function () {
				location.reload();
			},
			error: function (xhr, status, error) {
				var err = eval("(" + xhr.responseText + ")");
				var message =
					typeof err.msg !== "undefined" && err.msg != null && err.msg != "" ? err.msg : "Cannot delete subscription";
				$("#subscriptions_result").html("<p><b>" + message + "</b></p>");
			},
		});
		return false;
	});
});
