$(document).ready(function () {
    $("input[id^='create_subscription']").click(
        function () {
            var id = this.id.replace("create_subscription_","");
            $.ajax({
                url: `/api/customer/${id}/subscription`,
                type: "POST",
                data: $('#subscription_form').serialize(),
                success: function () {
                    location.replace(`admin/customer/${id}/subscriptions/`);
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Incorrect Data'
                    $('#subscription_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='delete_subscription']").click(
        function () {
            var id = this.id.replace("delete_subscription_","");
            var customerId = this.customerId;
            $.ajax({
                url: `/api/customer${customerId}/subscription/${id}`,
                type: "DELETE",
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#subscriptions_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='enable_subscription']").click(
        function () {
            var id = this.id.replace("enable_subscription_","");
            var subscription = this.title.replace("Enable subscription ", "")
            $.ajax({
                url: `/api/customer/${customerId}/subscription/${id}`,
                type: "PATCH",
                data: {name: subscription, status: "on"},
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#subscriptions_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='disable_subscription']").click(
        function () {
            var id = this.id.replace("disable_subscription_","");
            var subscription = this.title.replace("Disable subscription ", "")
            $.ajax({
                url: `/api/subscription/${id}`,
                type: "PATCH",
                data: {name: subscription, status: "off"},
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#subscriptions_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );
});
